package hostingde

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &recordResource{}
	_ resource.ResourceWithConfigure   = &recordResource{}
	_ resource.ResourceWithImportState = &recordResource{}
)

// NewRecordResource is a helper function to simplify the provider implementation.
func NewRecordResource() resource.Resource {
	return &recordResource{}
}

// recordResource is the resource implementation.
type recordResource struct {
	client *Client
}

// recordResourceModel maps the DNSRecord resource schema data.
type recordResourceModel struct {
	ID      types.String `tfsdk:"id"`
	ZoneID  types.String `tfsdk:"zone_id"`
	Name    types.String `tfsdk:"name"`
	Type    types.String `tfsdk:"type"`
	Content types.String `tfsdk:"content"`
	TTL     types.Int64  `tfsdk:"ttl"`
}

// Metadata returns the resource type name.
func (r *recordResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_record"
}

// Schema defines the schema for the resource.
func (r *recordResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "DNS record ID",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"zone_id": schema.StringAttribute{
				Description: "ID of DNS zone that the record belongs to.",
				Required:    true,
			},
			"name": schema.StringAttribute{
				Description: "Name of the record. Example: mail.example.com.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of the DNS record. Valid types are A, AAAA, ALIAS, CAA, CERT, CNAME, DNSKEY, DS, MX, NS, NSEC, NSEC3, NSEC3PARAM, NULLMX, OPENPGPKEY, PTR, RRSIG, SRV, SSHFP, TLSA, and TXT.",
				Required:    true,
			},
			"content": schema.StringAttribute{
				Description: "Content of the DNS record.",
				Required:    true,
			},
			"ttl": schema.Int64Attribute{
				Description: "TTL of the DNS record in seconds.",
				Computed:    true,
				Required:    false,
				Optional:    true,
				Default:     int64default.StaticInt64(3600),
			},
		},
	}
}

// Create a new resource
func (r *recordResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan recordResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	record := DNSRecord{
		Name:    plan.Name.ValueString(),
		ZoneID:  plan.ZoneID.ValueString(),
		Type:    plan.Type.ValueString(),
		Content: plan.Content.ValueString(),
		TTL:     int(plan.TTL.ValueInt64()),
	}

	recordReq := RecordsUpdateRequest{
		BaseRequest:  &BaseRequest{},
		ZoneConfigId: plan.ZoneID.ValueString(),
		RecordsToAdd: []DNSRecord{record},
	}

	recordResp, err := r.client.updateRecords(recordReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating records",
			"Could not update records, unexpected error: "+err.Error(),
		)
		return
	}

	var returnedRecord DNSRecord
	for _, r := range recordResp.Response.Records {
		if r.Name == record.Name && r.Type == record.Type && r.Content == record.Content && r.TTL == record.TTL {
			returnedRecord = r
		}
	}

	// Overwrite DNS record with refreshed state
	plan.ZoneID = types.StringValue(recordResp.Response.ZoneConfig.ID)
	plan.ID = types.StringValue(returnedRecord.ID)
	plan.Name = types.StringValue(returnedRecord.Name)
	plan.Type = types.StringValue(returnedRecord.Type)
	plan.Content = types.StringValue(returnedRecord.Content)
	plan.TTL = types.Int64Value(int64(returnedRecord.TTL))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *recordResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state recordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	recordReq := RecordsFindRequest{
		BaseRequest: &BaseRequest{},
		Filter: FilterOrChain{Filter: Filter{
			Field: "RecordId",
			Value: state.ID.ValueString(),
		}},
		Limit: 1,
		Page:  1,
	}

	// Get refreshed DNS record from hostingde
	recordResp, err := r.client.listRecords(recordReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading hosting.de DNS zone",
			"Could not read hosting.de DNS zone ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite DNS record with refreshed state
	state.ZoneID = types.StringValue(recordResp.Response.Data[0].ZoneID)
	state.ID = types.StringValue(recordResp.Response.Data[0].ID)
	state.Name = types.StringValue(recordResp.Response.Data[0].Name)
	state.Type = types.StringValue(recordResp.Response.Data[0].Type)
	state.Content = types.StringValue(recordResp.Response.Data[0].Content)
	state.TTL = types.Int64Value(int64(recordResp.Response.Data[0].TTL))

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *recordResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan recordResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	record := DNSRecord{
		Name:    plan.Name.ValueString(),
		ID:      plan.ID.ValueString(),
		ZoneID:  plan.ZoneID.ValueString(),
		Type:    plan.Type.ValueString(),
		Content: plan.Content.ValueString(),
		TTL:     int(plan.TTL.ValueInt64()),
	}

	recordReq := RecordsUpdateRequest{
		BaseRequest:     &BaseRequest{},
		ZoneConfigId:    plan.ZoneID.ValueString(),
		RecordsToModify: []DNSRecord{record},
	}

	recordResp, err := r.client.updateRecords(recordReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating records",
			"Could not update records, unexpected error: "+err.Error(),
		)
		return
	}

	var returnedRecord DNSRecord
	for _, r := range recordResp.Response.Records {
		if r.Name == record.Name && r.Type == record.Type && r.Content == record.Content && r.TTL == record.TTL {
			returnedRecord = r
		}
	}

	// Overwrite DNS record with refreshed state
	plan.ZoneID = types.StringValue(recordResp.Response.ZoneConfig.ID)
	plan.ID = types.StringValue(returnedRecord.ID)
	plan.Name = types.StringValue(returnedRecord.Name)
	plan.Type = types.StringValue(returnedRecord.Type)
	plan.Content = types.StringValue(returnedRecord.Content)
	plan.TTL = types.Int64Value(int64(returnedRecord.TTL))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *recordResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state recordResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	record := DNSRecord{
		ID:   state.ID.ValueString(),
		Name: state.Name.ValueString(),
		Type: state.Type.ValueString(),
	}

	recordReq := RecordsUpdateRequest{
		BaseRequest:     &BaseRequest{},
		ZoneConfigId:    state.ZoneID.ValueString(),
		RecordsToDelete: []DNSRecord{record},
	}

	// Delete existing record
	_, err := r.client.updateRecords(recordReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Record",
			"Could not delete record, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *recordResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *recordResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
