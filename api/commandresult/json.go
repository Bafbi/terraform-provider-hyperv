package commandresult

import (
	"encoding/json"
	"fmt"
	"strings"
)

// DecodeJSON validates command execution output and decodes JSON stdout into result.
func DecodeJSON(exitStatus int, stdout, stderr, command string, result interface{}) error {
	stdout = strings.TrimSpace(stdout)

	if exitStatus != 0 {
		return fmt.Errorf("exitStatus:%d\nstdOut:%s\nstdErr:%s\ncommand:%s", exitStatus, stdout, stderr, command)
	}

	if stdout == "" {
		return fmt.Errorf("empty stdout from remote command - exitStatus:%d\nstdOut:%s\nstdErr:%s\ncommand:%s", exitStatus, stdout, stderr, command)
	}

	if err := json.Unmarshal([]byte(stdout), result); err != nil {
		return fmt.Errorf("failed to unmarshal JSON result - exitStatus:%d\nstdOut:%s\nstdErr:%s\nerr:%s\ncommand:%s", exitStatus, stdout, stderr, err, command)
	}

	return nil
}
