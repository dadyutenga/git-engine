package detectors

import (
	"fmt"
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// DockerStrategy deploys docker-compose based projects.
type DockerStrategy struct {
	Exec application.RemoteExecutor
	FS   application.RemoteFileSystem
}

// Name returns the strategy identifier.
func (d DockerStrategy) Name() string { return "docker" }

// Detect checks for docker compose files.
func (d DockerStrategy) Detect(fs application.RemoteFileSystem, project domain.Project) (bool, error) {
	paths := []string{project.DeployDir + "/docker-compose.yml", project.DeployDir + "/compose.yml"}
	for _, p := range paths {
		ok, err := fs.Exists(p)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

// Deploy runs docker compose build+up.
func (d DockerStrategy) Deploy(project domain.Project, exec application.RemoteExecutor) error {
	_, err := exec.Run(fmt.Sprintf("cd %s && docker compose down && docker compose up -d --build", shell.Escape(project.DeployDir)))
	return err
}

// Restart restarts docker compose services.
func (d DockerStrategy) Restart(project domain.Project, exec application.RemoteExecutor) error {
	_, err := exec.Run(fmt.Sprintf("cd %s && docker compose restart", shell.Escape(project.DeployDir)))
	return err
}

// Status reports running state via docker compose.
func (d DockerStrategy) Status(project domain.Project, exec application.RemoteExecutor) (bool, error) {
	out, err := exec.Run(fmt.Sprintf("cd %s && docker compose ps --status running", shell.Escape(project.DeployDir)))
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "running"), nil
}

var _ application.DeploymentStrategy = DockerStrategy{}
