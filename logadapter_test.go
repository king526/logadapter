package log_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	stdlog "log"

	"github.com/king526/logadapter"
)

func Test_Default(t *testing.T) {
	log.Info("hello")
}

func Test_NewLogger(t *testing.T) {
	logger := log.NewLogger(func(lev string, tick time.Time, caller string, msg string) {
		fmt.Println(caller, lev, msg)
	})
	logger.Info("hello,my logger")
}

func Test_Warp(t *testing.T) {
	logger := stdlog.New(os.Stderr, "mystdlog ", stdlog.Ltime|stdlog.Llongfile)
	warpper := log.NewByWarp(func(lev string, msg string) {
		logger.Output(2+log.CallDepth, lev+","+msg)
	})
	warpper.Info("hello,my warpper logger")
}
