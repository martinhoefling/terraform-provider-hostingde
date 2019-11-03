package hostingde

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceZone() *schema.Resource {
	return &schema.Resource{
		Create: resourceZoneCreate,
		Read:   resourceZoneRead,
		Update: resourceZoneUpdate,
		Delete: resourceZoneDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: false,
				Optional: true,
			},
		},
	}
}

func resourceZoneCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	name := d.Get("name").(string)
	ztype := d.Get("type").(string)
	if ztype == "" {
		ztype = "NATIVE"
	}

	req := ZoneCreateRequest{
		BaseRequest:             &BaseRequest{},
		UseDefaultNameserverSet: true,
		ZoneConfig: ZoneConfig{
			Name: name,
			Type: ztype,
		},
	}
	resp, err := c.createZone(req)
	if err != nil {
		return err
	}

	d.SetId(resp.Response.ZoneConfig.ID)
	return resourceZoneRead(d, m)
}

func resourceZoneRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	req := ZoneConfigsFindRequest{
		BaseRequest: &BaseRequest{},
		Filter: FilterOrChain{Filter: Filter{
			Field: "ZoneConfigId",
			Value: d.Id(),
		}},
		Limit: 1,
		Page:  1,
	}
	resp, err := c.getZone(req)
	if err != nil {
		d.SetId("")
		return nil
	}
	_ = d.Set("name", resp.Name)
	d.SetId(resp.ID)
	return nil
}

func resourceZoneUpdate(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)

	zoneConfig, err := getZoneConfig(d.Id(), c)
	if err != nil {
		d.SetId("")
		return nil
	}
	zoneConfig.Name = d.Get("name").(string)

	req := ZoneUpdateRequest{
		BaseRequest: &BaseRequest{},
		ZoneConfig:  *zoneConfig,
	}
	resp, err := c.updateZone(req)
	if err != nil {
		return err
	}

	d.SetId(resp.Response.ZoneConfig.ID)
	return resourceZoneRead(d, m)
}

func resourceZoneDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*Client)
	req := ZoneDeleteRequest{
		BaseRequest:  &BaseRequest{},
		ZoneConfigId: d.Id(),
	}
	_, err := c.deleteZone(req)

	if err != nil {
		return err
	}

	return nil
}

func getZoneConfig(zoneId string, c *Client) (*ZoneConfig, error) {
	zoneFindReq := ZoneConfigsFindRequest{
		BaseRequest: &BaseRequest{},
		Filter: FilterOrChain{Filter: Filter{
			Field: "ZoneConfigId",
			Value: zoneId,
		}},
		Limit: 1,
		Page:  1,
	}
	zoneConfig, err := c.getZone(zoneFindReq)
	return zoneConfig, err
}
