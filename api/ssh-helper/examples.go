package ssh_helper

// Example usage patterns for SSH helper

/*
Basic SSH Connection with Password:

	config := &ClientConfig{
		Host:     "192.168.1.100",
		Port:     22,
		User:     "admin",
		Password: "password",
		Timeout:  30 * time.Second,
	}

	provider, err := New(config)
	if err != nil {
		return err
	}

Basic SSH Connection with Private Key:

	config := &ClientConfig{
		Host:           "192.168.1.100",
		Port:           22,
		User:           "admin",
		PrivateKeyPath: "/home/user/.ssh/id_rsa",
		Timeout:        30 * time.Second,
	}

	provider, err := New(config)
	if err != nil {
		return err
	}

SSH with Privilege Escalation (sudo):

	config := &ClientConfig{
		Host:            "192.168.1.100",
		Port:            22,
		User:            "admin",
		Password:        "password",
		ElevatedUser:    "root",
		ElevatedCommand: "sudo",
		Timeout:         30 * time.Second,
	}

	provider, err := New(config)
	if err != nil {
		return err
	}

Running PowerShell Commands via SSH (Linux host with PowerShell):

	scriptTemplate := template.Must(template.New("ps").Parse(`
		pwsh -Command "Get-VM -Name {{.VMName}} | ConvertTo-Json"
	`))

	args := struct {
		VMName string
	}{
		VMName: "test-vm",
	}

	var result map[string]interface{}
	err := config.RunScriptWithResult(ctx, scriptTemplate, args, &result)

Uploading and Executing PowerShell Scripts:

	// Upload the script
	remotePath, err := config.UploadFile(ctx,
		"/local/scripts/setup.ps1",
		"/tmp/setup.ps1",
	)
	if err != nil {
		return err
	}

	// Execute it
	scriptTemplate := template.Must(template.New("exec").Parse(`
		pwsh -File {{.ScriptPath}}
	`))

	args := struct {
		ScriptPath string
	}{
		ScriptPath: remotePath,
	}

	err = config.RunFireAndForgetScript(ctx, scriptTemplate, args)

File Transfer Operations:

	// Upload single file
	remotePath, err := config.UploadFile(ctx,
		"/local/config.xml",
		"/etc/hyperv/config.xml",
	)

	// Upload directory
	remoteRoot, files, err := config.UploadDirectory(ctx,
		"/local/scripts",
		[]string{"*.tmp", "*.log", ".git"},
	)

	// Check if file exists
	exists, err := config.FileExists(ctx, "/etc/hyperv/config.xml")

	// Check if directory exists
	exists, err := config.DirectoryExists(ctx, "/etc/hyperv")

	// Delete file or directory
	err = config.DeleteFileOrDirectory(ctx, "/tmp/old-configs")

Environment Variables:

	config := &ClientConfig{
		Host:     "192.168.1.100",
		Port:     22,
		User:     "admin",
		Password: "password",
		Vars:     "export HYPERV_HOST=localhost; export HYPERV_PORT=5985",
	}

Integration with Existing Provider Code:

	// The Provider wrapper makes it compatible with api.Client interface
	provider, err := ssh_helper.New(config)
	if err != nil {
		return err
	}

	// Can be used anywhere api.Client is expected
	var apiClient api.Client = provider.Client

	// Now use with existing HyperV API functions
	vmClient := hyperv_api.NewVmClient(apiClient)
*/
