//nolint:forcetypeassert // Terraform schema/path state funcs provide expected concrete types.
package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func PathStateFunc(val interface{}) string {
	return api.NormalizePath(val.(string))
}

func normalizeComparablePath(path string) string {
	normalized := api.NormalizePath(path)
	if len(normalized) > 2 && normalized[1] == ':' {
		normalized = normalized[2:]
		if !strings.HasPrefix(normalized, "/") {
			normalized = "/" + normalized
		}
	}

	return normalized
}

func PathDiffSuppress(k, old, newValue string, d *schema.ResourceData) bool {
	if newValue == "" {
		return true
	}
	oldNormalized := normalizeComparablePath(old)
	newNormalized := normalizeComparablePath(newValue)
	return strings.EqualFold(oldNormalized, newNormalized)
}

func PathDiffSuppressWithMachineName(k, old, newValue string, d *schema.ResourceData) bool {
	if newValue == "" {
		return true
	}

	oldNormalized := normalizeComparablePath(old)
	newNormalized := normalizeComparablePath(newValue)

	name := d.Get("name").(string)
	computedPath := newNormalized
	if !strings.HasSuffix(computedPath, "/") {
		computedPath += "/"
	}
	computedPath += name

	if strings.EqualFold(computedPath, oldNormalized) {
		return true
	}

	return strings.EqualFold(oldNormalized, newNormalized)
}

func CaseInsensitiveDiffSuppress(k, old, newValue string, d *schema.ResourceData) bool {
	return strings.EqualFold(old, newValue)
}

func ZeroUuidDiffSuppress(k, old, newValue string, d *schema.ResourceData) bool {
	zeroUuid := "00000000-0000-0000-0000-000000000000"
	if newValue == "" || strings.EqualFold(newValue, zeroUuid) {
		return strings.EqualFold(old, zeroUuid) || old == ""
	}
	return strings.EqualFold(old, newValue)
}
