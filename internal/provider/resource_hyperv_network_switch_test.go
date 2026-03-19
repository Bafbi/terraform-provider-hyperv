//go:build integration
// +build integration

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHyperVNetworkSwitch_basic(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	name := fmt.Sprintf("acc-switch-%d", randInt())

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHyperVNetworkSwitchConfigBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hyperv_network_switch.this", "name", name),
					resource.TestCheckResourceAttr("hyperv_network_switch.this", "switch_type", "Private"),
				),
			},
		},
	})
}

func testAccHyperVNetworkSwitchConfigBasic(name string) string {
	return fmt.Sprintf(`
resource "hyperv_network_switch" "this" {
  name        = "%s"
  switch_type = "Private"
  allow_management_os = false
}
`, escapeForHcl(name))
}
