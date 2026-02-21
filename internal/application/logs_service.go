package application

import (
	"fmt"
	"io"

	"github.com/dadyutenga/git-engine/internal/domain"
	"github.com/dadyutenga/git-engine/internal/shared/shell"
)

// LogsService streams deployment logs.
type LogsService struct {
	Exec RemoteExecutor
}

// Tail streams the last N lines and follows updates.
func (s LogsService) Tail(projectName string, lines int, follow bool, writer io.Writer) error {
	project := domain.NewProject(projectName)
	cmd := fmt.Sprintf("tail -n %d %s", lines, shell.Escape(project.LogFile))
	if follow {
		cmd = fmt.Sprintf("tail -n %d -F %s", lines, shell.Escape(project.LogFile))
	}
	return s.Exec.RunStream(cmd, writer)
}
