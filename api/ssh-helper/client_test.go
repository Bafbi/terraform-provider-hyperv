package ssh_helper

import (
	"context"
	"encoding/base64"
	"os"
	"strings"
	"testing"
	"text/template"
	"time"
	"unicode/utf16"
)

// TestClientConfig_Basic tests basic SSH client configuration
func TestClientConfig_Basic(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

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
			command:     listDirectoryCommand(config.IsWindows),
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
			t.Parallel()

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
	t.Parallel()

	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()
	existingFilePath, missingFilePath := fileTestPaths(config.IsWindows)

	// Test file that should exist
	exists, err := config.FileExists(ctx, existingFilePath)
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	if !exists {
		t.Errorf("Expected %s to exist", existingFilePath)
	}

	// Test file that should not exist
	exists, err = config.FileExists(ctx, missingFilePath)
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	if exists {
		t.Errorf("Expected %s to not exist", missingFilePath)
	}
}

// TestClientConfig_DirectoryExists tests directory existence checking
func TestClientConfig_DirectoryExists(t *testing.T) {
	t.Parallel()

	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()
	existingDirectoryPath, missingDirectoryPath := directoryTestPaths(config.IsWindows)

	// Test directory that should exist
	exists, err := config.DirectoryExists(ctx, existingDirectoryPath)
	if err != nil {
		t.Fatalf("Failed to check directory existence: %v", err)
	}
	if !exists {
		t.Errorf("Expected %s to exist", existingDirectoryPath)
	}

	// Test directory that should not exist
	exists, err = config.DirectoryExists(ctx, missingDirectoryPath)
	if err != nil {
		t.Fatalf("Failed to check directory existence: %v", err)
	}
	if exists {
		t.Errorf("Expected %s to not exist", missingDirectoryPath)
	}
}

// TestClientConfig_RunScriptWithResult tests script execution with JSON result
func TestClientConfig_RunScriptWithResult(t *testing.T) {
	t.Parallel()

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
	t.Parallel()

	config := getTestConfig(t)
	if config == nil {
		return
	}

	ctx := context.Background()
	baseDirectory := testBaseDirectory(config.IsWindows)
	testDirectory := baseDirectory + "/hyperv-test-12345"
	if config.IsWindows {
		testDirectory = baseDirectory + "\\hyperv-test-12345"
	}
	testFilePath := testDirectory + "/test.txt"
	if config.IsWindows {
		testFilePath = testDirectory + "\\test.txt"
	}

	// Create a simple script template
	scriptTemplate := template.Must(template.New("test").Parse(fireAndForgetTemplate(config.IsWindows)))

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
	exists, err := config.FileExists(ctx, testFilePath)
	if err != nil {
		t.Fatalf("Failed to check file existence: %v", err)
	}
	if !exists {
		t.Error("Expected test file to be created")
	}

	// Cleanup
	err = config.DeleteFileOrDirectory(ctx, testDirectory)
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

	config := &ClientConfig{
		Host:           host,
		Port:           port,
		User:           user,
		Password:       password,
		PrivateKeyPath: privateKeyPath,
		Timeout:        30 * time.Second,
	}

	config.IsWindows = detectWindowsHost(t, config)

	return config
}

func detectWindowsHost(t *testing.T, config *ClientConfig) bool {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	stdout, _, exitCode, err := config.runCommand(ctx, "powershell -NoProfile -NonInteractive -Command \"[System.Environment]::OSVersion.Platform\"")
	if err != nil || exitCode != 0 {
		t.Logf("Could not detect Windows host via PowerShell probe, defaulting to Unix mode: err=%v exitCode=%d", err, exitCode)
		return false
	}

	return strings.Contains(strings.ToLower(strings.TrimSpace(stdout)), "win32nt")
}

func listDirectoryCommand(isWindows bool) string {
	if isWindows {
		return "Get-ChildItem C:\\Windows\\Temp"
	}

	return "ls -la /tmp"
}

func fileTestPaths(isWindows bool) (string, string) {
	if isWindows {
		return "C:\\Windows\\System32\\drivers\\etc\\hosts", "C:\\Windows\\Temp\\nonexistent-file-12345"
	}

	return "/etc/passwd", "/tmp/nonexistent-file-12345"
}

func directoryTestPaths(isWindows bool) (string, string) {
	if isWindows {
		return "C:\\Windows\\Temp", "C:\\Windows\\Temp\\nonexistent-directory-12345"
	}

	return "/tmp", "/nonexistent-directory-12345"
}

