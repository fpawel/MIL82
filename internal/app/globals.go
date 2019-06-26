package app

import (
	"context"
	"github.com/powerman/structlog"
)

var (
	log    = structlog.New()
	ctxApp context.Context
)
