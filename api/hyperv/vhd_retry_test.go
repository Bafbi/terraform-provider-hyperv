package hyperv

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"
)

func TestIsVhdResourceBusyError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "nil error",
			err:  nil,
			want: false,
		},
		{
			name: "object in use message",
			err:  errors.New("Get-VHD failed: The operation cannot be performed while the object is in use."),
			want: true,
		},
		{
			name: "object in use code",
			err:  errors.New("FullyQualifiedErrorId : ObjectInUse,Microsoft.Vhd.PowerShell.Cmdlets.GetVHD"),
			want: true,
		},
		{
			name: "resource busy",
			err:  errors.New("CategoryInfo : ResourceBusy: (:) [Get-VHD], VirtualizationException"),
			want: true,
		},
		{
			name: "other error",
			err:  errors.New("Access denied"),
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := isVhdResourceBusyError(tc.err)
			if got != tc.want {
				t.Fatalf("isVhdResourceBusyError() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestRunVhdOperationWithRetrySucceedsAfterTransientBusyError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	attempts := 0

	err := runVhdOperationWithRetry(ctx, `C:\temp\disk.vhdx`, "GetVhd", time.Millisecond, 100*time.Millisecond, func() error {
		attempts++
		if attempts < 3 {
			return errors.New("ObjectInUse")
		}
		return nil
	})

	if err != nil {
		t.Fatalf("expected retry to eventually succeed, got error: %v", err)
	}

	if attempts != 3 {
		t.Fatalf("expected 3 attempts, got %d", attempts)
	}
}

func TestRunVhdOperationWithRetryReturnsNonTransientErrorImmediately(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	attempts := 0

	err := runVhdOperationWithRetry(ctx, `C:\temp\disk.vhdx`, "GetVhd", time.Millisecond, 100*time.Millisecond, func() error {
		attempts++
		return errors.New("access denied")
	})

	if err == nil {
		t.Fatal("expected error")
	}

	if attempts != 1 {
		t.Fatalf("expected one attempt for non-transient error, got %d", attempts)
	}
}

func TestRunVhdOperationWithRetryTimesOutOnPersistentBusyError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	attempts := 0

	err := runVhdOperationWithRetry(ctx, `C:\temp\disk.vhdx`, "GetVhd", time.Millisecond, 5*time.Millisecond, func() error {
		attempts++
		return errors.New("The operation cannot be performed while the object is in use")
	})

	if err == nil {
		t.Fatal("expected timeout error")
	}

	if !strings.Contains(err.Error(), "timed out") {
		t.Fatalf("expected timeout error message, got: %v", err)
	}

	if attempts < 2 {
		t.Fatalf("expected retry attempts, got %d", attempts)
	}
}

func TestRunVhdOperationWithRetryHonorsContextCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	attempts := 0

	err := runVhdOperationWithRetry(ctx, `C:\temp\disk.vhdx`, "GetVhd", 50*time.Millisecond, time.Second, func() error {
		attempts++
		cancel()
		return errors.New("ObjectInUse")
	})

	if err == nil {
		t.Fatal("expected cancellation error")
	}

	if !strings.Contains(err.Error(), "canceled") {
		t.Fatalf("expected cancellation error message, got: %v", err)
	}

	if attempts != 1 {
		t.Fatalf("expected one attempt before cancellation, got %d", attempts)
	}
}
