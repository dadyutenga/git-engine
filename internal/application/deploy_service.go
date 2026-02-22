package application

import (
	"fmt"
	"log"
	"time"

	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// DeployService orchestrates push deployments.
type DeployService struct {
	Exec       RemoteExecutor
	FS         RemoteFileSystem
	Lock       LockManager
	Strategies []DeploymentStrategy
	Branch     string
}

// Deploy executes a deployment pipeline for the given project.
func (s DeployService) Deploy(projectName string) (domain.DeploymentResult, error) {
	project := domain.NewProject(projectName)
	now := time.Now()
	result := domain.DeploymentResult{ProjectName: project.Name, Timestamp: now, LogFile: project.LogFile}

	exists, err := s.FS.Exists(project.DeployDir)
	if err != nil {
		result.Message = "failed to check project path"
		return result, err
	}
	if !exists {
		result.Message = "project not found"
		return result, domain.ErrProjectNotFound
	}

	acquired, err := s.Lock.Acquire(project)
	if err != nil {
		result.Message = "failed to acquire deployment lock"
		return result, err
	}
	if !acquired {
		result.Message = "deployment lock unavailable"
		return result, domain.ErrLockUnavailable
	}
	defer func() {
		if err := s.Lock.Release(project); err != nil {
			log.Printf("WARNING: failed to release lock for %s: %v", project.Name, err)
		}
	}()

	if err := s.FS.Mkdir(project.BackupDir, true); err != nil {
		result.Message = "failed to ensure backup directory"
		return result, err
	}

	backupName := fmt.Sprintf("%s/%s-%d.tgz", project.BackupDir, project.Name, now.Unix())
	if _, err := s.Exec.Run(fmt.Sprintf("tar -czf %s -C %s .", shell.Escape(backupName), shell.Escape(project.DeployDir))); err != nil {
		result.Message = "failed to create backup"
		return result, err
	}

	branch := s.Branch
	if branch == "" {
		branch = "main"
	}
	if _, err := s.Exec.Run(fmt.Sprintf("cd %s && git fetch origin %s && git reset --hard origin/%s", shell.Escape(project.DeployDir), shell.Escape(branch), shell.Escape(branch))); err != nil {
		result.Message = "failed to update sources"
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
	if strategy == nil {
		result.Message = "unsupported project type"
		return result, domain.ErrUnsupportedProject
	}

	if err := strategy.Deploy(project, s.Exec); err != nil {
		result.Message = fmt.Sprintf("%s deployment failed", strategy.Name())
		return result, err
	}

	result.Success = true
	result.Status = "deployed"
	result.Message = fmt.Sprintf("%s deployed with %s strategy", project.Name, strategy.Name())
	return result, nil
}
