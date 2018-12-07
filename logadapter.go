package log

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	CallDepth = 2 // logger Public Method call depth.
)

var (
	console = NewLogger(func(lev string, tick time.Time, caller string, msg string) {
		timeMsg := tick.Format("01-02 15:04:05.999")
		fmt.Fprintf(os.Stderr, "[%-18s][%-5s] %s (%s)\r\n", timeMsg, lev, msg, caller)
	})
)

const (
	LevDEBUG = "DEBUG"
	LevINFO  = "INFO"
	LevWARN  = "WARN"
	LevERROR = "ERROR"
	LevFATAL = "FATAL"
)

type logger struct {
	logFunc    func(lev string, tick time.Time, caller string, msg string)
	warpFunc   func(lev string, msg string)
	callerSkip int
}

// NewLogger New console with user implement write log msg function, set callerSkip if call is wrapped.
func NewLogger(logFunc func(lev string, tick time.Time, caller string, msg string), callerSkip ...int) *logger {
	logger := &logger{logFunc: logFunc}
	if len(callerSkip) != 0 {
		logger.callerSkip = callerSkip[0]
	}
	return logger
}

// NewByWarp warp another lib (add time,caller info and other by self), the lib should can set caller skip depth.
// an simple example:
//
//  import "log" ...
//  logger=log.New(os.Stderr, "official log ", log.Ltime|log.Llongfile)
//  warpper := NewByWarp(func(lev string, msg string) {
//		logger.Output(2+CallDepth, msg)
//  })
//  warpper.Info("msg")
//
func NewByWarp(logFunc func(lev string, msg string)) *logger {
	logger := &logger{warpFunc: logFunc}
	return logger
}

// Console get the default Console logger.
func Console() *logger {
	return console
}
func (l *logger) Debug(args ...interface{}) {
	l.log(LevDEBUG, "", args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log(LevDEBUG, format, args...)
}

// Infof log infomation msg
func (l *logger) Info(args ...interface{}) {
	l.log(LevINFO, "", args...)
}

// Infof log formatted infomation msg
func (l *logger) Infof(format string, args ...interface{}) {
	l.log(LevINFO, format, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.log(LevWARN, "", args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.log(LevWARN, format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.log(LevERROR, "", args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log(LevERROR, format, args...)
}

// Fatal log fatal msg,then call os.Exit(1)
func (l *logger) Fatal(args ...interface{}) {
	l.log(LevFATAL, "", args...)
	os.Exit(1)
}

// Fatal log formatted fatal msg,then call os.Exit(1)
func (l *logger) Fatalf(format string, args ...interface{}) {
	l.log(LevFATAL, format, args...)
	os.Exit(1)
}

func (l *logger) log(lev string, format string, fmtArgs ...interface{}) {
	msg := format
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(format, fmtArgs...)
	}
	if l.warpFunc != nil {
		l.warpFunc(lev, msg)
	} else {
		caller := l.caller(runtime.Caller(l.callerSkip + CallDepth))
		l.logFunc(lev, time.Now(), caller, msg)
	}
}

// caller get caller path(include the last package dir name).
func (l *logger) caller(pc uintptr, file string, line int, ok bool) string {
	if !ok {
		return "undefined"
	}
	idx := strings.LastIndexByte(file, '/')
	if idx != -1 {
		idx = strings.LastIndexByte(file[:idx], '/')
	}
	return file[idx+1:] + ":" + strconv.Itoa(line)
}

func Debug(args ...interface{}) {
	console.log(LevDEBUG, "", args...)
}

func Debugf(format string, args ...interface{}) {
	console.log(LevDEBUG, format, args...)
}

// Info use the default logger log INFO level msg.default logger print msg to stderr.
func Info(args ...interface{}) {
	console.log(LevINFO, "", args...)
}

// Infof use the default logger log formatted INFO level msg. default logger print msg to stderr.
func Infof(format string, args ...interface{}) {
	console.log(LevINFO, format, args...)
}

func Warn(args ...interface{}) {
	console.log(LevWARN, "", args...)
}

func Warnf(format string, args ...interface{}) {
	console.log(LevWARN, format, args...)
}

func Error(args ...interface{}) {
	console.log(LevERROR, "", args...)
}

func Errorf(format string, args ...interface{}) {
	console.log(LevERROR, format, args...)
}

// Fatal use the default logger log fatal msg,then call os.Exit(1). default logger print msg to stderr.
func Fatal(args ...interface{}) {
	console.log(LevFATAL, "", args...)
	os.Exit(1)
}

// Fatal use the default logger log formatted fatal msg,then call os.Exit(1). default logger print msg to stderr.
func Fatalf(format string, args ...interface{}) {
	console.log(LevFATAL, format, args...)
	os.Exit(1)
}