func testBaseDirectory(isWindows bool) string {
	if isWindows {
		return "C:\\Windows\\Temp"
	}

	return "/tmp"
}

func fireAndForgetTemplate(isWindows bool) string {
	if isWindows {
		return `
		New-Item -ItemType Directory -Force -Path 'C:\\Windows\\Temp\\hyperv-test-{{.TestID}}' | Out-Null
		Set-Content -Path 'C:\\Windows\\Temp\\hyperv-test-{{.TestID}}\\test.txt' -Value 'Test content'
	`
	}

	return `
		mkdir -p /tmp/hyperv-test-{{.TestID}}
		echo "Test content" > /tmp/hyperv-test-{{.TestID}}/test.txt
	`
}

func TestWrapPowerShellEncodedCommand(t *testing.T) {
	command := "$ErrorActionPreference = 'Stop'\nWrite-Output '{\"ok\":true}'"
	wrapped := wrapPowerShellEncodedCommand(command)

	prefix := "powershell -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand "
	if !strings.HasPrefix(wrapped, prefix) {
		t.Fatalf("expected command prefix %q, got %q", prefix, wrapped)
	}

	encoded := strings.TrimPrefix(wrapped, prefix)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	if len(decoded)%2 != 0 {
		t.Fatalf("expected UTF-16LE byte sequence, got odd length %d", len(decoded))
	}

	utf16Data := make([]uint16, len(decoded)/2)
	for i := 0; i < len(utf16Data); i++ {
		utf16Data[i] = uint16(decoded[i*2]) | uint16(decoded[i*2+1])<<8
	}

	expectedCommand := "if (Test-Path variable:global:ProgressPreference) { $ProgressPreference = 'SilentlyContinue' }; " + command
	if got := string(utf16.Decode(utf16Data)); got != expectedCommand {
		t.Fatalf("expected decoded command %q, got %q", expectedCommand, got)
	}
}

func TestWrapPowerShellEncodedCommandAlreadyHasProgressPreference(t *testing.T) {
	command := "$ProgressPreference = 'SilentlyContinue'\n$ErrorActionPreference = 'Stop'\nWrite-Output 'test'"
	wrapped := wrapPowerShellEncodedCommand(command)

	prefix := "powershell -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand "
	if !strings.HasPrefix(wrapped, prefix) {
		t.Fatalf("expected command prefix %q, got %q", prefix, wrapped)
	}

	encoded := strings.TrimPrefix(wrapped, prefix)
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("failed to decode base64: %v", err)
	}

	utf16Data := make([]uint16, len(decoded)/2)
	for i := 0; i < len(utf16Data); i++ {
		utf16Data[i] = uint16(decoded[i*2]) | uint16(decoded[i*2+1])<<8
	}

	if got := string(utf16.Decode(utf16Data)); got != command {
		t.Fatalf("expected original command unchanged %q, got %q", command, got)
	}
}

func TestPrepareCommandWindowsWrapsPowerShell(t *testing.T) {
	config := &ClientConfig{IsWindows: true}

	prepared, err := config.prepareCommand("$x = 1\nWrite-Output $x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.HasPrefix(prepared, "powershell -NoProfile -NonInteractive -ExecutionPolicy Bypass -EncodedCommand ") {
		t.Fatalf("expected encoded powershell command, got %q", prepared)
	}
}

func TestPrepareCommandWindowsKeepsExplicitPowerShell(t *testing.T) {
	config := &ClientConfig{IsWindows: true}
	explicit := `powershell -Command "Write-Output 'ok'"`

	prepared, err := config.prepareCommand(explicit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prepared != explicit {
		t.Fatalf("expected explicit powershell command unchanged, got %q", prepared)
	}
}

func TestPrepareCommandWindowsKeepsQuotedPwshPath(t *testing.T) {
	config := &ClientConfig{IsWindows: true}
	explicit := `"C:\Program Files\PowerShell\7\pwsh.exe" -Command "Write-Output 'ok'"`

	prepared, err := config.prepareCommand(explicit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prepared != explicit {
		t.Fatalf("expected quoted pwsh path command unchanged, got %q", prepared)
	}
}

func TestPrepareCommandWindowsKeepsCallOperatorQuotedPwshPath(t *testing.T) {
	config := &ClientConfig{IsWindows: true}
	explicit := `& "C:\Program Files\PowerShell\7\pwsh.exe" -Command "Write-Output 'ok'"`

	prepared, err := config.prepareCommand(explicit)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if prepared != explicit {
		t.Fatalf("expected call-operator quoted pwsh path command unchanged, got %q", prepared)
	}
}
