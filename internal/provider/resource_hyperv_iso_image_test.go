package provider

import (
	"encoding/hex"
	"os"
	"path/filepath"
	"testing"
)

func TestComputeFileSHA256(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	path := filepath.Join(dir, "test.iso")

	content := []byte("hello iso content")
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	got, err := computeFileSHA256(path)
	if err != nil {
		t.Fatalf("computeFileSHA256 returned error: %v", err)
	}

	if len(got) != hex.EncodedLen(32) {
		t.Errorf("expected 64-char hex digest, got %d chars: %s", len(got), got)
	}

	// Same content should produce the same hash.
	got2, err := computeFileSHA256(path)
	if err != nil {
		t.Fatalf("second computeFileSHA256 returned error: %v", err)
	}
	if got != got2 {
		t.Errorf("hash is non-deterministic: %s != %s", got, got2)
	}

	// Different content should produce a different hash.
	if err := os.WriteFile(path, []byte("different content"), 0600); err != nil {
		t.Fatalf("failed to overwrite temp file: %v", err)
	}
	got3, err := computeFileSHA256(path)
	if err != nil {
		t.Fatalf("computeFileSHA256 on updated file returned error: %v", err)
	}
	if got == got3 {
		t.Errorf("expected different hash after file update, got same: %s", got)
	}
}

func TestComputeFileSHA256_NonExistent(t *testing.T) {
	t.Parallel()

	_, err := computeFileSHA256("/nonexistent/path/file.iso")
	if err == nil {
		t.Error("expected an error for non-existent file, got nil")
	}
}
