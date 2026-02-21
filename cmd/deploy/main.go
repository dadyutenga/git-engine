package main

import (
	"fmt"
	"os"

	"github.com/dadyutenga/git-engine/internal/application"
	"github.com/dadyutenga/git-engine/internal/infrastructure/detectors"
	"github.com/dadyutenga/git-engine/internal/infrastructure/logger"
	"github.com/dadyutenga/git-engine/internal/infrastructure/remote"
	"github.com/dadyutenga/git-engine/internal/infrastructure/ssh"
	"github.com/dadyutenga/git-engine/internal/interfaces/cli"
)

func main() {
	configPath := os.Getenv("DEPLOY_CONFIG")
	if configPath == "" {
		configPath = "configs/config.yaml"
	}

	cfg, err := cli.LoadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	client, err := ssh.New(cfg.SSH)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect via ssh: %v\n", err)
		os.Exit(1)
	}
	defer client.Close() // nolint:errcheck

	exec := remote.Executor{Client: client}
	fs := remote.FileSystem{Exec: exec}
	lockManager := remote.LockManager{Exec: exec}

	strategies := []application.DeploymentStrategy{
		detectors.DockerStrategy{Exec: exec, FS: fs},
		detectors.NodeStrategy{},
		detectors.LaravelStrategy{},
		detectors.PythonStrategy{},
		detectors.StaticStrategy{},
	}

	log := logger.New(os.Stdout)
	app := cli.CLI{
		InitService:     application.InitService{Exec: exec, FS: fs},
		DeployService:   application.DeployService{Exec: exec, FS: fs, Lock: lockManager, Strategies: strategies},
		RollbackService: application.RollbackService{Exec: exec, FS: fs, Strategies: strategies},
		StatusService:   application.StatusService{Exec: exec, FS: fs, Strategies: strategies},
		LogsService:     application.LogsService{Exec: exec},
		Logger:          log,
	}

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
