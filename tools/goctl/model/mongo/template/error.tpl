package model

import (
	"errors"

	"github.com/lemongoff/hexas/core/stores/mon"
)

var (
	ErrNotFound        = mon.ErrNotFound
	ErrInvalidObjectId = errors.New("invalid objectId")
)
