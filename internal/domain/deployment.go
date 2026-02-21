package domain

import "time"

// DeploymentResult captures the outcome of a deployment execution.
type DeploymentResult struct {
	ProjectName string
	Success     bool
	Status      string
	Message     string
	LogFile     string
	Timestamp   time.Time
	Details     map[string]string
}

// InitResult represents the output of an init operation.
type InitResult struct {
	Project   Project
	Success   bool
	Message   string
	Timestamp time.Time
}

// RollbackResult represents the outcome of a rollback.
type RollbackResult struct {
	ProjectName string
	Success     bool
	Restored    string
	Message     string
	Timestamp   time.Time
}

// StatusResult describes the remote state of an application.
type StatusResult struct {
	ProjectName string
	Exists      bool
	Running     bool
	Strategy    string
	Message     string
	Timestamp   time.Time
}
