package hyperv

import (
	"context"
	"text/template"

	"github.com/taliesins/terraform-provider-hyperv/api"
)

func New(clientConfig *ClientConfig) (*api.Provider, error) {
	return &api.Provider{
		Client: clientConfig,
	}, nil
}

type ClientConfig struct {
	ScriptRunner ScriptRunner
}

// ScriptRunner is a transport-agnostic remote script/file execution interface.
// Implemented by both WinRM and SSH helpers.
type ScriptRunner interface {
	RunFireAndForgetScript(ctx context.Context, script *template.Template, args interface{}) error
	RunScriptWithResult(ctx context.Context, script *template.Template, args interface{}, result interface{}) (err error)
	UploadFile(ctx context.Context, filePath string, remoteFilePath string) (resolvedRemoteFilePath string, err error)
	UploadDirectory(ctx context.Context, rootPath string, excludeList []string) (remoteRootPath string, remoteAbsoluteFilePaths []string, err error)
	FileExists(ctx context.Context, remoteFilePath string) (exists bool, err error)
	DirectoryExists(ctx context.Context, remoteDirectoryPath string) (exists bool, err error)
	DeleteFileOrDirectory(ctx context.Context, remotePath string) (err error)
}
