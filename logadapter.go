package logadapter

import (
	"fmt"
	"io"
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
	console = NewSimple(LevDEBUG, os.Stderr)
)

func defaultFormat(lev uint8, name string, tick time.Time, caller string, msg string) string {
	timeMsg := tick.Format("01-02 15:04:05.999")
	dest := fmt.Sprintf("[%-18s][%-5s] ", timeMsg, Level(lev))
	if name != "" {
		dest += "<" + name + ">"
	}
	return dest + msg + " (" + caller + ")\r\n"
}

const (
	LevDEBUG uint8 = 0
	LevVERBO uint8 = 1 // high than debug,lower than info.it convenient to debug one problem when debug msg is too much.
	LevINFO  uint8 = 2
	LevWARN  uint8 = 3
	LevERROR uint8 = 4
	LevFATAL uint8 = 5
)

func Level(l uint8) string {
	switch l {
	case LevDEBUG:
		return "DEBUG"
	case LevINFO:
		return "INFO"
	case LevERROR:
		return "ERROR"
	case LevVERBO:
		return "VERBO"
	case LevFATAL:
		return "FATAL"
	default:
		return "WARN"
	}
}

type Logger struct {
	logFunc    func(lev uint8, name string, tick time.Time, caller string, msg string)
	warpFunc   func(lev uint8, name, msg string)
	callerSkip int
	rootLev    uint8
	name       string
}

// NewLogger New console with user implement write log msg function, set callerSkip if call is wrapped.
func NewLogger(rootLev uint8, logFunc func(lev uint8, name string, tick time.Time, caller string, msg string), callerSkip ...int) *Logger {
	if rootLev > LevFATAL {
		panic("invalid root level:" + strconv.Itoa(int(rootLev)))
	}
	logger := &Logger{logFunc: logFunc, rootLev: rootLev}
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
//  warpper := NewByWarp(LevINFO,func(lev uint8, msg string) {
//		logger.Output(2+CallDepth, msg)
//  })
//  warpper.Info("msg")
//
func NewByWarp(rootLev uint8, logFunc func(lev uint8, name, msg string)) *Logger {
	if rootLev > LevFATAL {
		panic("invalid root level:" + strconv.Itoa(int(rootLev)))
	}
	logger := &Logger{warpFunc: logFunc, rootLev: rootLev}
	return logger
}

// NewSimple write as console log to writer.
func NewSimple(rootLev uint8, w io.Writer) *Logger {
	return NewLogger(rootLev, func(lev uint8, name string, tick time.Time, caller string, msg string) {
		fmt.Fprintf(w, defaultFormat(lev, name, tick, caller, msg))
	})
}

// Console get the default Console logger.
func Console() *Logger {
	return console
}

// Named special a name for the log for user to classify.keep for reuse is recommand.
func (l *Logger) Named(name string) *Logger {
	if name == l.name {
		return l
	}
	n := *l
	n.name = name
	return &n
}

func (l *Logger) Debug(args ...interface{}) {
	l.log(LevDEBUG, "", args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log(LevDEBUG, format, args...)
}

// Verbose verbose is design for get output for one problem when debug msg is too much.
// delete or set level to debug after resolved is recommend.
func (l *Logger) Verbose(args ...interface{}) {
	l.log(LevVERBO, "", args...)
}

// Verbosef verbose is design for get output for one problem when debug msg is too much.
// delete or set level to debug after resolved is recommend.
func (l *Logger) Verbosef(format string, args ...interface{}) {
	l.log(LevVERBO, format, args...)
}

// Infof log infomation msg
func (l *Logger) Info(args ...interface{}) {
	l.log(LevINFO, "", args...)
}

// Infof log formatted infomation msg
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log(LevINFO, format, args...)
}

func (l *Logger) Warn(args ...interface{}) {
	l.log(LevWARN, "", args...)
}

func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log(LevWARN, format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.log(LevERROR, "", args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log(LevERROR, format, args...)
}

// Fatal log fatal msg,then call os.Exit(1)
func (l *Logger) Fatal(args ...interface{}) {
	l.log(LevFATAL, "", args...)
	os.Exit(1)
}

// Fatal log formatted fatal msg,then call os.Exit(1)
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log(LevFATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) log(lev uint8, format string, fmtArgs ...interface{}) {
	if lev < l.rootLev {
		return
	}
	msg := format
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(format, fmtArgs...)
	}
	if l.warpFunc != nil {
		l.warpFunc(lev, l.name, msg)
	} else {
		caller := l.caller(runtime.Caller(l.callerSkip + CallDepth))
		l.logFunc(lev, l.name, time.Now(), caller, msg)
	}
}

// caller get caller path(include the last package dir name).
func (l *Logger) caller(pc uintptr, file string, line int, ok bool) string {
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

// Verbose verbose is design for get output for one problem when debug msg is too much.
// delete or set level to debug after resolved is recommend.
func Verbose(args ...interface{}) {
	console.log(LevVERBO, "", args...)
}

// Verbosef verbose is design for get output for one problem when debug msg is too much.
// delete or set level to debug after resolved is recommend.
func Verbosef(format string, args ...interface{}) {
	console.log(LevVERBO, format, args...)
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

// Named special a name for the log for user to classify.keep for reuse is recommand.
func Named(name string) *Logger {
	return console.Named(name)
}

// CurrentStack get stack of current goroutine.
func CurrentStack() string {
	n := 4096
	var trace []byte
	for i := 0; i < 5; i++ {
		trace = make([]byte, n)
		nbytes := runtime.Stack(trace, false)
		if nbytes < len(trace) {
			return string(trace[:nbytes])
		}
		n *= 2
	}
	return string(trace)
}
