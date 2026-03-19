package provider

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func TestAllowedIsoVolumeName(t *testing.T) {
	t.Parallel()

	validator := AllowedIsoVolumeName()

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{name: "valid uppercase and underscore", input: "MY_ISO_01", wantError: false},
		{name: "valid empty string", input: "", wantError: false},
		{name: "too long", input: "ABCDEFGHIJKLMNOP", wantError: true},
		{name: "lowercase not allowed", input: "iso_name", wantError: true},
		{name: "dot not allowed", input: "ISO.NAME", wantError: true},
		{name: "wrong type", input: 42, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := validator(tt.input, cty.Path{})
			if hasErrorDiag(diags) != tt.wantError {
				t.Fatalf("AllowedIsoVolumeName(%v) error=%t, want %t", tt.input, hasErrorDiag(diags), tt.wantError)
			}
		})
	}
}

func TestStringKeyInMap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		validator func() func(interface{}, cty.Path) diag.Diagnostics
		input     interface{}
		wantError bool
	}{
		{
			name: "exact key match",
			validator: func() func(interface{}, cty.Path) diag.Diagnostics {
				return StringKeyInMap(map[string]int{"alpha": 1, "beta": 2}, false)
			},
			input:     "alpha",
			wantError: false,
		},
		{
			name: "case insensitive key match",
			validator: func() func(interface{}, cty.Path) diag.Diagnostics {
				return StringKeyInMap(map[string]int{"alpha": 1}, true)
			},
			input:     "ALPHA",
			wantError: false,
		},
		{
			name: "missing key",
			validator: func() func(interface{}, cty.Path) diag.Diagnostics {
				return StringKeyInMap(map[string]int{"alpha": 1}, false)
			},
			input:     "ALPHA",
			wantError: true,
		},
		{
			name: "invalid map argument",
			validator: func() func(interface{}, cty.Path) diag.Diagnostics {
				return StringKeyInMap([]string{"alpha"}, false)
			},
			input:     "alpha",
			wantError: true,
		},
		{
			name: "wrong input type",
			validator: func() func(interface{}, cty.Path) diag.Diagnostics {
				return StringKeyInMap(map[string]int{"alpha": 1}, false)
			},
			input:     100,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := tt.validator()(tt.input, cty.Path{})
			if hasErrorDiag(diags) != tt.wantError {
				t.Fatalf("StringKeyInMap input=%v error=%t, want %t", tt.input, hasErrorDiag(diags), tt.wantError)
			}
		})
	}
}

func TestIntInSlice(t *testing.T) {
	t.Parallel()

	validator := IntInSlice([]int{1, 3, 5})

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{name: "value in slice", input: 3, wantError: false},
		{name: "value not in slice", input: 2, wantError: true},
		{name: "wrong type", input: "3", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := validator(tt.input, cty.Path{})
			if hasErrorDiag(diags) != tt.wantError {
				t.Fatalf("IntInSlice(%v) error=%t, want %t", tt.input, hasErrorDiag(diags), tt.wantError)
			}
		})
	}
}

func TestIntBetween(t *testing.T) {
	t.Parallel()

	validator := IntBetween(10, 20)

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{name: "minimum", input: 10, wantError: false},
		{name: "maximum", input: 20, wantError: false},
		{name: "inside range", input: 15, wantError: false},
		{name: "below range", input: 9, wantError: true},
		{name: "above range", input: 21, wantError: true},
		{name: "wrong type", input: "15", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := validator(tt.input, cty.Path{})
			if hasErrorDiag(diags) != tt.wantError {
				t.Fatalf("IntBetween(%v) error=%t, want %t", tt.input, hasErrorDiag(diags), tt.wantError)
			}
		})
	}
}

func TestValueOrIntBetween(t *testing.T) {
	t.Parallel()

	validator := ValueOrIntBetween(0, 10, 20)

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{name: "special value allowed", input: 0, wantError: false},
		{name: "inside range", input: 15, wantError: false},
		{name: "minimum in range", input: 10, wantError: false},
		{name: "maximum in range", input: 20, wantError: false},
		{name: "outside range", input: 9, wantError: true},
		{name: "wrong type", input: "10", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := validator(tt.input, cty.Path{})
			if hasErrorDiag(diags) != tt.wantError {
				t.Fatalf("ValueOrIntBetween(%v) error=%t, want %t", tt.input, hasErrorDiag(diags), tt.wantError)
			}
		})
	}
}

func TestIsDivisibleBy(t *testing.T) {
	t.Parallel()

	validator := IsDivisibleBy(4096)

	tests := []struct {
		name      string
		input     interface{}
		wantError bool
	}{
		{name: "perfectly divisible", input: 8192, wantError: false},
		{name: "not divisible", input: 8193, wantError: true},
		{name: "zero divisible", input: 0, wantError: false},
		{name: "wrong type", input: "8192", wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			diags := validator(tt.input, cty.Path{})
			if hasErrorDiag(diags) != tt.wantError {
				t.Fatalf("IsDivisibleBy(%v) error=%t, want %t", tt.input, hasErrorDiag(diags), tt.wantError)
			}
		})
	}
}

func hasErrorDiag(diags diag.Diagnostics) bool {
	for _, d := range diags {
		if d.Severity == diag.Error {
			return true
		}
	}

	return false
}
