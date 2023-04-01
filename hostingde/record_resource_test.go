package hostingde

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccRecordResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "hostingde_zone" "test" {
  name = "example2.test"
  type = "NATIVE"
  email = "hostmaster@example2.test"
}
resource "hostingde_record" "test" {
  zone_id = hostingde_zone.test.id
  name = "test.example2.test"
  type = "CNAME"
  content = "www.example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name attribute.
					resource.TestCheckResourceAttr("hostingde_record.test", "name", "test.example2.test"),
					// Verify type attribute.
					resource.TestCheckResourceAttr("hostingde_record.test", "type", "CNAME"),
					// Verify email attribute.
					resource.TestCheckResourceAttr("hostingde_record.test", "content", "www.example.com"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("hostingde_record.test", "id"),
					resource.TestCheckResourceAttrSet("hostingde_record.test", "zone_id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "hostingde_zone.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: providerConfig + `
resource "hostingde_zone" "test" {
  name = "example2.test"
  type = "NATIVE"
  email = "hostmaster@example2.test"
}
resource "hostingde_record" "test" {
  zone_id = hostingde_zone.test.id
  name = "test.example2.test"
  type = "CNAME"
  content = "www2.example.com"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify content attribute.
					resource.TestCheckResourceAttr("hostingde_record.test", "content", "www2.example.com"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
