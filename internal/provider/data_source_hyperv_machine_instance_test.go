package provider

import (
	"testing"
)

// TestDataSourceHyperVMachineInstanceSchema_PathFieldsDiffSuppressFunc verifies that
// smart_paging_file_path and snapshot_file_location in the data source schema have
// a DiffSuppressFunc configured (added in this PR).
func TestDataSourceHyperVMachineInstanceSchema_PathFieldsDiffSuppressFunc(t *testing.T) {
	t.Parallel()

	resource := dataSourceHyperVMachineInstance()

	fields := []string{"smart_paging_file_path", "snapshot_file_location"}
	for _, field := range fields {
		t.Run(field, func(t *testing.T) {
			t.Parallel()

			s, ok := resource.Schema[field]
			if !ok {
				t.Fatalf("schema field %q not found in dataSourceHyperVMachineInstance", field)
			}
			if s.DiffSuppressFunc == nil {
				t.Errorf("schema field %q has nil DiffSuppressFunc; expected PathDiffSuppress to be set", field)
			}
		})
	}
}

// TestDataSourceHyperVMachineInstanceSchema_SmartPagingFilePathDiffSuppress verifies that
// the DiffSuppressFunc on smart_paging_file_path behaves like PathDiffSuppress.
func TestDataSourceHyperVMachineInstanceSchema_SmartPagingFilePathDiffSuppress(t *testing.T) {
	t.Parallel()

	s := dataSourceHyperVMachineInstance().Schema["smart_paging_file_path"]
	if s.DiffSuppressFunc == nil {
		t.Fatal("DiffSuppressFunc is nil for smart_paging_file_path")
	}

	tests := []struct {
		name     string
		old      string
		newValue string
		want     bool
	}{
		{
			name:     "same path with backslash vs forward slash is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `C:/ProgramData/Microsoft/Windows/Hyper-V`,
			want:     true,
		},
		{
			name:     "case-insensitive same path is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `c:\programdata\microsoft\windows\hyper-v`,
			want:     true,
		},
		{
			name:     "empty new value is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: "",
			want:     true,
		},
		{
			name:     "genuinely different path is not suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `D:\VMs\SmartPaging`,
			want:     false,
		},
		{
			name:     "default windows path vs default linux-style path is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `/ProgramData/Microsoft/Windows/Hyper-V`,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := s.DiffSuppressFunc("smart_paging_file_path", tt.old, tt.newValue, nil)
			if got != tt.want {
				t.Errorf("DiffSuppressFunc(%q, %q) = %t, want %t", tt.old, tt.newValue, got, tt.want)
			}
		})
	}
}

// TestDataSourceHyperVMachineInstanceSchema_SnapshotFileLocationDiffSuppress verifies that
// the DiffSuppressFunc on snapshot_file_location behaves like PathDiffSuppress.
func TestDataSourceHyperVMachineInstanceSchema_SnapshotFileLocationDiffSuppress(t *testing.T) {
	t.Parallel()

	s := dataSourceHyperVMachineInstance().Schema["snapshot_file_location"]
	if s.DiffSuppressFunc == nil {
		t.Fatal("DiffSuppressFunc is nil for snapshot_file_location")
	}

	tests := []struct {
		name     string
		old      string
		newValue string
		want     bool
	}{
		{
			name:     "same path with backslash vs forward slash is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `C:/ProgramData/Microsoft/Windows/Hyper-V`,
			want:     true,
		},
		{
			name:     "case-insensitive same path is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `c:\programdata\microsoft\windows\hyper-v`,
			want:     true,
		},
		{
			name:     "empty new value is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: "",
			want:     true,
		},
		{
			name:     "genuinely different path is not suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `D:\Snapshots\VM`,
			want:     false,
		},
		{
			name:     "default windows path vs default linux-style path is suppressed",
			old:      `C:\ProgramData\Microsoft\Windows\Hyper-V`,
			newValue: `/ProgramData/Microsoft/Windows/Hyper-V`,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := s.DiffSuppressFunc("snapshot_file_location", tt.old, tt.newValue, nil)
			if got != tt.want {
				t.Errorf("DiffSuppressFunc(%q, %q) = %t, want %t", tt.old, tt.newValue, got, tt.want)
			}
		})
	}
}

// TestDataSourceHyperVMachineInstanceSchema_PathFieldsAlsoHaveStateFunc verifies that
// smart_paging_file_path and snapshot_file_location also retain their StateFunc (pre-existing).
func TestDataSourceHyperVMachineInstanceSchema_PathFieldsAlsoHaveStateFunc(t *testing.T) {
	t.Parallel()

	resource := dataSourceHyperVMachineInstance()

	fields := []string{"smart_paging_file_path", "snapshot_file_location"}
	for _, field := range fields {
		t.Run(field, func(t *testing.T) {
			t.Parallel()

			s, ok := resource.Schema[field]
			if !ok {
				t.Fatalf("schema field %q not found", field)
			}
			if s.StateFunc == nil {
				t.Errorf("schema field %q has nil StateFunc; PathStateFunc should still be set", field)
			}
		})
	}
}
