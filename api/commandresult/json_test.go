package commandresult

import (
	"strings"
	"testing"
)

func TestDecodeJSONSuccess(t *testing.T) {
	t.Parallel()

	var result struct {
		Value string `json:"value"`
	}

	err := DecodeJSON(0, `{"value":"ok"}`, "", "test-command", &result)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Value != "ok" {
		t.Fatalf("expected value %q, got %q", "ok", result.Value)
	}
}

func TestDecodeJSONEmptyStdout(t *testing.T) {
	t.Parallel()

	var result map[string]interface{}

	err := DecodeJSON(0, "", "", "test-command", &result)
	if err == nil {
		t.Fatal("expected error for empty stdout")
	}

	if got := err.Error(); !strings.HasPrefix(got, "empty stdout") {
		t.Fatalf("expected empty stdout error, got %q", err.Error())
	}
}

func TestDecodeJSONExitStatusError(t *testing.T) {
	t.Parallel()

	var result map[string]interface{}

	err := DecodeJSON(1, `{"value":"ok"}`, "failed", "test-command", &result)
	if err == nil {
		t.Fatal("expected error for non-zero exit status")
	}

	if got := err.Error(); !strings.HasPrefix(got, "exitStatus:1") {
		t.Fatalf("expected exit status error, got %q", err.Error())
	}
}

func TestDecodeJSONInvalidJSON(t *testing.T) {
	t.Parallel()

	var result map[string]interface{}

	err := DecodeJSON(0, "not-json", "", "test-command", &result)
	if err == nil {
		t.Fatal("expected error for invalid json")
	}

	if got := err.Error(); !strings.HasPrefix(got, "failed to unmarshal JSON") {
		t.Fatalf("expected unmarshal error, got %q", err.Error())
	}
}
