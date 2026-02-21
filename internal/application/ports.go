package application

import (
	"io"

	"github.com/dadyutenga/git-engine/internal/domain"
)

// RemoteExecutor abstracts remote command execution over SSH.
type RemoteExecutor interface {
	Run(command string) (string, error)
	RunStream(command string, writer io.Writer) error
}

// RemoteFileSystem offers simple remote file operations.
type RemoteFileSystem interface {
	Exists(path string) (bool, error)
	Mkdir(path string, recursive bool) error
	List(path string) ([]string, error)
}

// LockManager coordinates distributed deployment locks.
type LockManager interface {
	Acquire(project domain.Project) (bool, error)
	Release(project domain.Project) error
}

// DeploymentStrategy implements detection and deployment for a project type.
type DeploymentStrategy interface {
	Name() string
	Detect(fs RemoteFileSystem, project domain.Project) (bool, error)
	Deploy(project domain.Project, exec RemoteExecutor) error
	Restart(project domain.Project, exec RemoteExecutor) error
	Status(project domain.Project, exec RemoteExecutor) (bool, error)
}
