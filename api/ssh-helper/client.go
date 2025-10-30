package ssh_helper

// SSH Helper for executing commands on remote systems via SSH.
// Supports both Linux/Unix and Windows remote hosts.
//
// For Windows hosts (IsWindows=true):
// - Uses PowerShell commands (Test-Path, Remove-Item, New-Item, etc.)
// - Paths use Windows-style backslashes (C:\Temp\...)
//
// For Linux/Unix hosts (IsWindows=false):
// - Uses standard Unix commands (test, rm, mkdir, etc.)
// - Paths use forward slashes (/tmp/...)

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// New creates a new SSH provider
func New(clientConfig *ClientConfig) (*Provider, error) {
	return &Provider{
		Client: clientConfig,
	}, nil
}

// ClientConfig holds the SSH connection configuration
type ClientConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	PrivateKey      string
	PrivateKeyPath  string
	Timeout         time.Duration
	KeepAlive       time.Duration
	ElevatedUser    string
	ElevatedCommand string // Command to use for privilege escalation (e.g., "sudo", "doas")
	Vars            string // Environment variables to set
	IsWindows       bool   // True if remote host is Windows (uses PowerShell instead of bash)
}

// getSSHClient creates and returns an SSH client connection
func (c *ClientConfig) getSSHClient() (*ssh.Client, error) {
	var authMethods []ssh.AuthMethod

	// Add password authentication if provided
	if c.Password != "" {
		authMethods = append(authMethods, ssh.Password(c.Password))
	}

	// Add private key authentication if provided
	if c.PrivateKey != "" {
		signer, err := ssh.ParsePrivateKey([]byte(c.PrivateKey))
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	} else if c.PrivateKeyPath != "" {
		// Expand ~ to home directory
		keyPath := c.PrivateKeyPath
		if strings.HasPrefix(keyPath, "~/") {
			homeDir, err := os.UserHomeDir()
			if err != nil {
				return nil, fmt.Errorf("failed to get home directory: %w", err)
			}
			keyPath = filepath.Join(homeDir, keyPath[2:])
		}

		keyBytes, err := os.ReadFile(keyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read private key file: %w", err)
		}
		signer, err := ssh.ParsePrivateKey(keyBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to parse private key from file: %w", err)
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no authentication method provided (password or private key required)")
	}

	config := &ssh.ClientConfig{
		User:            c.User,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO: Make this configurable for production
		Timeout:         c.Timeout,
	}

	addr := fmt.Sprintf("%s:%d", c.Host, c.Port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH: %w", err)
	}

	return client, nil
}

// runCommand executes a command over SSH and returns the output
func (c *ClientConfig) runCommand(ctx context.Context, command string) (stdout, stderr string, exitCode int, err error) {
	client, err := c.getSSHClient()
	if err != nil {
		return "", "", -1, err
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", "", -1, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdoutBuf, stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	// Add environment variables if provided
	if c.Vars != "" {
		command = fmt.Sprintf("%s; %s", c.Vars, command)
	}

	// Use privilege escalation if configured
	if c.ElevatedUser != "" && c.ElevatedCommand != "" {
		command = fmt.Sprintf("%s -u %s bash -c '%s'", c.ElevatedCommand, c.ElevatedUser, strings.ReplaceAll(command, "'", "'\\''"))
	}

	log.Printf("[DEBUG] Executing SSH command: %s", command)

	err = session.Run(command)
	stdout = stdoutBuf.String()
	stderr = stderrBuf.String()

	if err != nil {
		if exitErr, ok := err.(*ssh.ExitError); ok {
			exitCode = exitErr.ExitStatus()
		} else {
			return stdout, stderr, -1, fmt.Errorf("command execution failed: %w", err)
		}
	}

	return stdout, stderr, exitCode, nil
}

// RunFireAndForgetScript executes a script without waiting for or processing results
func (c *ClientConfig) RunFireAndForgetScript(ctx context.Context, script *template.Template, args interface{}) error {
	var scriptRendered bytes.Buffer
	err := script.Execute(&scriptRendered, args)
	if err != nil {
		return fmt.Errorf("failed to render script template: %w", err)
	}

	command := scriptRendered.String()
	log.Printf("[DEBUG] Running fire and forget script:\n%s\n", command)

	_, stderr, exitCode, err := c.runCommand(ctx, command)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("command failed with exit code %d: %s", exitCode, stderr)
	}

	return nil
}

// RunScriptWithResult executes a script and unmarshals JSON output into result
func (c *ClientConfig) RunScriptWithResult(ctx context.Context, script *template.Template, args interface{}, result interface{}) error {
	var scriptRendered bytes.Buffer
	err := script.Execute(&scriptRendered, args)
	if err != nil {
		return fmt.Errorf("failed to render script template: %w", err)
	}

	command := scriptRendered.String()
	log.Printf("[DEBUG] Running script with result:\n%s\n", command)

	stdout, stderr, exitCode, err := c.runCommand(ctx, command)
	if err != nil {
		return err
	}

	stdout = strings.TrimSpace(stdout)

	if exitCode != 0 {
		return fmt.Errorf("exitStatus:%d\nstdOut:%s\nstdErr:%s\ncommand:%s", exitCode, stdout, stderr, command)
	}

	err = json.Unmarshal([]byte(stdout), &result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON result - exitStatus:%d\nstdOut:%s\nstdErr:%s\nerr:%s\ncommand:%s", exitCode, stdout, stderr, err, command)
	}

	return nil
}

// UploadFile uploads a local file to the remote system
// Tries SFTP first, falls back to writing via PowerShell/shell commands
func (c *ClientConfig) UploadFile(ctx context.Context, filePath string, remoteFilePath string) (string, error) {
	client, err := c.getSSHClient()
	if err != nil {
		return "", err
	}
	defer client.Close()

	log.Printf("[DEBUG] Uploading file %s to %s", filePath, remoteFilePath)

	// Read local file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read local file: %w", err)
	}

	// If remote path is a directory or not specified, append the filename
	if remoteFilePath == "" || strings.HasSuffix(remoteFilePath, "/") {
		remoteFilePath = filepath.Join(remoteFilePath, filepath.Base(filePath))
	}

	// Try SFTP first
	err = c.uploadViaSFTP(client, fileData, remoteFilePath)
	if err == nil {
		log.Printf("[DEBUG] Successfully uploaded file via SFTP to %s", remoteFilePath)
		return remoteFilePath, nil
	}

	log.Printf("[DEBUG] SFTP upload failed: %v, trying command-based upload", err)

	// Fallback to command-based upload
	err = c.uploadViaCommands(client, fileData, remoteFilePath)
	if err != nil {
		return "", fmt.Errorf("all upload methods failed: %w", err)
	}

	log.Printf("[DEBUG] Successfully uploaded file via commands to %s", remoteFilePath)
	return remoteFilePath, nil
}

