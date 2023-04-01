package hostingde

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccZoneResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
resource "hostingde_zone" "test" {
  name = "example.test"
  type = "NATIVE"
  email = "hostmaster@example.test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify name attribute.
					resource.TestCheckResourceAttr("hostingde_zone.test", "name", "example.test"),
					// Verify type attribute.
					resource.TestCheckResourceAttr("hostingde_zone.test", "type", "NATIVE"),
					// Verify email attribute.
					resource.TestCheckResourceAttr("hostingde_zone.test", "email", "hostmaster@example.test"),
					// Verify dynamic values have any value set in the state.
					resource.TestCheckResourceAttrSet("hostingde_zone.test", "id"),
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
  name = "example.test"
  type = "NATIVE"
  email = "test@example.test"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					// Verify email attribute.
					resource.TestCheckResourceAttr("hostingde_zone.test", "email", "test@example.test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}
