package hostingde

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceRecord() *schema.Resource {
	return &schema.Resource{
		Create: resourceRecordCreate,
		Read:   resourceRecordRead,
		Update: resourceRecordUpdate,
		Delete: resourceRecordDelete,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeInt,
				Required: false,
				Optional: true,
				Default:  60,
			},
		},
	}
}

func resourceRecordCreate(d *schema.ResourceData, m interface{}) error {
	return resourceRecordCreateOrUpdate(d, m, "")
}

func resourceRecordRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	req := RecordsFindRequest{
		BaseRequest: &BaseRequest{},
		Filter: FilterOrChain{
			SubFilterConnective: "AND",
			SubFilter: []Filter{
				{
					Field: "ZoneConfigId",
					Value: d.Get("zone_id").(string),
				},
				{
					Field: "RecordId",
					Value: d.Id(),
				},
			},
		},
		Limit: 1,
		Page:  1,
	}
	record, err := c.getRecord(req)
	if err != nil {
		return err
	}
	_ = d.Set("name", record.Name)
	_ = d.Set("type", record.Type)
	_ = d.Set("content", record.Content)
	_ = d.Set("ttl", record.TTL)
	d.SetId(record.ID)
	return nil
}

func resourceRecordUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceRecordCreateOrUpdate(d, m, d.Id())
}

func resourceRecordDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	zoneConfig, err := getZoneConfig(d.Get("zone_id").(string), c)
	if err != nil {
		return err
	}

	req := ZoneUpdateRequest{
		BaseRequest:     &BaseRequest{},
		ZoneConfig:      *zoneConfig,
		RecordsToDelete: []DNSRecord{{ID: d.Id()}},
	}

	resp, err := c.updateZone(req)
	if err != nil {
		return err
	}

	for _, r := range resp.Response.Records {
		if r.ID == d.Id() {
			return fmt.Errorf("deleted record still in response from zone update")
		}
	}

	return nil
}

func resourceRecordCreateOrUpdate(d *schema.ResourceData, m interface{}, record_id string) error {
	c := m.(*Client)
	newRecord := DNSRecord{
		Name:    d.Get("name").(string),
		Type:    d.Get("type").(string),
		Content: d.Get("content").(string),
		TTL:     d.Get("ttl").(int),
	}

	zoneConfig, err := getZoneConfig(d.Get("zone_id").(string), c)
	if err != nil {
		d.SetId("")
		return nil
	}

	zoneUpdateReq := ZoneUpdateRequest{
		BaseRequest:     &BaseRequest{},
		ZoneConfig:      *zoneConfig,
		RecordsToAdd:    []DNSRecord{newRecord},
		RecordsToDelete: nil,
	}

	if record_id != "" {
		zoneUpdateReq.RecordsToDelete = []DNSRecord{{ID: record_id}}
	}

	resp, err := c.updateZone(zoneUpdateReq)
	if err != nil {
		return err
	}

	for _, r := range resp.Response.Records {
		if r.Name != newRecord.Name {
			continue
		}
		if r.Type != newRecord.Type {
			continue
		}
		if r.Content != newRecord.Content {
			continue
		}
		if r.TTL != newRecord.TTL {
			continue
		}
		d.SetId(r.ID)
		return resourceRecordRead(d, m)
	}

	return fmt.Errorf("response from server did not contain created record %w", resp.Response.Records)
}
