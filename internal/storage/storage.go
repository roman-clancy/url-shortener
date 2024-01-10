package storage

import "errors"

var (
	ErrAliasExists   = errors.New("alias exists")
	ErrAliasNotFound = errors.New("alias not found")
)
