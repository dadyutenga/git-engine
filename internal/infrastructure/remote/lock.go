package remote

import (
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
)

// LockManager uses remote lock files to serialize deployments.
type LockManager struct {
	Exec application.RemoteExecutor
}

// Acquire tries to acquire a lock for the project.
func (l LockManager) Acquire(project domain.Project) (bool, error) {
	cmd := `sh -c 'if ( set -o noclobber; echo $$ > ` + project.LockFile + ` ) 2>/dev/null; then echo acquired; else echo busy; fi'`
	out, err := l.Exec.Run(cmd)
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "acquired"), nil
}

// Release frees the lock file.
func (l LockManager) Release(project domain.Project) error {
	_, err := l.Exec.Run("rm -f " + project.LockFile)
	return err
}

var _ application.LockManager = LockManager{}
