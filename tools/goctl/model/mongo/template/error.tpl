package model

import (
	"errors"

	"github.com/JellyGoFF/FF-Hexas/core/stores/mon"
)

var (
	ErrNotFound        = mon.ErrNotFound
	ErrInvalidObjectId = errors.New("invalid objectId")
)
