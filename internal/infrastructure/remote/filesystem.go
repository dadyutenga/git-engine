package remote

import (
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// FileSystem implements basic remote file operations.
type FileSystem struct {
	Exec application.RemoteExecutor
}

// Exists reports if path exists remotely.
func (fs FileSystem) Exists(path string) (bool, error) {
	out, err := fs.Exec.Run("test -e " + shell.Escape(path) + " && echo exists || echo missing")
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "exists"), nil
}

// Mkdir ensures directory exists.
func (fs FileSystem) Mkdir(path string, recursive bool) error {
	cmd := "mkdir "
	if recursive {
		cmd += "-p "
	}
	_, err := fs.Exec.Run(cmd + shell.Escape(path))
	return err
}

// List returns entries within a directory.
func (fs FileSystem) List(path string) ([]string, error) {
	out, err := fs.Exec.Run("ls -1 " + shell.Escape(path))
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

var _ application.RemoteFileSystem = FileSystem{}