// uploadViaSFTP uploads a file using SFTP protocol
func (c *ClientConfig) uploadViaSFTP(client *ssh.Client, fileData []byte, remoteFilePath string) error {
	// Create SFTP client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		return fmt.Errorf("failed to create SFTP client: %w", err)
	}
	defer sftpClient.Close()

	// Create remote directory if needed
	remoteDir := filepath.Dir(remoteFilePath)
	if remoteDir != "" && remoteDir != "." {
		err = sftpClient.MkdirAll(remoteDir)
		if err != nil {
			log.Printf("[DEBUG] Failed to create remote directory via SFTP: %v", err)
			// Continue anyway, the directory might exist
		}
	}

	// Create remote file
	remoteFile, err := sftpClient.Create(remoteFilePath)
	if err != nil {
		return fmt.Errorf("failed to create remote file: %w", err)
	}
	defer remoteFile.Close()

	// Write file content
	_, err = remoteFile.Write(fileData)
	if err != nil {
		return fmt.Errorf("failed to write file content: %w", err)
	}

	return nil
}

// uploadViaCommands uploads a file by writing it via shell commands
func (c *ClientConfig) uploadViaCommands(client *ssh.Client, fileData []byte, remoteFilePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	// Encode file data as base64 for safe transport
	encoded := base64.StdEncoding.EncodeToString(fileData)

	var command string
	if c.IsWindows {
		// Use PowerShell to decode base64 and write file
		// Create directory first
		remoteDir := filepath.Dir(remoteFilePath)
		if remoteDir != "" && remoteDir != "." {
			dirSession, _ := client.NewSession()
			if dirSession != nil {
				dirCmd := fmt.Sprintf("powershell -Command \"New-Item -ItemType Directory -Force -Path '%s' | Out-Null\"", remoteDir)
				dirSession.Run(dirCmd)
				dirSession.Close()
			}
		}

		// Write file via PowerShell
		command = fmt.Sprintf(
			"powershell -Command \"$bytes = [System.Convert]::FromBase64String('%s'); [System.IO.File]::WriteAllBytes('%s', $bytes)\"",
			encoded,
			remoteFilePath,
		)
	} else {
		// Use Unix commands
		remoteDir := filepath.Dir(remoteFilePath)
		if remoteDir != "" && remoteDir != "." {
			dirSession, _ := client.NewSession()
			if dirSession != nil {
				dirSession.Run(fmt.Sprintf("mkdir -p '%s'", remoteDir))
				dirSession.Close()
			}
		}

		// Write file via base64 decode
		command = fmt.Sprintf("echo '%s' | base64 -d > '%s'", encoded, remoteFilePath)
	}

	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("failed to execute upload command: %w", err)
	}

	return nil
}

