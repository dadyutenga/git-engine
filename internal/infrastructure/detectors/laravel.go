package detectors

import (
	"fmt"
	"strings"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// LaravelStrategy handles Laravel/PHP deployments.
type LaravelStrategy struct {
	Exec application.RemoteExecutor
}

// Name returns strategy identifier.
func (LaravelStrategy) Name() string { return "laravel" }

// Detect checks for artisan and composer.json using a single batched command.
func (l LaravelStrategy) Detect(fs application.RemoteFileSystem, project domain.Project) (bool, error) {
	cmd := fmt.Sprintf("( [ -f %s ] && [ -f %s ] ) && echo found || echo missing",
		shell.Escape(project.DeployDir+"/artisan"),
		shell.Escape(project.DeployDir+"/composer.json"))
	out, err := l.Exec.Run(cmd)
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "found"), nil
}

// Deploy installs composer deps and optimizes.
func (LaravelStrategy) Deploy(project domain.Project, exec application.RemoteExecutor) error {
	cmd := fmt.Sprintf("cd %s && composer install --no-dev --optimize-autoloader && php artisan config:cache && php artisan migrate --force", shell.Escape(project.DeployDir))
	_, err := exec.Run(cmd)
	return err
}

// Restart reloads php-fpm if available.
func (LaravelStrategy) Restart(project domain.Project, exec application.RemoteExecutor) error {
	_, err := exec.Run("systemctl restart php-fpm || true")
	return err
}

// Status checks php-fpm activity.
func (LaravelStrategy) Status(project domain.Project, exec application.RemoteExecutor) (bool, error) {
	out, err := exec.Run("systemctl is-active php-fpm || true")
	if err != nil {
		return false, err
	}
	return strings.Contains(out, "active"), nil
}

var _ application.DeploymentStrategy = LaravelStrategy{}
