package domain

import (
	"fmt"
	"path/filepath"
)

const (
	repoBasePath   = "/var/repo"
	deployBasePath = "/var/www"
	backupBasePath = "/var/backups"
	lockBasePath   = "/var/locks"
	logBasePath    = "/var/log/deploy"
)

// Project models a deployable project and the required remote paths.
type Project struct {
	Name      string
	RepoPath  string
	DeployDir string
	BackupDir string
	LockFile  string
	LogFile   string
}

// NewProject builds a project with opinionated remote paths.
func NewProject(name string) Project {
	return Project{
		Name:      name,
		RepoPath:  filepath.Join(repoBasePath, fmt.Sprintf("%s.git", name)),
		DeployDir: filepath.Join(deployBasePath, name),
		BackupDir: filepath.Join(backupBasePath, name),
		LockFile:  filepath.Join(lockBasePath, fmt.Sprintf("%s.lock", name)),
		LogFile:   filepath.Join(logBasePath, fmt.Sprintf("%s.log", name)),
	}
}
