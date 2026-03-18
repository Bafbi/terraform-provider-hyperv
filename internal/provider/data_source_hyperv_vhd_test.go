//go:build integration
// +build integration

package provider

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func TestHyperVDataSourceVhd(t *testing.T) {
	// Skip if -short flag exist
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	//tempDirectory := os.TempDir() uses short name ;<
	tempDirectory, _ := filepath.Abs(".")
	path, _ := filepath.Abs(filepath.Join(tempDirectory, fmt.Sprintf("testhypervdatasourcevhd_%d.vhdx", randInt())))

	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testHyperDataSourceVVhdConfig(path),
				Check: resource.ComposeTestCheckFunc(
					testCheckPathEquivalent("data.hyperv_vhd.this", "path", path),
				),
			},
		},
	})
}

func testCheckPathEquivalent(resourceName, attributeName, expectedPath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found in state: %s", resourceName)
		}

		actualPath, ok := resourceState.Primary.Attributes[attributeName]
		if !ok {
			return fmt.Errorf("attribute not found in state: %s.%s", resourceName, attributeName)
		}

		expectedNormalized := strings.ToLower(api.NormalizePath(expectedPath))
		actualNormalized := strings.ToLower(api.NormalizePath(actualPath))

		if expectedNormalized == actualNormalized {
			return nil
		}

		if len(actualNormalized) > 2 && actualNormalized[1] == ':' {
			withoutDrive := actualNormalized[2:]
			if !strings.HasPrefix(withoutDrive, "/") {
				withoutDrive = "/" + withoutDrive
			}

			if withoutDrive == expectedNormalized {
				return nil
			}
		}

		return fmt.Errorf("attribute '%s' expected %q or Windows-drive equivalent, got %q", attributeName, expectedNormalized, actualNormalized)
	}
}

func testHyperDataSourceVVhdConfig(path string) string {
	return fmt.Sprintf(`
resource "hyperv_vhd" "this" {
	path = "%s"
	size = 4001792
}

data "hyperv_vhd" "this" {
	path = hyperv_vhd.this.path
}
	`, escapeForHcl(path))
}