// UploadDirectory uploads a local directory to the remote system
func (c *ClientConfig) UploadDirectory(ctx context.Context, rootPath string, excludeList []string) (remoteRootPath string, remoteAbsoluteFilePaths []string, err error) {
	log.Printf("[DEBUG] Uploading directory %s", rootPath)

	// Create a temporary remote directory
	if c.IsWindows {
		remoteRootPath = fmt.Sprintf("C:\\Temp\\hyperv-upload-%d", time.Now().Unix())
	} else {
		remoteRootPath = fmt.Sprintf("/tmp/hyperv-upload-%d", time.Now().Unix())
	}

	client, err := c.getSSHClient()
	if err != nil {
		return "", nil, err
	}
	defer client.Close()

	// Create remote directory
	session, err := client.NewSession()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create session: %w", err)
	}

	var mkdirCmd string
	if c.IsWindows {
		mkdirCmd = fmt.Sprintf("powershell -Command \"New-Item -ItemType Directory -Force -Path '%s'\"", remoteRootPath)
	} else {
		mkdirCmd = fmt.Sprintf("mkdir -p %s", remoteRootPath)
	}

	err = session.Run(mkdirCmd)
	session.Close()
	if err != nil {
		return "", nil, fmt.Errorf("failed to create remote directory: %w", err)
	}

	// Walk through local directory and upload files
	remoteAbsoluteFilePaths = []string{}
	err = filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file should be excluded
		relPath, _ := filepath.Rel(rootPath, path)
		for _, exclude := range excludeList {
			if matched, _ := filepath.Match(exclude, relPath); matched {
				log.Printf("[DEBUG] Skipping excluded file: %s", relPath)
				return nil
			}
		}

		// Upload file
		var remotePath string
		if c.IsWindows {
			// Convert to Windows path
			remotePath = filepath.Join(remoteRootPath, relPath)
			remotePath = filepath.ToSlash(remotePath)              // Convert to forward slashes for consistency
			remotePath = strings.ReplaceAll(remotePath, "/", "\\") // Then to backslashes for Windows
		} else {
			remotePath = filepath.Join(remoteRootPath, relPath)
		}

		var remoteDir string
		if c.IsWindows {
			remoteDir = filepath.Dir(remotePath)
		} else {
			remoteDir = filepath.Dir(remotePath)
		}

		// Create remote directory structure
		session, err := client.NewSession()
		if err != nil {
			return err
		}

		var mkdirCmd string
		if c.IsWindows {
			mkdirCmd = fmt.Sprintf("powershell -Command \"New-Item -ItemType Directory -Force -Path '%s'\"", remoteDir)
		} else {
			mkdirCmd = fmt.Sprintf("mkdir -p %s", remoteDir)
		}

		session.Run(mkdirCmd)
		session.Close()

		// Upload the file
		_, err = c.UploadFile(ctx, path, remotePath)
		if err != nil {
			return err
		}

		remoteAbsoluteFilePaths = append(remoteAbsoluteFilePaths, remotePath)
		return nil
	})

	if err != nil {
		return "", nil, fmt.Errorf("failed to upload directory: %w", err)
	}

	log.Printf("[DEBUG] Successfully uploaded directory to %s with %d files", remoteRootPath, len(remoteAbsoluteFilePaths))
	return remoteRootPath, remoteAbsoluteFilePaths, nil
}

