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

	portProducts = comport.NewReadWriter(func() comport.Config {
		return comport.Config{
			Baud:        9600,
			ReadTimeout: time.Millisecond,
			Name:        cfg.Get().ComportProducts,
		}
	}, func() comm.Config {
		return comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     1000,
			MaxAttemptsRead:       3,
		}
	})

	portGas = comport.NewReadWriter(func() comport.Config {
		return comport.Config{
			Baud:        9600,
			ReadTimeout: time.Millisecond,
			Name:        cfg.Get().ComportGas,
		}
	}, func() comm.Config {
		return comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     1000,
			MaxAttemptsRead:       3,
		}
	})

	portTemp = comport.NewReadWriter(func() comport.Config {
		return comport.Config{
			Baud:        9600,
			ReadTimeout: time.Millisecond,
			Name:        cfg.Get().ComportTemperature,
		}
	}, func() comm.Config {
		return comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     1000,
			MaxAttemptsRead:       3,
		}
	})
)
