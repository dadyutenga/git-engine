package detectors

import (
	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/domain"
)

// StaticStrategy is a no-op fallback deployer.
type StaticStrategy struct{}

// Name returns identifier.
func (StaticStrategy) Name() string { return "static" }

// Detect always returns true as a fallback.
func (StaticStrategy) Detect(_ application.RemoteFileSystem, _ domain.Project) (bool, error) {
	return true, nil
}

// Deploy performs a no-op to keep interface parity.
func (StaticStrategy) Deploy(_ domain.Project, _ application.RemoteExecutor) error { return nil }

// Restart performs nothing for static assets.
func (StaticStrategy) Restart(_ domain.Project, _ application.RemoteExecutor) error { return nil }

// Status always returns true for static files.
func (StaticStrategy) Status(_ domain.Project, _ application.RemoteExecutor) (bool, error) {
	return true, nil
}

var _ application.DeploymentStrategy = StaticStrategy{}
