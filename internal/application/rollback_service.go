package application

import (
	"fmt"
	"sort"
	"time"

	"github.com/dadyutenga/git-engine/internal/domain"
)

// RollbackService restores applications from backups.
type RollbackService struct {
	Exec       RemoteExecutor
	FS         RemoteFileSystem
	Strategies []DeploymentStrategy
}

// Rollback restores the specified or latest backup.
func (s RollbackService) Rollback(projectName, backup string) (domain.RollbackResult, error) {
	project := domain.NewProject(projectName)
	now := time.Now()
	result := domain.RollbackResult{ProjectName: project.Name, Timestamp: now}

	files, err := s.FS.List(project.BackupDir)
	if err != nil {
		result.Message = "failed to list backups"
		return result, err
	}
	if len(files) == 0 {
		result.Message = "no backups available"
		return result, fmt.Errorf("no backups found for %s", project.Name)
	}

	sort.Strings(files)
	chosen := backup
	if chosen == "" {
		chosen = files[len(files)-1]
	}

	if _, err := s.Exec.Run(fmt.Sprintf("tar -xzf %s/%s -C %s", project.BackupDir, chosen, project.DeployDir)); err != nil {
		result.Message = "failed to restore backup"
		return result, err
	}

	var strategy DeploymentStrategy
	for _, st := range s.Strategies {
		ok, derr := st.Detect(s.FS, project)
		if derr != nil {
			continue
		}
		if ok {
			strategy = st
			break
		}
	}

	if strategy != nil {
		_ = strategy.Restart(project, s.Exec)
	}

	result.Success = true
	result.Restored = chosen
	result.Message = fmt.Sprintf("rollback complete using %s", chosen)
	return result, nil
}
