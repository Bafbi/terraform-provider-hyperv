package provider

import (
	"context"
	"strings"
	"testing"
)

type destinationPathValidationClient struct {
	exists bool
	err    error
	path   string
}

func (c *destinationPathValidationClient) RemoteDirectoryExists(ctx context.Context, path string) (bool, error) {
	c.path = path
	return c.exists, c.err
}

func TestDestinationDriveRoot(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string
		path string
		want string
	}{
		{name: "Windows drive root", path: `V:\images\debian.iso`, want: `V:\`},
		{name: "Windows drive slash", path: `v:/images/debian.iso`, want: `V:\`},
		{name: "Quoted path", path: `'D:/iso/file.iso'`, want: `D:\`},
		{name: "Env path", path: `$env:TEMP\file.iso`, want: ""},
		{name: "UNC path", path: `\\server\share\file.iso`, want: ""},
		{name: "Empty", path: "", want: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := destinationDriveRoot(tc.path)
			if got != tc.want {
				t.Fatalf("destinationDriveRoot(%q) = %q, want %q", tc.path, got, tc.want)
			}
		})
	}
}

func TestValidateDestinationDriveExists(t *testing.T) {
	t.Parallel()

	t.Run("Skips non-drive paths", func(t *testing.T) {
		t.Parallel()

		client := &destinationPathValidationClient{}
		err := validateDestinationDriveExists(context.Background(), client, `$env:TEMP\debian.iso`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.path != "" {
			t.Fatalf("expected no directory check call, got path %q", client.path)
		}
	})

	t.Run("Fails when drive missing", func(t *testing.T) {
		t.Parallel()

		client := &destinationPathValidationClient{exists: false}
		err := validateDestinationDriveExists(context.Background(), client, `V:/images/debian.iso`)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), `destination drive "V:" does not exist`) {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.path != `V:\` {
			t.Fatalf("expected drive root check path %q, got %q", `V:\`, client.path)
		}
	})

	t.Run("Passes when drive exists", func(t *testing.T) {
		t.Parallel()

		client := &destinationPathValidationClient{exists: true}
		err := validateDestinationDriveExists(context.Background(), client, `D:\images\debian.iso`)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client.path != `D:\` {
			t.Fatalf("expected drive root check path %q, got %q", `D:\`, client.path)
		}
	})
}
