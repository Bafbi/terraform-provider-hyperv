package api

import "testing"

func TestNormalizePath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty",
			in:   "",
			want: "",
		},
		{
			name: "windows separators",
			in:   `C:\HyperV\Virtual Hard Disks\disk.vhdx`,
			want: "C:/HyperV/Virtual Hard Disks/disk.vhdx",
		},
		{
			name: "already normalized",
			in:   "C:/HyperV/Virtual Hard Disks/disk.vhdx",
			want: "C:/HyperV/Virtual Hard Disks/disk.vhdx",
		},
		{
			name: "mixed separators",
			in:   `C:/HyperV\Virtual Hard Disks/disk.vhdx`,
			want: "C:/HyperV/Virtual Hard Disks/disk.vhdx",
		},
		{
			name: "unc path",
			in:   `\\server\share\disks\disk.vhdx`,
			want: "//server/share/disks/disk.vhdx",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := NormalizePath(tc.in)
			if got != tc.want {
				t.Fatalf("NormalizePath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestToWindowsPath(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty",
			in:   "",
			want: "",
		},
		{
			name: "unix separators",
			in:   "C:/HyperV/Virtual Hard Disks/disk.vhdx",
			want: `C:\HyperV\Virtual Hard Disks\disk.vhdx`,
		},
		{
			name: "already windows",
			in:   `C:\HyperV\Virtual Hard Disks\disk.vhdx`,
			want: `C:\HyperV\Virtual Hard Disks\disk.vhdx`,
		},
		{
			name: "mixed separators",
			in:   `C:\HyperV/Virtual Hard Disks/disk.vhdx`,
			want: `C:\HyperV\Virtual Hard Disks\disk.vhdx`,
		},
		{
			name: "network path",
			in:   "//server/share/disks/disk.vhdx",
			want: `\\server\share\disks\disk.vhdx`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := ToWindowsPath(tc.in)
			if got != tc.want {
				t.Fatalf("ToWindowsPath(%q) = %q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
