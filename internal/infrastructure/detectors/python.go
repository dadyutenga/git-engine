package detectors

import (
	"fmt"
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// PythonStrategy deploys Python services.
type PythonStrategy struct{}

// Name returns identifier.
func (PythonStrategy) Name() string { return "python" }

// Detect checks for requirements.txt or pyproject.
func (PythonStrategy) Detect(fs application.RemoteFileSystem, project domain.Project) (bool, error) {
	paths := []string{project.DeployDir + "/requirements.txt", project.DeployDir + "/pyproject.toml"}
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

// Deploy installs dependencies.
func (PythonStrategy) Deploy(project domain.Project, exec application.RemoteExecutor) error {
	cmd := fmt.Sprintf("cd %s && if [ -f requirements.txt ]; then pip install -r requirements.txt; fi", shell.Escape(project.DeployDir))
	_, err := exec.Run(cmd)
	return err
}

// Restart attempts to restart a systemd service matching project name.
func (PythonStrategy) Restart(project domain.Project, exec application.RemoteExecutor) error {
	_, err := exec.Run(fmt.Sprintf("systemctl restart %s || true", shell.Escape(project.Name)))
	return err
}

// Status checks systemd service state.
func (PythonStrategy) Status(project domain.Project, exec application.RemoteExecutor) (bool, error) {
	out, err := exec.Run(fmt.Sprintf("systemctl is-active %s || true", shell.Escape(project.Name)))
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "active"), nil
}

var _ application.DeploymentStrategy = PythonStrategy{}
