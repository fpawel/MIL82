package app

import (
	"context"
	"fmt"
	"github.com/ansel1/merry"
	"github.com/fpawel/comm"
	"github.com/fpawel/comm/comport"
	"github.com/fpawel/gohelp"
	"github.com/fpawel/gohelp/helpstr"
	"github.com/fpawel/mil82/internal/api/notify"
	"github.com/fpawel/mil82/internal/api/types"
	"github.com/fpawel/mil82/internal/cfg"
	"github.com/fpawel/mil82/internal/charts"
	"github.com/powerman/structlog"
	"sync"
)

type WorkFunc = func() error

func runWork(parentCtx context.Context, createNewChart bool, workName string, work WorkFunc) {

	cancelFunc()
	wgWork.Wait()
	wgWork = sync.WaitGroup{}
	ctxWork, cancelFunc = context.WithCancel(parentCtx)

	wgWork.Add(1)

	log = gohelp.NewLogWithSuffixKeys("работа", workName)

	notify.WorkStarted(log, workName)

	if createNewChart {
		charts.CreateNewBucket(workName)
	}

	go func() {

		r := types.WorkResultInfo{Work: workName}

		defer func() {
			if createNewChart {
				charts.SaveLastBucket()
			}
			notify.WorkComplete(log, r)
			log.Printf("%+v", r.Result)
			log.ErrIfFail(readerProducts.Close)
			log.ErrIfFail(responseReaderGasBlock.reader.Close)
			log.ErrIfFail(responseReaderTemperature.reader.Close)
			log = structlog.New()
			wgWork.Done()
		}()

		err := work()

		if err == nil {
			log.Info("выполнено успешно")
			r.Message = "выполнено успешно"
			r.Result = types.WrOk
			return
		}
		if merry.Is(err, context.Canceled) {
			r.Message = "выполнение прервано"
			r.Result = types.WrCanceled
			return
		}

		r.Message = err.Error()
		r.Result = types.WrError

		var kvs []interface{}
		for k, v := range merry.Values(err) {
			strK := fmt.Sprintf("%v", k)
			if strK != "stack" && strK != "msg" && strK != "message" {
				kvs = append(kvs, k, v)
			}
		}
		kvs = append(kvs, "stack", helpstr.FormatMerryStacktrace(err))
		log.PrintErr(err, kvs...)
	}()
}

type reader struct {
	reader       *comport.ReadWriter
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
		Log:            logger,
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
