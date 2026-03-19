package api

import "testing"

func TestDiffSuppressVmHardDiskPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		old      string
		newValue string
		want     bool
	}{
		{
			name:     "empty new value suppressed",
			old:      `C:\HyperV\Virtual Hard Disks\disk.vhdx`,
			newValue: "",
			want:     true,
		},
		{
			name:     "same path with different separators suppressed",
			old:      `C:\HyperV\Virtual Hard Disks\disk.vhdx`,
			newValue: "C:/HyperV/Virtual Hard Disks/disk.vhdx",
			want:     true,
		},
		{
			name:     "same path case-insensitive suppressed",
			old:      `C:\HYPERV\VIRTUAL HARD DISKS\DISK.VHDX`,
			newValue: "c:/hyperv/virtual hard disks/disk.vhdx",
			want:     true,
		},
		{
			name:     "snapshot avhdx to base vhdx suppressed",
			old:      `C:\vhdx\web_server_g2_B63C9D15-F9A3-4F63-A896-FFD80BC7754C.avhdx`,
			newValue: `C:\vhdx\web_server_g2.vhdx`,
			want:     true,
		},
		{
			name:     "snapshot suppression requires underscore prefix",
			old:      `C:\vhdx\web_server_g2B63C9D15-F9A3-4F63-A896-FFD80BC7754C.avhdx`,
			newValue: `C:\vhdx\web_server_g2.vhdx`,
			want:     false,
		},
		{
			name:     "different paths not suppressed",
			old:      `C:\HyperV\Virtual Hard Disks\disk1.vhdx`,
			newValue: `C:\HyperV\Virtual Hard Disks\disk2.vhdx`,
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DiffSuppressVmHardDiskPath("hard_disk_drives.0.path", tc.old, tc.newValue, nil)
			if got != tc.want {
				t.Fatalf("DiffSuppressVmHardDiskPath(old=%q, new=%q) = %v, want %v", tc.old, tc.newValue, got, tc.want)
			}
		})
	}
}

func TestDiffSuppressVmIntegrationServices(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		key      string
		old      string
		newValue string
		want     bool
	}{
		{
			name:     "element count key is always suppressed",
			key:      "integration_services.%",
			old:      "2",
			newValue: "3",
			want:     true,
		},
		{
			name:     "empty new value suppressed",
			key:      "integration_services.VSS",
			old:      "false",
			newValue: "",
			want:     true,
		},
		{
			name:     "matching bool values suppressed",
			key:      "integration_services.Shutdown",
			old:      "true",
			newValue: "TRUE",
			want:     true,
		},
		{
			name:     "different bool values not suppressed",
			key:      "integration_services.Heartbeat",
			old:      "false",
			newValue: "true",
			want:     false,
		},
		{
			name:     "new parse fallback uses default true",
			key:      "integration_services.VSS",
			old:      "true",
			newValue: "not-a-bool",
			want:     true,
		},
		{
			name:     "old parse fallback uses default false",
			key:      "integration_services.Guest Service Interface",
			old:      "not-a-bool",
			newValue: "false",
			want:     true,
		},
		{
			name:     "both parse fallback unknown key defaults false",
			key:      "integration_services.Unknown Service",
			old:      "bad-old",
			newValue: "bad-new",
			want:     true,
		},
		{
			name:     "parse fallback then mismatch not suppressed",
			key:      "integration_services.Time Synchronization",
			old:      "invalid",
			newValue: "false",
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DiffSuppressVmIntegrationServices(tc.key, tc.old, tc.newValue, nil)
			if got != tc.want {
				t.Fatalf("DiffSuppressVmIntegrationServices(key=%q, old=%q, new=%q) = %v, want %v", tc.key, tc.old, tc.newValue, got, tc.want)
			}
		})
	}
}

func TestDiffSuppressVmStaticMacAddress(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		old      string
		newValue string
		want     bool
	}{
		{
			name:     "empty new value suppressed",
			old:      "00155D012345",
			newValue: "",
			want:     true,
		},
		{
			name:     "equal values suppressed",
			old:      "00155D012345",
			newValue: "00155D012345",
			want:     true,
		},
		{
			name:     "different values not suppressed",
			old:      "00155D012345",
			newValue: "00155D543210",
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DiffSuppressVmStaticMacAddress("network_adaptors.0.static_mac_address", tc.old, tc.newValue, nil)
			if got != tc.want {
				t.Fatalf("DiffSuppressVmStaticMacAddress(old=%q, new=%q) = %v, want %v", tc.old, tc.newValue, got, tc.want)
			}
		})
	}
}

func TestDiffSuppressVmProcessorMaximumCountPerNumaNode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		old      string
		newValue string
		want     bool
	}{
		{
			name:     "zero new value suppressed",
			old:      "4",
			newValue: "0",
			want:     true,
		},
		{
			name:     "equal non-zero values suppressed",
			old:      "8",
			newValue: "8",
			want:     true,
		},
		{
			name:     "different non-zero values not suppressed",
			old:      "8",
			newValue: "16",
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DiffSuppressVmProcessorMaximumCountPerNumaNode("vm_processor.0.maximum_count_per_numa_node", tc.old, tc.newValue, nil)
			if got != tc.want {
				t.Fatalf("DiffSuppressVmProcessorMaximumCountPerNumaNode(old=%q, new=%q) = %v, want %v", tc.old, tc.newValue, got, tc.want)
			}
		})
	}
}

func TestDiffSuppressVmProcessorMaximumCountPerNumaSocket(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		old      string
		newValue string
		want     bool
	}{
		{
			name:     "zero new value suppressed",
			old:      "4",
			newValue: "0",
			want:     true,
		},
		{
			name:     "equal non-zero values suppressed",
			old:      "12",
			newValue: "12",
			want:     true,
		},
		{
			name:     "different non-zero values not suppressed",
			old:      "12",
			newValue: "24",
			want:     false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := DiffSuppressVmProcessorMaximumCountPerNumaSocket("vm_processor.0.maximum_count_per_numa_socket", tc.old, tc.newValue, nil)
			if got != tc.want {
				t.Fatalf("DiffSuppressVmProcessorMaximumCountPerNumaSocket(old=%q, new=%q) = %v, want %v", tc.old, tc.newValue, got, tc.want)
			}
		})
	}
}
