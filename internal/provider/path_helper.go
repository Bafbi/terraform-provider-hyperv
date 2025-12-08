package provider

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/taliesins/terraform-provider-hyperv/api"
)

func PathStateFunc(val interface{}) string {
	return api.NormalizePath(val.(string))
}

func PathDiffSuppress(k, old, new string, d *schema.ResourceData) bool {
	if new == "" {
		return true
	}
	oldNormalized := api.NormalizePath(old)
	newNormalized := api.NormalizePath(new)
	return strings.EqualFold(oldNormalized, newNormalized)
}
