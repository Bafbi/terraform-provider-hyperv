package provider

import "testing"

func TestProviderInternalValidate(t *testing.T) {
	t.Parallel()

	p := New("test-version", "test-commit")()
	if err := p.InternalValidate(); err != nil {
		t.Fatalf("provider InternalValidate failed: %v", err)
	}
}