// FileExists checks if a file exists on the remote system
func (c *ClientConfig) FileExists(ctx context.Context, remoteFilePath string) (bool, error) {
	log.Printf("[DEBUG] Checking if file exists: %s", remoteFilePath)

	var command string
	if c.IsWindows {
		command = fmt.Sprintf("powershell -Command \"Test-Path -Path '%s' -PathType Leaf\"", remoteFilePath)
	} else {
		command = fmt.Sprintf("test -f '%s' && echo 'true' || echo 'false'", remoteFilePath)
	}

	stdout, _, exitCode, err := c.runCommand(ctx, command)
	if err != nil {
		return false, err
	}

	if exitCode != 0 {
		return false, nil
	}

	stdout = strings.TrimSpace(stdout)
	var exists bool
	if c.IsWindows {
		exists = strings.EqualFold(stdout, "true")
	} else {
		exists = stdout == "true"
	}

	if exists {
		log.Printf("[DEBUG] File exists: %s", remoteFilePath)
	} else {
		log.Printf("[DEBUG] File does not exist: %s", remoteFilePath)
	}

	return exists, nil
}

// DirectoryExists checks if a directory exists on the remote system
func (c *ClientConfig) DirectoryExists(ctx context.Context, remoteDirectoryPath string) (bool, error) {
	log.Printf("[DEBUG] Checking if directory exists: %s", remoteDirectoryPath)

	var command string
	if c.IsWindows {
		command = fmt.Sprintf("powershell -Command \"Test-Path -Path '%s' -PathType Container\"", remoteDirectoryPath)
	} else {
		command = fmt.Sprintf("test -d '%s' && echo 'true' || echo 'false'", remoteDirectoryPath)
	}

	stdout, _, exitCode, err := c.runCommand(ctx, command)
	if err != nil {
		return false, err
	}

	if exitCode != 0 {
		return false, nil
	}

	stdout = strings.TrimSpace(stdout)
	var exists bool
	if c.IsWindows {
		exists = strings.EqualFold(stdout, "true")
	} else {
		exists = stdout == "true"
	}

	if exists {
		log.Printf("[DEBUG] Directory exists: %s", remoteDirectoryPath)
	} else {
		log.Printf("[DEBUG] Directory does not exist: %s", remoteDirectoryPath)
	}

	return exists, nil
}

// DeleteFileOrDirectory removes a file or directory from the remote system
func (c *ClientConfig) DeleteFileOrDirectory(ctx context.Context, remotePath string) error {
	log.Printf("[DEBUG] Deleting file or directory: %s", remotePath)

	var command string
	if c.IsWindows {
		command = fmt.Sprintf("powershell -Command \"Remove-Item -Path '%s' -Recurse -Force -ErrorAction SilentlyContinue\"", remotePath)
	} else {
		command = fmt.Sprintf("rm -rf '%s'", remotePath)
	}

	_, stderr, exitCode, err := c.runCommand(ctx, command)
	if err != nil {
		return err
	}

	if exitCode != 0 {
		return fmt.Errorf("failed to delete %s: %s", remotePath, stderr)
	}

	log.Printf("[DEBUG] Successfully deleted: %s", remotePath)
	return nil
}
