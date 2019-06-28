package hard

import (
	"context"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/comport"
	"github.com/fpawel/comm/modbus"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/mil82/internal/api"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/powerman/structlog"
	"time"
)

func ResetContext(parentCtx context.Context) {
	ctx, cancelFunc = context.WithCancel(parentCtx)
}

func ReadProductVar(addr modbus.Addr, Var modbus.Var, ctx context.Context, log *structlog.Logger) (float64, error) {

	log = gohelp.LogWithKeys(log, "адрес", addr, "var", Var)
	value, err := modbus.Read3BCD(log, responseReaderProducts(ctx), addr, Var)
	if err == nil {
		notify.ReadVar(log, api.AddrVarValue{Addr: addr, Var: Var, Value: value})
	}
	return value, err
}

func Close() {
	log.ErrIfFail(readerProducts.Close)
	log.ErrIfFail(responseReaderGasBlock.reader.Close)
	log.ErrIfFail(responseReaderTemperature.reader.Close)
}

type reader struct {
	reader       *comport.Reader
	config       comm.Config
	portNameFunc func() string
	ctx          context.Context
}

func (x reader) GetResponse(logger *structlog.Logger, request []byte, responseParser comm.ResponseParser) ([]byte, error) {
	if !x.reader.Opened() {
		if err := x.reader.Open(x.portNameFunc()); err != nil {
			return nil, err
		}
	}
	return x.reader.GetResponse(comm.Request{
		Logger:         logger,
		Bytes:          request,
		Config:         x.config,
		ResponseParser: responseParser,
	}, x.ctx)
}

func responseReaderProducts(ctx context.Context) reader {
	return reader{
		reader: readerProducts,
		config: comm.Config{
			ReadByteTimeoutMillis: 50,
			ReadTimeoutMillis:     1000,
			MaxAttemptsRead:       3,
		},
		portNameFunc: func() string {
			return cfg.Get().ComportProducts
		},
		ctx: ctx,
	}
}

var (
	log           = structlog.New()
	cancelFunc    = func() {}
	skipDelayFunc = func() {}
	ctx           = context.TODO()

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
