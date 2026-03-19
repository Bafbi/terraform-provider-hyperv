//go:build integration
// +build integration

package provider

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary
	version string = "0.0.0"

	// goreleaser can also pass the specific commit if you want
	commit string = ""
)

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"hyperv": func() (*schema.Provider, error) {
		return New(version, commit)(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := New(version, commit)().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	t.Helper()

	if os.Getenv("TF_ACC") == "" {
		t.Fatalf("TF_ACC must be set for acceptance tests")
	}

	if strings.EqualFold(os.Getenv("HYPERV_SSH"), "true") {
		if !hasAnyEnv("HYPERV_SSH_HOST", "HYPERV_HOST") {
			t.Fatalf("acceptance tests require HYPERV_SSH_HOST or HYPERV_HOST when HYPERV_SSH=true")
		}

		if !hasAnyEnv("HYPERV_SSH_USER", "HYPERV_USER") {
			t.Fatalf("acceptance tests require HYPERV_SSH_USER or HYPERV_USER when HYPERV_SSH=true")
		}

		if !hasAnyEnv("HYPERV_SSH_PASSWORD", "HYPERV_PASSWORD", "HYPERV_SSH_PRIVATE_KEY", "HYPERV_SSH_PRIVATE_KEY_PATH") {
			t.Fatalf("acceptance tests require SSH authentication: set one of HYPERV_SSH_PASSWORD, HYPERV_PASSWORD, HYPERV_SSH_PRIVATE_KEY, or HYPERV_SSH_PRIVATE_KEY_PATH")
		}

		if keyPath := os.Getenv("HYPERV_SSH_PRIVATE_KEY_PATH"); keyPath != "" {
			expandedKeyPath := keyPath
			if strings.HasPrefix(keyPath, "~/") {
				home, err := os.UserHomeDir()
				if err != nil {
					t.Fatalf("failed to resolve home directory for HYPERV_SSH_PRIVATE_KEY_PATH: %v", err)
				}
				expandedKeyPath = filepath.Join(home, keyPath[2:])
			}

			if _, err := os.Stat(expandedKeyPath); err != nil {
				t.Fatalf("HYPERV_SSH_PRIVATE_KEY_PATH points to a missing file: %s (%v)", expandedKeyPath, err)
			}
		}

		return
	}

	if !hasAnyEnv("HYPERV_PASSWORD") {
		t.Fatalf("acceptance tests require HYPERV_PASSWORD when HYPERV_SSH is not enabled")
	}
}

func hasAnyEnv(keys ...string) bool {
	for _, key := range keys {
		if os.Getenv(key) != "" {
			return true
		}
	}

	return false
}

func escapeForHcl(value string) string {
	return strings.ReplaceAll(value, "\\", "\\\\")
}

func randInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 100
	max := 999
	return rand.Intn(max-min+1) + min
}
