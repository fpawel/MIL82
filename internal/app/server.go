package app

import (
	"fmt"
	"github.com/fpawel/gohelp/must"
	"github.com/fpawel/mil82/internal/api"
	"github.com/fpawel/mil82/internal/charts"
	"github.com/fpawel/mil82/internal/data"
	"github.com/fpawel/mil82/internal/mil82"
	"github.com/powerman/rpc-codec/jsonrpc2"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/html"
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

	http.HandleFunc("/chart", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Accept", "application/octet-stream")
		bucketID, _ := strconv.ParseInt(r.URL.Query().Get("bucket"), 10, 64)
		charts.WritePointsResponse(w, bucketID)
	})

	http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		mw := minifyHtml.ResponseWriter(w, r)
		defer log.ErrIfFail(mw.Close)
		w = mw

		w.WriteHeader(200)
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Accept", "text/html")
		partyID, _ := strconv.ParseInt(r.URL.Query().Get("party_id"), 10, 64)
		mil82.WriteViewParty(w, partyID)
	})

	http.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("assets"))))

	srv := new(http.Server)
	lnHTTP, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	addr := "http://" + lnHTTP.Addr().String()
	log.Info(fmt.Sprintf("%s/report?party_id=%d", addr, data.LastParty().PartyID))
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

var minifyHtml = minify.New()

func init() {
	minifyHtml.AddFunc("text/html", html.Minify)
}
