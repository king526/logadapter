package logadapter_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	stdlog "log"

	"github.com/king526/logadapter"
)

func Test_Default(t *testing.T) {
	logadapter.Debug("this msg should be show")
	logadapter.Info("hello")
}

func Test_NewLogger(t *testing.T) {
	logger := logadapter.NewLogger(logadapter.LevVERBO, func(lev uint8, name string, tick time.Time, caller string, msg string) {
		fmt.Println(caller, logadapter.Level(lev), msg)
	})
	logger.Debug("this msg should not be show")
	logger.Info("hello,my logger")
}

func Test_Warp(t *testing.T) {
	logger := stdlog.New(os.Stderr, "mystdlog ", stdlog.Ltime|stdlog.Llongfile)
	warpper := logadapter.NewByWarp(logadapter.LevVERBO, func(lev uint8, name, msg string) {
		logger.Output(2+logadapter.CallDepth, logadapter.Level(lev)+","+msg)
	})
	warpper.Debug("this msg should not be show")
	warpper.Info("hello,my warpper logger")
}

func Test_NewSimple(t *testing.T) {
	fs, _ := os.OpenFile("simple.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	logger := logadapter.NewSimple(logadapter.LevINFO, fs)
	logger.Debug("debug log to NewSimple")
	logger.Info("info log to NewSimple")
	fs.Close()
}

func Test_Named(t *testing.T) {
	logadapter.Named("named").Info("hello")
}

func Test_GetStack(t *testing.T) {
	logadapter.Error("error,the stack:", logadapter.CurrentStack())
}
