package application

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/dadyutenga/git-engine/internal/domain"
)

// InitService provisions remote directories and bare repositories for a new project.
type InitService struct {
	Exec RemoteExecutor
	FS   RemoteFileSystem
}

// Init creates the remote scaffold and validates SSH connectivity.
func (s InitService) Init(projectName string) (domain.InitResult, error) {
	project := domain.NewProject(projectName)
	now := time.Now()

	uidOut, err := s.Exec.Run("id -u")
	if err != nil {
		return domain.InitResult{Project: project, Success: false, Message: "failed to verify remote user", Timestamp: now}, err
	}
	if strings.TrimSpace(uidOut) == "0" {
		return domain.InitResult{Project: project, Success: false, Message: "refusing to run as root user", Timestamp: now}, fmt.Errorf("remote user is root")
	}

	paths := []string{
		project.RepoPath,
		project.DeployDir,
		project.BackupDir,
		filepath.Dir(project.LogFile),
		filepath.Dir(project.LockFile),
	}

	for _, p := range paths {
		if err := s.FS.Mkdir(p, true); err != nil {
			return domain.InitResult{Project: project, Success: false, Message: fmt.Sprintf("unable to create %s", p), Timestamp: now}, err
		}
	}

	if _, err := s.Exec.Run(fmt.Sprintf("test -d %s || git init --bare %s", project.RepoPath, project.RepoPath)); err != nil {
		return domain.InitResult{Project: project, Success: false, Message: "failed to initialize bare repository", Timestamp: now}, err
	}

	message := fmt.Sprintf("project %s initialized at %s", project.Name, project.DeployDir)
	return domain.InitResult{Project: project, Success: true, Message: message, Timestamp: now}, nil
}
