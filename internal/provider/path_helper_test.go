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

func TestCaseInsensitiveDiffSuppress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		oldValue string
		newValue string
		want     bool
	}{
		{
			name:     "same case",
			oldValue: "On",
			newValue: "On",
			want:     true,
		},
		{
			name:     "different case suppressed",
			oldValue: "On",
			newValue: "on",
			want:     true,
		},
		{
			name:     "different case suppressed reverse",
			oldValue: "OFF",
			newValue: "off",
			want:     true,
		},
		{
			name:     "different values not suppressed",
			oldValue: "On",
			newValue: "Off",
			want:     false,
		},
		{
			name:     "mixed case suppressed",
			oldValue: "On",
			newValue: "ON",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := CaseInsensitiveDiffSuppress("field", tt.oldValue, tt.newValue, nil)
			if got != tt.want {
				t.Fatalf("CaseInsensitiveDiffSuppress(%q, %q) = %t, want %t", tt.oldValue, tt.newValue, got, tt.want)
			}
		})
	}
}

func TestZeroUuidDiffSuppress(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		oldValue string
		newValue string
		want     bool
	}{
		{
			name:     "same uuid",
			oldValue: "12345678-1234-1234-1234-123456789012",
			newValue: "12345678-1234-1234-1234-123456789012",
			want:     true,
		},
		{
			name:     "different uuids not suppressed",
			oldValue: "12345678-1234-1234-1234-123456789012",
			newValue: "87654321-4321-4321-4321-210987654321",
			want:     false,
		},
		{
			name:     "old zero uuid new empty suppressed",
			oldValue: "00000000-0000-0000-0000-000000000000",
			newValue: "",
			want:     true,
		},
		{
			name:     "old empty new zero uuid suppressed",
			oldValue: "",
			newValue: "00000000-0000-0000-0000-000000000000",
			want:     true,
		},
		{
			name:     "both zero uuid suppressed",
			oldValue: "00000000-0000-0000-0000-000000000000",
			newValue: "00000000-0000-0000-0000-000000000000",
			want:     true,
		},
		{
			name:     "old has value new empty not suppressed",
			oldValue: "12345678-1234-1234-1234-123456789012",
			newValue: "",
			want:     false,
		},
		{
			name:     "old zero uuid new has value not suppressed",
			oldValue: "00000000-0000-0000-0000-000000000000",
			newValue: "12345678-1234-1234-1234-123456789012",
			want:     false,
		},
		{
			name:     "case insensitive uuid suppressed",
			oldValue: "ABCDEF12-3456-7890-abcd-ef1234567890",
			newValue: "abcdef12-3456-7890-ABCD-EF1234567890",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ZeroUuidDiffSuppress("qos_policy_id", tt.oldValue, tt.newValue, nil)
			if got != tt.want {
				t.Fatalf("ZeroUuidDiffSuppress(%q, %q) = %t, want %t", tt.oldValue, tt.newValue, got, tt.want)
			}
		})
	}
}
