package hyperv

import (
	"github.com/taliesins/terraform-provider-hyperv/api"
	ssh_helper "github.com/taliesins/terraform-provider-hyperv/api/ssh-helper"
	winrm_helper "github.com/taliesins/terraform-provider-hyperv/api/winrm-helper"
)

func New(clientConfig *ClientConfig) (*api.Provider, error) {
	return &api.Provider{
		Client: clientConfig,
	}, nil
}

// NewSSH creates a new provider using SSH as the transport
// This allows the same HyperV API to work over SSH instead of WinRM
func NewSSH(clientConfig *ClientConfig) (*api.Provider, error) {
	return &api.Provider{
		Client: clientConfig,
	}, nil
}

type ClientConfig struct {
	WinRmClient winrm_helper.Client
}

// ClientConfigSSH wraps an SSH client to be compatible with the HyperV API
type ClientConfigSSH struct {
	SSHClient ssh_helper.Client
}
