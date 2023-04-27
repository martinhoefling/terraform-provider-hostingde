package hostingde

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &zoneResource{}
	_ resource.ResourceWithConfigure   = &zoneResource{}
	_ resource.ResourceWithImportState = &zoneResource{}
)

// NewZoneResource is a helper function to simplify the provider implementation.
func NewZoneResource() resource.Resource {
	return &zoneResource{}
}

// zoneResource is the resource implementation.
type zoneResource struct {
	client *Client
}

// zoneResourceModel maps the ZoneConfig resource schema data.
type zoneResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Type         types.String `tfsdk:"type"`
	EMailAddress types.String `tfsdk:"email"`
}

// Metadata returns the resource type name.
func (r *zoneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_zone"
}

// Schema defines the schema for the resource.
func (r *zoneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Numeric identifier of the zone.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Description: "Domain name (top-level domain) of the zone.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The zone type. Valid types are NATIVE, MASTER, and SLAVE. Changing this forces re-creation of the zone.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"email": schema.StringAttribute{
				Description: "The hostmaster email address. Only relevant if the type is NATIVE or MASTER. If the field is left empty, the default is hostmaster@name.",
				Computed:    true,
				Required:    false,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Create a new resource
func (r *zoneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan zoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := plan.Name.ValueString()
	ztype := plan.Type.ValueString()
	if ztype == "" {
		ztype = "NATIVE"
	}
	email := plan.EMailAddress.ValueString()

	// Generate API request body from plan
	zoneReq := ZoneCreateRequest{
		BaseRequest:             &BaseRequest{},
		UseDefaultNameserverSet: true,
		ZoneConfig: ZoneConfig{
			Name:         name,
			Type:         ztype,
			EMailAddress: email,
		},
		Records: []DNSRecord{},
	}
	zone, err := r.client.createZone(zoneReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating zone",
			"Could not create zone, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(zone.Response.ZoneConfig.ID)
	plan.Name = types.StringValue(zone.Response.ZoneConfig.Name)
	plan.Type = types.StringValue(zone.Response.ZoneConfig.Type)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *zoneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state zoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneReq := ZonesFindRequest{
		BaseRequest: &BaseRequest{},
		Filter: FilterOrChain{Filter: Filter{
			Field: "ZoneConfigId",
			Value: state.ID.ValueString(),
		}},
		Limit: 1,
		Page:  1,
	}

	// Get refreshed zone value from hosting.de
	zone, err := r.client.listZones(zoneReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading hosting.de DNS zone",
			"Could not read hosting.de DNS zone ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	state.ID = types.StringValue(zone.Response.Data[0].ZoneConfig.ID)
	state.Name = types.StringValue(zone.Response.Data[0].ZoneConfig.Name)
	state.Type = types.StringValue(zone.Response.Data[0].ZoneConfig.Type)
	state.EMailAddress = types.StringValue(zone.Response.Data[0].ZoneConfig.EMailAddress)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *zoneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan zoneResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneFindReq := ZonesFindRequest{
		BaseRequest: &BaseRequest{},
		Filter: FilterOrChain{Filter: Filter{
			Field: "ZoneConfigId",
			Value: plan.ID.ValueString(),
		}},
		Limit: 1,
		Page:  1,
	}

	// Get refreshed zone value from hosting.de
	zoneFindResp, err := r.client.listZones(zoneFindReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading hosting.de DNS zone",
			"Could not read hosting.de DNS zone ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	zoneConfig := zoneFindResp.Response.Data[0].ZoneConfig
	zoneConfig.ID = plan.ID.ValueString()
	zoneConfig.Name = plan.Name.ValueString()
	zoneConfig.Type = plan.Type.ValueString()
	zoneConfig.EMailAddress = plan.EMailAddress.ValueString()

	// Generate API request body from plan
	zoneReq := ZoneUpdateRequest{
		BaseRequest: &BaseRequest{},
		ZoneConfig:  zoneConfig,
	}
	zone, err := r.client.updateZone(zoneReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating zone",
			"Could not update zone, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(zone.Response.ZoneConfig.ID)
	plan.Name = types.StringValue(zone.Response.ZoneConfig.Name)
	plan.Type = types.StringValue(zone.Response.ZoneConfig.Type)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *zoneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state zoneResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	zoneReq := ZoneDeleteRequest{
		BaseRequest:  &BaseRequest{},
		ZoneConfigId: state.ID.ValueString(),
	}

	// Delete existing zone
	_, err := r.client.deleteZone(zoneReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting hosting.de Zone",
			"Could not delete zone, unexpected error: "+err.Error(),
		)
		return
	}

	// Purge restorable zone
	_, purgeErr := r.client.purgeZone(zoneReq)
	if purgeErr != nil {
		resp.Diagnostics.AddError(
			"Error Purging hosting.de Zone",
			"Could not purge zone, unexpected error: "+err.Error(),
		)
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *zoneResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*Client)
}

func (r *zoneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
