package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestPathStateFunc(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "windows separators",
			input: `C:\VMs\test.vhdx`,
			want:  "C:/VMs/test.vhdx",
		},
		{
			name:  "already normalized",
			input: "/var/lib/hyperv/test.vhdx",
			want:  "/var/lib/hyperv/test.vhdx",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := PathStateFunc(tt.input)
			if got != tt.want {
				t.Fatalf("PathStateFunc(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestPathDiffSuppress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		oldValue string
		newValue string
		want     bool
	}{
		{
			name:     "empty new value is suppressed",
			oldValue: `C:\VMs\disk.vhdx`,
			newValue: "",
			want:     true,
		},
		{
			name:     "case-insensitive and separator-insensitive equality",
			oldValue: `C:\VMs\Disk.vhdx`,
			newValue: "c:/vms/disk.vhdx",
			want:     true,
		},
		{
			name:     "drive letter and linux-style equivalent",
			oldValue: `C:\VMs\disk.vhdx`,
			newValue: "/VMs/disk.vhdx",
			want:     true,
		},
		{
			name:     "different path is not suppressed",
			oldValue: `C:\VMs\disk1.vhdx`,
			newValue: "c:/vms/disk2.vhdx",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := PathDiffSuppress("path", tt.oldValue, tt.newValue, nil)
			if got != tt.want {
				t.Fatalf("PathDiffSuppress(%q, %q) = %t, want %t", tt.oldValue, tt.newValue, got, tt.want)
			}
		})
	}
}

func TestPathDiffSuppressWithMachineName(t *testing.T) {
	t.Parallel()

	resourceSchema := map[string]*schema.Schema{
		"name": {
			Type:     schema.TypeString,
			Optional: true,
		},
	}

	tests := []struct {
		name         string
		resource     map[string]interface{}
		oldValue     string
		newValue     string
		wantSuppress bool
	}{
		{
			name: "computed path appends machine name",
			resource: map[string]interface{}{
				"name": "vm01",
			},
			oldValue:     `C:\VMs\vm01`,
			newValue:     `C:/VMs`,
			wantSuppress: true,
		},
		{
			name: "exact path equality without appending name",
			resource: map[string]interface{}{
				"name": "ignored",
			},
			oldValue:     `C:\VMs\existing.vhdx`,
			newValue:     `c:/vms/existing.vhdx`,
			wantSuppress: true,
		},
		{
			name: "different path is not suppressed",
			resource: map[string]interface{}{
				"name": "vm01",
			},
			oldValue:     `C:\VMs\vm02`,
			newValue:     `C:/VMs`,
			wantSuppress: false,
		},
		{
			name: "empty new value is suppressed",
			resource: map[string]interface{}{
				"name": "vm01",
			},
			oldValue:     `C:\VMs\vm02`,
			newValue:     "",
			wantSuppress: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := schema.TestResourceDataRaw(t, resourceSchema, tt.resource)

			got := PathDiffSuppressWithMachineName("path", tt.oldValue, tt.newValue, d)
			if got != tt.wantSuppress {
				t.Fatalf("PathDiffSuppressWithMachineName(%q, %q) = %t, want %t", tt.oldValue, tt.newValue, got, tt.wantSuppress)
			}
		})
	}
}
