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
	log.Debug("this msg should be show")
	log.Info("hello")
}

func Test_NewLogger(t *testing.T) {
	logger := log.NewLogger(log.LevVERBO, func(lev uint8, name string, tick time.Time, caller string, msg string) {
		fmt.Println(caller, log.Level(lev), msg)
	})
	logger.Debug("this msg should not be show")
	logger.Info("hello,my logger")
}

func Test_Warp(t *testing.T) {
	logger := stdlog.New(os.Stderr, "mystdlog ", stdlog.Ltime|stdlog.Llongfile)
	warpper := log.NewByWarp(log.LevVERBO, func(lev uint8, name, msg string) {
		logger.Output(2+log.CallDepth, log.Level(lev)+","+msg)
	})
	warpper.Debug("this msg should not be show")
	warpper.Info("hello,my warpper logger")
}

func Test_Named(t *testing.T) {
	log.Named("named").Info("hello")
}

func Test_GetStack(t *testing.T) {
	log.Error("error,the stack:", log.CurrentStack())
}
