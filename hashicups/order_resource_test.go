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
						my_boolean = false
						my_default_boolean = true
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hashicups_order.test", "computed_value", "1"),
					resource.TestCheckResourceAttr("hashicups_order.test", "my_default_boolean", "true"),
				),
			},
			{
				Config: providerConfig + `
					resource "hashicups_order" "test" {
						my_boolean = false
					}
				`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("hashicups_order.test", "computed_value", "2"),
					resource.TestCheckResourceAttr("hashicups_order.test", "my_default_boolean", "false"),
				),
			},
		},
	})
}
