package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/infrastructure/logger"
)

// CLI wires command flags to application services.
type CLI struct {
	InitService     application.InitService
	DeployService   application.DeployService
	RollbackService application.RollbackService
	StatusService   application.StatusService
	LogsService     application.LogsService
	Logger          logger.Logger
}

// Run parses args and dispatches to the correct service.
func (c CLI) Run(args []string) error {
	if len(args) == 0 {
		c.usage()
		return fmt.Errorf("no command supplied")
	}

	switch args[0] {
	case "init":
		return c.handleInit(args[1:])
	case "push":
		return c.handleDeploy(args[1:])
	case "rollback":
		return c.handleRollback(args[1:])
	case "status":
		return c.handleStatus(args[1:])
	case "logs":
		return c.handleLogs(args[1:])
	default:
		c.usage()
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func (c CLI) handleInit(args []string) error {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() < 1 {
		return fmt.Errorf("project name required")
	}
	project := fs.Arg(0)
	result, err := c.InitService.Init(project)
	if err != nil {
		return err
	}
	c.Logger.Info(result.Message)
	return nil
}

func (c CLI) handleDeploy(args []string) error {
	fs := flag.NewFlagSet("push", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() < 1 {
		return fmt.Errorf("project name required")
	}
	project := fs.Arg(0)
	result, err := c.DeployService.Deploy(project)
	if err != nil {
		return err
	}
	if !result.Success {
		return fmt.Errorf(result.Message)
	}
	c.Logger.Info(result.Message)
	return nil
}

func (c CLI) handleRollback(args []string) error {
	fs := flag.NewFlagSet("rollback", flag.ExitOnError)
	backup := fs.String("backup", "", "backup filename to restore")
	fs.Parse(args)
	if fs.NArg() < 1 {
		return fmt.Errorf("project name required")
	}
	project := fs.Arg(0)
	result, err := c.RollbackService.Rollback(project, *backup)
	if err != nil {
		return err
	}
	c.Logger.Info(result.Message)
	return nil
}

func (c CLI) handleStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	fs.Parse(args)
	if fs.NArg() < 1 {
		return fmt.Errorf("project name required")
	}
	project := fs.Arg(0)
	result, err := c.StatusService.Status(project)
	if err != nil {
		return err
	}
	state := "stopped"
	if result.Running {
		state = "running"
	}
	c.Logger.Info("project=%s exists=%t strategy=%s state=%s", result.ProjectName, result.Exists, result.Strategy, state)
	return nil
}

func (c CLI) handleLogs(args []string) error {
	fs := flag.NewFlagSet("logs", flag.ExitOnError)
	follow := fs.Bool("f", false, "follow log output")
	lines := fs.Int("n", 100, "number of lines")
	fs.Parse(args)
	if fs.NArg() < 1 {
		return fmt.Errorf("project name required")
	}
	project := fs.Arg(0)
	if *lines <= 0 {
		*lines = 100
	}
	return c.LogsService.Tail(project, *lines, *follow, os.Stdout)
}

func (c CLI) usage() {
	msg := `deploy CLI

Usage:
  deploy init <project>
  deploy push <project>
  deploy rollback [-backup filename] <project>
  deploy status <project>
  deploy logs [-f] [-n 100] <project>
`
	_, _ = fmt.Fprintln(os.Stderr, msg)
}
