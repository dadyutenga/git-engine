package remote

import (
	"fmt"
	"log"
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// LockManager uses remote lock files to serialize deployments.
type LockManager struct {
	Exec application.RemoteExecutor
}

// Acquire tries to acquire a lock for the project.
// If the existing lock is stale (older than 60 minutes), it is removed and
// a single retry is attempted.
func (l LockManager) Acquire(project domain.Project) (bool, error) {
	cmd := `sh -c 'if ( set -o noclobber; echo $$ > ` + shell.Escape(project.LockFile) + ` ) 2>/dev/null; then echo acquired; else echo busy; fi'`
	out, err := l.Exec.Run(cmd)
	if err != nil {
		return false, err
	}
	if strings.Contains(out, "acquired") {
		return true, nil
	}

	// Check if the existing lock is stale (older than 60 minutes).
	checkStale := fmt.Sprintf("find %s -mmin +60 2>/dev/null", shell.Escape(project.LockFile))
	staleOut, staleErr := l.Exec.Run(checkStale)
	if staleErr != nil {
		log.Printf("WARNING: failed to check lock staleness for %s: %v", project.LockFile, staleErr)
		return false, nil
	}
	if strings.TrimSpace(staleOut) == "" {
		return false, nil
	}

	// Lock is stale; remove it and retry once.
	_ = l.Release(project)
	out, err = l.Exec.Run(cmd)
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "acquired"), nil
}

// Release frees the lock file.
func (l LockManager) Release(project domain.Project) error {
	_, err := l.Exec.Run("rm -f " + shell.Escape(project.LockFile))
	return err
}

var _ application.LockManager = LockManager{}
