package ssh_helper

import (
	"context"
	"os"
	"testing"
	"text/template"
	"time"
)

// TestClientConfig_Basic tests basic SSH client configuration
func TestClientConfig_Basic(t *testing.T) {
	// Skip if SSH credentials are not provided
	host := os.Getenv("SSH_TEST_HOST")
	if host == "" {
		t.Skip("SSH_TEST_HOST not set, skipping integration test")
	}

	user := os.Getenv("SSH_TEST_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("SSH_TEST_PASSWORD")
	privateKeyPath := os.Getenv("SSH_TEST_KEY_PATH")

	if password == "" && privateKeyPath == "" {
		t.Skip("Neither SSH_TEST_PASSWORD nor SSH_TEST_KEY_PATH set, skipping test")
	}

	port := 22
	if portEnv := os.Getenv("SSH_TEST_PORT"); portEnv != "" {
		// Parse port if needed
		port = 22
	}

	config := &ClientConfig{
		Host:           host,
		Port:           port,
		User:           user,
		Password:       password,
		PrivateKeyPath: privateKeyPath,
		Timeout:        30 * time.Second,
	}

	// Test connection
	client, err := config.getSSHClient()
	if err != nil {
		t.Fatalf("Failed to create SSH client: %v", err)
	}
	defer client.Close()

	t.Log("Successfully connected via SSH")
}

// TestClientConfig_RunCommand tests running basic commands
func TestClientConfig_RunCommand(t *testing.T) {
	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()

	tests := []struct {
		name        string
		command     string
		expectError bool
	}{
		{
			name:        "Echo command",
			command:     "echo 'Hello, World!'",
			expectError: false,
		},
		{
			name:        "List directory",
			command:     "ls -la /tmp",
			expectError: false,
		},
		{
			name:        "Check hostname",
			command:     "hostname",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, exitCode, err := config.runCommand(ctx, tt.command)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			t.Logf("Command: %s\nExit Code: %d\nStdout: %s\nStderr: %s", tt.command, exitCode, stdout, stderr)
		})
	}
}

// TestClientConfig_FileExists tests file existence checking
func TestClientConfig_FileExists(t *testing.T) {
	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()

	// Test file that should exist
	exists, err := config.FileExists(ctx, "/etc/passwd")
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	if !exists {
		t.Error("Expected /etc/passwd to exist")
	}

	// Test file that should not exist
	exists, err = config.FileExists(ctx, "/tmp/nonexistent-file-12345")
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	if exists {
		t.Error("Expected /tmp/nonexistent-file-12345 to not exist")
	}
}

// TestClientConfig_DirectoryExists tests directory existence checking
func TestClientConfig_DirectoryExists(t *testing.T) {
	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()

	// Test directory that should exist
	exists, err := config.DirectoryExists(ctx, "/tmp")
	if err != nil {
		t.Fatalf("Failed to check directory existence: %v", err)
	}
	if !exists {
		t.Error("Expected /tmp to exist")
	}

	// Test directory that should not exist
	exists, err = config.DirectoryExists(ctx, "/nonexistent-directory-12345")
	if err != nil {
		t.Fatalf("Failed to check directory existence: %v", err)
	}
	if exists {
		t.Error("Expected /nonexistent-directory-12345 to not exist")
	}
}

// TestClientConfig_RunScriptWithResult tests script execution with JSON result
func TestClientConfig_RunScriptWithResult(t *testing.T) {
	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()

	// Create a script template that returns JSON
	scriptTemplate := template.Must(template.New("test").Parse(`
		echo '{"message": "{{.Message}}", "value": {{.Value}}}'
	`))

	args := struct {
		Message string
		Value   int
	}{
		Message: "test message",
		Value:   42,
	}

	var result struct {
		Message string `json:"message"`
		Value   int    `json:"value"`
	}

	err := config.RunScriptWithResult(ctx, scriptTemplate, args, &result)
	if err != nil {
		t.Fatalf("Failed to run script with result: %v", err)
	}

	if result.Message != args.Message {
		t.Errorf("Expected message %q, got %q", args.Message, result.Message)
	}

	if result.Value != args.Value {
		t.Errorf("Expected value %d, got %d", args.Value, result.Value)
	}
}

// TestClientConfig_RunFireAndForgetScript tests fire-and-forget script execution
func TestClientConfig_RunFireAndForgetScript(t *testing.T) {
	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()

	// Create a simple script template
	scriptTemplate := template.Must(template.New("test").Parse(`
		mkdir -p /tmp/hyperv-test-{{.TestID}}
		echo "Test content" > /tmp/hyperv-test-{{.TestID}}/test.txt
	`))

	args := struct {
		TestID string
	}{
		TestID: "12345",
	}

	err := config.RunFireAndForgetScript(ctx, scriptTemplate, args)
	if err != nil {
		t.Fatalf("Failed to run fire and forget script: %v", err)
	}

	// Verify the file was created
	exists, err := config.FileExists(ctx, "/tmp/hyperv-test-12345/test.txt")
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	if !exists {
		t.Error("Expected test file to be created")
	}

	// Cleanup
	err = config.DeleteFileOrDirectory(ctx, "/tmp/hyperv-test-12345")
	if err != nil {
		t.Logf("Warning: failed to cleanup test directory: %v", err)
	}
}

// getTestConfig creates a test configuration from environment variables
// Returns nil if required environment variables are not set
func getTestConfig(t *testing.T) *ClientConfig {
	host := os.Getenv("SSH_TEST_HOST")
	if host == "" {
		t.Skip("SSH_TEST_HOST not set, skipping integration test")
		return nil
	}

	user := os.Getenv("SSH_TEST_USER")
	if user == "" {
		user = "root"
	}

	password := os.Getenv("SSH_TEST_PASSWORD")
	privateKeyPath := os.Getenv("SSH_TEST_KEY_PATH")

	if password == "" && privateKeyPath == "" {
		t.Skip("Neither SSH_TEST_PASSWORD nor SSH_TEST_KEY_PATH set, skipping test")
		return nil
	}

	port := 22
	if portEnv := os.Getenv("SSH_TEST_PORT"); portEnv != "" {
		// Parse port if needed
		port = 22
	}

	return &ClientConfig{
		Host:           host,
		Port:           port,
		User:           user,
		Password:       password,
		PrivateKeyPath: privateKeyPath,
		Timeout:        30 * time.Second,
	}
}
