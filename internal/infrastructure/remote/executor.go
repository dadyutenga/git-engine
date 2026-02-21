package remote

import (
	"io"

	"github.com/dadyutenga/git-engine/internal/application"
	sshclient "github.com/dadyutenga/git-engine/internal/infrastructure/ssh"
)

// Executor implements application.RemoteExecutor using SSH.
type Executor struct {
	Client *sshclient.Client
}

// Run executes a remote command.
func (e Executor) Run(command string) (string, error) {
	return e.Client.Run(command)
}

// RunStream streams command output to the provided writer.
func (e Executor) RunStream(command string, writer io.Writer) error {
	return e.Client.RunStream(command, writer)
}

var _ application.RemoteExecutor = Executor{}
