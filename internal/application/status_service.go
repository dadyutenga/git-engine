package application

import (
	"time"

	"github.com/dadyutenga/git-engine/internal/domain"
)

// StatusService inspects deployment state.
type StatusService struct {
	Exec       RemoteExecutor
	FS         RemoteFileSystem
	Strategies []DeploymentStrategy
}

// Status returns project status details.
func (s StatusService) Status(projectName string) (domain.StatusResult, error) {
	project := domain.NewProject(projectName)
	now := time.Now()
	result := domain.StatusResult{ProjectName: project.Name, Timestamp: now}

	exists, err := s.FS.Exists(project.DeployDir)
	if err != nil {
		result.Message = "failed to check project path"
		return result, err
	}
	if !exists {
		result.Exists = false
		result.Message = "project not found"
		return result, nil
	}
	result.Exists = true

	for _, st := range s.Strategies {
		ok, derr := st.Detect(s.FS, project)
		if derr != nil || !ok {
			continue
		}
		running, serr := st.Status(project, s.Exec)
		result.Running = running
		result.Strategy = st.Name()
		result.Message = "status retrieved"
		return result, serr
	}

	result.Message = "unknown project type"
	return result, domain.ErrUnsupportedProject
}
