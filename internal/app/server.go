package app

import (
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/mil82/internal/api"
	"github.com/fpawel/mil82/internal/charts"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"golang.org/x/sys/windows/registry"
	"net"
	"net/http"
	"net/rpc"
	"strconv"
)

func startHttpServer() func() {

	for _, svcObj := range []interface{}{
		new(api.LastPartySvc),
		new(api.ConfigSvc),
		api.NewPeerSvc(peerNotifier{}),
		api.NewRunnerSvc(runner{}),
		new(api.ChartsSvc),
	} {
		must.AbortIf(rpc.Register(svcObj))
	}

	// Server provide a HTTP transport on /rpc endpoint.
	http.Handle("/rpc", jsonrpc2.HTTPHandler(nil))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		must.Write(w, []byte("hello world"))
	})

	http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Accept", "application/octet-stream")
		bucketID, _ := strconv.ParseInt(r.URL.Query().Get("bucket"), 10, 64)
		charts.WritePointsResponse(w, bucketID)
	})

	srv := new(http.Server)
	lnHTTP, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	addr := "http://" + lnHTTP.Addr().String()
	log.Info(addr)
	key, _, err := registry.CreateKey(registry.CURRENT_USER, `mil82\http`, registry.ALL_ACCESS)
	if err != nil {
		panic(err)
	}
	if err := key.SetStringValue("addr", addr); err != nil {
		panic(err)
	}
	log.ErrIfFail(key.Close)

	go func() {
		err := srv.Serve(lnHTTP)
		if err == http.ErrServerClosed {
			return
		}
		log.PrintErr(err)
		log.ErrIfFail(lnHTTP.Close)
	}()

	return func() {
		if err := srv.Shutdown(ctxApp); err != nil {
			log.PrintErr(err)
		}
	}
}
