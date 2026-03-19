package powershell

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func TestWinPath(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "empty path",
			input: "",
			want:  "",
		},
		{
			name:  "slash conversion",
			input: "C:/Temp/file.ps1",
			want:  `C:\Temp\file.ps1`,
		},
		{
			name:  "quote path with spaces",
			input: "C:/Program Files/Test/script.ps1",
			want:  `'C:\Program Files\Test\script.ps1'`,
		},
		{
			name:  "trim existing quotes before wrapping",
			input: "\"C:/Program Files/Test/script.ps1\"",
			want:  `'C:\Program Files\Test\script.ps1'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := winPath(tt.input)
			if got != tt.want {
				t.Fatalf("winPath(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestGetFilesInDirectory(t *testing.T) {
	t.Parallel()

	root := t.TempDir()

	pathsToCreate := []string{
		filepath.Join(root, "keep.txt"),
		filepath.Join(root, "ignore.log"),
		filepath.Join(root, "nested", "keep.ps1"),
		filepath.Join(root, "nested", "skip.tmp"),
		filepath.Join(root, "excluded-dir", "inside.txt"),
	}

	for _, p := range pathsToCreate {
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatalf("MkdirAll(%q) failed: %v", p, err)
		}
		if err := os.WriteFile(p, []byte("x"), 0o600); err != nil {
			t.Fatalf("WriteFile(%q) failed: %v", p, err)
		}
	}

	t.Run("returns all files when exclude list empty", func(t *testing.T) {
		t.Parallel()

		got, err := getFilesInDirectory(root, nil)
		if err != nil {
			t.Fatalf("getFilesInDirectory returned error: %v", err)
		}

		relGot := make([]string, 0, len(got))
		for _, p := range got {
			rel, relErr := filepath.Rel(root, p)
			if relErr != nil {
				t.Fatalf("filepath.Rel failed: %v", relErr)
			}
			relGot = append(relGot, filepath.ToSlash(rel))
		}

		sort.Strings(relGot)
		want := []string{
			"excluded-dir/inside.txt",
			"ignore.log",
			"keep.txt",
			"nested/keep.ps1",
			"nested/skip.tmp",
		}

		if len(relGot) != len(want) {
			t.Fatalf("unexpected file count: got %d, want %d; got=%v", len(relGot), len(want), relGot)
		}

		for i := range want {
			if relGot[i] != want[i] {
				t.Fatalf("unexpected files: got=%v want=%v", relGot, want)
			}
		}
	})

	t.Run("excludes files and directories via regex", func(t *testing.T) {
		t.Parallel()

		excludes := []string{`ignore\.log$`, `\.tmp$`, `excluded-dir`}
		got, err := getFilesInDirectory(root, excludes)
		if err != nil {
			t.Fatalf("getFilesInDirectory returned error: %v", err)
		}

		relGot := make([]string, 0, len(got))
		for _, p := range got {
			rel, relErr := filepath.Rel(root, p)
			if relErr != nil {
				t.Fatalf("filepath.Rel failed: %v", relErr)
			}
			relGot = append(relGot, filepath.ToSlash(rel))
		}

		sort.Strings(relGot)
		want := []string{
			"keep.txt",
			"nested/keep.ps1",
		}

		if len(relGot) != len(want) {
			t.Fatalf("unexpected file count: got %d, want %d; got=%v", len(relGot), len(want), relGot)
		}

		for i := range want {
			if relGot[i] != want[i] {
				t.Fatalf("unexpected files: got=%v want=%v", relGot, want)
			}
		}
	})
}
