package detectors

import (
	"fmt"
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// NodeStrategy deploys Node/pm2 projects.
type NodeStrategy struct{}

// Name returns strategy name.
func (NodeStrategy) Name() string { return "node" }

// Detect determines if package.json exists.
func (NodeStrategy) Detect(fs application.RemoteFileSystem, project domain.Project) (bool, error) {
	return fs.Exists(project.DeployDir + "/package.json")
}

// Deploy installs dependencies and restarts pm2.
func (NodeStrategy) Deploy(project domain.Project, exec application.RemoteExecutor) error {
	cmd := fmt.Sprintf("cd %s && npm install --production && (pm2 restart %s || pm2 start npm --name %s -- start)", shell.Escape(project.DeployDir), shell.Escape(project.Name), shell.Escape(project.Name))
	_, err := exec.Run(cmd)
	return err
}

// Restart restarts pm2 process.
func (NodeStrategy) Restart(project domain.Project, exec application.RemoteExecutor) error {
	_, err := exec.Run(fmt.Sprintf("pm2 restart %s", shell.Escape(project.Name)))
	return err
}

// Status returns pm2 process state.
func (NodeStrategy) Status(project domain.Project, exec application.RemoteExecutor) (bool, error) {
	out, err := exec.Run(fmt.Sprintf("pm2 describe %s", shell.Escape(project.Name)))
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "online"), nil
}

var _ application.DeploymentStrategy = NodeStrategy{}
