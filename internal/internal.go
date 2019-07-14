package internal

import (
	"os"
	"path/filepath"
)

const (
	EnvVarLogLevel = "MIL82_LOG_LEVEL"
)

func DataDir() string {
	dir := os.Getenv("MIL82_DATA_DIR")
	if len(dir) == 0 {
		dir = filepath.Dir(os.Args[0])
	}
	return dir
}
