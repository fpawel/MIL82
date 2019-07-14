package app

import (
	"context"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/comport"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/powerman/structlog"
	"sync"
	"time"
)

var (
	log    = structlog.New()
	ctxApp context.Context

	cancelFunc    = func() {}
	skipDelayFunc = func() {}
	ctxWork       = context.TODO()
	wgWork        sync.WaitGroup

	readerProducts = comport.NewReader(comport.Config{
		Baud:        9600,
		ReadTimeout: time.Millisecond,
	})

	responseReaderGasBlock = reader{
		reader: comport.NewReader(comport.Config{
			Baud:        9600,
			ReadTimeout: time.Millisecond,
		}),
		config: comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     1000,
			MaxAttemptsRead:       3,
		},
		portNameFunc: func() string {
			return cfg.Get().ComportGas
		},
	}

	responseReaderTemperature = reader{
		reader: comport.NewReader(comport.Config{
			Baud:        9600,
			ReadTimeout: time.Millisecond,
		}),
		config: comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     1000,
			MaxAttemptsRead:       3,
		},
		portNameFunc: func() string {
			return cfg.Get().ComportTemperature
		},
	}
)
