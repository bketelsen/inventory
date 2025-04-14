package inventory

import "errors"

var (
	ErrNoHostName  = errors.New("no hostname provided")
	ErrEmptyReport = errors.New("report cannot be nil")
	ErrEmptyQuery  = errors.New("query cannot be empty")
	ErrNoResults   = errors.New("no reports found")
)
