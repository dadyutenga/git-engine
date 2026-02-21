package domain

import "errors"

var (
	// ErrProjectNotFound is returned when the project cannot be located remotely.
	ErrProjectNotFound = errors.New("project not found")
	// ErrLockUnavailable signals that a deployment lock is already held.
	ErrLockUnavailable = errors.New("deployment lock unavailable")
	// ErrUnsupportedProject denotes an unknown project type.
	ErrUnsupportedProject = errors.New("unsupported project type")
)
