package core_postgres_pool

import "errors"

var (
	ErrNoRows             = errors.New("ErrNoRows")
	ErrViolatesForeignKey = errors.New("ErrViolatesForeignKey")
	ErrUnknown            = errors.New("ErrUnknown")
)
