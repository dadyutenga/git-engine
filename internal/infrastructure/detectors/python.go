package detectors

import (
	"fmt"
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// PythonStrategy deploys Python services.
type PythonStrategy struct {
	Exec application.RemoteExecutor
}

// Name returns identifier.
func (PythonStrategy) Name() string { return "python" }

// Detect checks for requirements.txt or pyproject.toml using a single batched command.
func (p PythonStrategy) Detect(fs application.RemoteFileSystem, project domain.Project) (bool, error) {
	cmd := fmt.Sprintf("( [ -f %s ] || [ -f %s ] ) && echo found || echo missing",
		shell.Escape(project.DeployDir+"/requirements.txt"),
		shell.Escape(project.DeployDir+"/pyproject.toml"))
	out, err := p.Exec.Run(cmd)
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "found"), nil
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
