# SSH Helper Package

This package provides SSH connectivity for the Terraform HyperV provider, replacing the WinRM-based communication with SSH.

## Overview

The `ssh-helper` package mirrors the interface of `winrm-helper`, providing the same operations over SSH instead of WinRM:

- Execute scripts (fire-and-forget or with result)
- Upload files and directories
- Check file/directory existence
- Delete files and directories

## Architecture

```
ssh-helper/
├── provider.go      - Interface definitions
├── client.go        - SSH client implementation
├── client_test.go   - Integration tests
└── README.md        - This file
```

## Usage

### Basic Configuration

```go
import (
    ssh_helper "github.com/taliesins/terraform-provider-hyperv/api/ssh-helper"
)

config := &ssh_helper.ClientConfig{
    Host:           "hyperv-host.example.com",
    Port:           22,
    User:           "admin",
    Password:       "password",  // or use PrivateKey/PrivateKeyPath
    Timeout:        30 * time.Second,
}

provider, err := ssh_helper.New(config)
if err != nil {
    log.Fatal(err)
}
```

### Authentication Methods

#### Password Authentication

```go
config := &ssh_helper.ClientConfig{
    Host:     "host",
    Port:     22,
    User:     "admin",
    Password: "password",
}
```

#### Private Key Authentication (from string)

```go
privateKeyContent := `-----BEGIN RSA PRIVATE KEY-----
...
-----END RSA PRIVATE KEY-----`

config := &ssh_helper.ClientConfig{
    Host:       "host",
    Port:       22,
    User:       "admin",
    PrivateKey: privateKeyContent,
}
```

#### Private Key Authentication (from file)

```go
config := &ssh_helper.ClientConfig{
    Host:           "host",
    Port:           22,
    User:           "admin",
    PrivateKeyPath: "/home/user/.ssh/id_rsa",
}
```

### Privilege Escalation

For operations requiring elevated privileges:

```go
config := &ssh_helper.ClientConfig{
    Host:            "host",
    Port:            22,
    User:            "admin",
    Password:        "password",
    ElevatedUser:    "root",
    ElevatedCommand: "sudo",  // or "doas", "su", etc.
}
```

### Running Scripts

#### Fire and Forget

```go
scriptTemplate := template.Must(template.New("setup").Parse(`
    mkdir -p {{.Path}}
    chmod 755 {{.Path}}
`))

args := struct {
    Path string
}{
    Path: "/opt/myapp",
}

err := config.RunFireAndForgetScript(ctx, scriptTemplate, args)
```

#### With JSON Result

```go
scriptTemplate := template.Must(template.New("query").Parse(`
    cat /etc/hostname | jq -R '{hostname: .}'
`))

var result struct {
    Hostname string `json:"hostname"`
}

err := config.RunScriptWithResult(ctx, scriptTemplate, nil, &result)
```

### File Operations

#### Upload File

```go
remotePath, err := config.UploadFile(ctx, "/local/file.txt", "/remote/destination/")
```

#### Upload Directory

```go
remoteRoot, files, err := config.UploadDirectory(
    ctx,
    "/local/directory",
    []string{"*.tmp", ".git"},  // exclusion patterns
)
```

#### Check Existence

```go
// Check file
exists, err := config.FileExists(ctx, "/path/to/file")

// Check directory
exists, err := config.DirectoryExists(ctx, "/path/to/directory")
```

#### Delete

```go
err := config.DeleteFileOrDirectory(ctx, "/path/to/remove")
```

## Interface Compatibility

The `ssh-helper` package implements the same `Client` interface as `winrm-helper`:

```go
type Client interface {
    RunFireAndForgetScript(ctx context.Context, script *template.Template, args interface{}) error
    RunScriptWithResult(ctx context.Context, script *template.Template, args interface{}, result interface{}) error
    UploadFile(ctx context.Context, filePath string, remoteFilePath string) (resolvedRemoteFilePath string, err error)
    UploadDirectory(ctx context.Context, rootPath string, excludeList []string) (remoteRootPath string, remoteAbsoluteFilePaths []string, err error)
    FileExists(ctx context.Context, remoteFilePath string) (exists bool, err error)
    DirectoryExists(ctx context.Context, remoteDirectoryPath string) (exists bool, err error)
    DeleteFileOrDirectory(ctx context.Context, remotePath string) error
}
```

This means you can swap between WinRM and SSH implementations with minimal code changes.

## Testing

### Running Tests

Integration tests require a real SSH server. Set environment variables:

```bash
export SSH_TEST_HOST="your-test-host"
export SSH_TEST_USER="testuser"
export SSH_TEST_PASSWORD="password"  # or SSH_TEST_KEY_PATH
export SSH_TEST_PORT="22"  # optional, defaults to 22

go test -v ./api/ssh-helper/
```

### Test Environment Variables

- `SSH_TEST_HOST` - Hostname or IP of SSH server (required)
- `SSH_TEST_USER` - SSH username (defaults to "root")
- `SSH_TEST_PASSWORD` - Password for authentication
- `SSH_TEST_KEY_PATH` - Path to private key file (alternative to password)
- `SSH_TEST_PORT` - SSH port (defaults to 22)

## Migration from WinRM

### Before (WinRM)

```go
import winrm_helper "github.com/taliesins/terraform-provider-hyperv/api/winrm-helper"

config := &winrm_helper.ClientConfig{
    WinRmClientPool:  pool,
    ElevatedUser:     "Administrator",
    ElevatedPassword: "password",
}
```

### After (SSH)

```go
import ssh_helper "github.com/taliesins/terraform-provider-hyperv/api/ssh-helper"

config := &ssh_helper.ClientConfig{
    Host:            "host",
    Port:            22,
    User:            "admin",
    Password:        "password",
    ElevatedUser:    "root",
    ElevatedCommand: "sudo",
}
```

## Differences from WinRM Helper

### Connection Management

- **WinRM**: Uses connection pooling with `WinRmClientPool`
- **SSH**: Creates new connection per operation (could be enhanced with pooling)

### Script Execution

- **WinRM**: Executes PowerShell scripts
- **SSH**: Executes bash/shell scripts

### File Transfer

- **WinRM**: Uses WinRM file transfer protocol
- **SSH**: Uses SCP protocol

### Privilege Escalation

- **WinRM**: Uses scheduled tasks
- **SSH**: Uses sudo/doas/su commands

## Security Considerations

### Current Implementation

⚠️ **Warning**: The current implementation uses `ssh.InsecureIgnoreHostKey()` which disables host key verification. This is suitable for development but **should be changed for production**.

### Production Recommendations

1. **Host Key Verification**: Implement proper host key checking
   ```go
   config.HostKeyCallback = ssh.FixedHostKey(hostKey)
   ```

2. **Private Key Authentication**: Prefer private keys over passwords

3. **SSH Agent**: Support SSH agent for key management

4. **Connection Pooling**: Implement connection pooling for better performance

5. **Timeout Configuration**: Set appropriate timeouts for your environment

## Future Enhancements

- [ ] Connection pooling (like WinRM helper)
- [ ] Proper host key verification
- [ ] SSH agent support
- [ ] SFTP for more efficient file transfers
- [ ] Port forwarding support
- [ ] Jump host / bastion support
- [ ] Better error handling and retry logic

## Dependencies

- `golang.org/x/crypto/ssh` - SSH client implementation

## Contributing

When adding new features, ensure:

1. Interface compatibility with `winrm-helper` is maintained
2. Tests are added for new functionality
3. Documentation is updated
4. Error handling is comprehensive
5. Logging follows the existing patterns

## License

Same as the main terraform-provider-hyperv project.
