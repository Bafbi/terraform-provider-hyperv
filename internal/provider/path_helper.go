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

func PathDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		return true
	}
	oldNormalized := normalizeComparablePath(old)
	newNormalized := normalizeComparablePath(new)
	return strings.EqualFold(oldNormalized, newNormalized)
}

func PathDiffSuppressWithMachineName(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		return true
	}

	oldNormalized := normalizeComparablePath(old)
	newNormalized := normalizeComparablePath(new)

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
