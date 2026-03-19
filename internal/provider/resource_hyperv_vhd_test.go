//go:build integration
// +build integration

package provider

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccHyperVVhd_lifecycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}

	tempDirectory, _ := filepath.Abs(".")
	vhdPath, _ := filepath.Abs(filepath.Join(tempDirectory, fmt.Sprintf("testacchypervvhd_%d.vhdx", randInt())))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccHyperVVhdConfigBasic(vhdPath, 8*1024*1024),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hyperv_vhd.this", "exists", "true"),
					resource.TestCheckResourceAttr("hyperv_vhd.this", "size", "8388608"),
				),
			},
			{
				Config: testAccHyperVVhdConfigBasic(vhdPath, 16*1024*1024),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("hyperv_vhd.this", "exists", "true"),
					resource.TestCheckResourceAttr("hyperv_vhd.this", "size", "16777216"),
				),
			},
		},
	})
}

func testAccHyperVVhdConfigBasic(path string, size int) string {
	return fmt.Sprintf(`
resource "hyperv_vhd" "this" {
  path = "%s"
  size = %d
}
`, escapeForHcl(path), size)
}
