package hashicups

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrderResourceCrash(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: providerConfig + `
					resource "hashicups_order" "test" {
					  myblock {
						optional = false
						optional_int = 10
					  }
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hashicups_order.test", "id"),
					resource.TestCheckResourceAttrSet("hashicups_order.test", "myblock.optional"),
				),
			},
			{
				Config: providerConfig + `
					resource "hashicups_order" "test" {

					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("hashicups_order.test", "id"),
					resource.TestCheckNoResourceAttr("hashicups_order.test", "myblock.test"),
				),
			},
		},
	})
}
