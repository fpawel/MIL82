package main

import (
	"flag"
	"github.com/fpawel/mil82/internal/app"
	"github.com/powerman/structlog"
	"os"
)

func main() {

	logLevel := flag.String("log.level", os.Getenv("MIL82_LOG_LEVEL"), "log `level` (debug|info|warn|err)")
	runPeer := flag.Bool("run.peer", true, "execute gui app (true|false)")
	flag.Parse()

	// Wrong log.level is not fatal, it will be reported and set to "debug".
	structlog.DefaultLogger.SetLogLevel(structlog.ParseLevel(*logLevel))

	app.Run(*runPeer)
}
