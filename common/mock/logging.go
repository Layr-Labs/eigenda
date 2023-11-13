package mock

import (
	"log"

	"github.com/Layr-Labs/eigenda/common"
	ethlog "github.com/ethereum/go-ethereum/log"
)

type Logger struct {
	print bool
}

func NewLogger(print bool) common.Logger {
	return &Logger{
		print: print,
	}
}

func (l *Logger) New(ctx ...interface{}) common.Logger {
	return &Logger{}
}

func (l *Logger) printLog(level ethlog.Lvl, msg string, ctx ...interface{}) {
	if l.print {
		info := []interface{}{
			level,
			msg,
		}
		info = append(info, ctx...)
		log.Println(info)
	}
}

func (l *Logger) SetHandler(h ethlog.Handler) {}

func (l *Logger) Trace(msg string, ctx ...interface{}) {
	l.printLog(ethlog.LvlTrace, msg, ctx...)
}

func (l *Logger) Debug(msg string, ctx ...interface{}) {
	l.printLog(ethlog.LvlDebug, msg, ctx...)
}

func (l *Logger) Info(msg string, ctx ...interface{}) {
	l.printLog(ethlog.LvlInfo, msg, ctx...)
}

func (l *Logger) Warn(msg string, ctx ...interface{}) {
	l.printLog(ethlog.LvlWarn, msg, ctx...)
}

func (l *Logger) Error(msg string, ctx ...interface{}) {
	l.printLog(ethlog.LvlError, msg, ctx...)
}

func (l *Logger) Crit(msg string, ctx ...interface{}) {
	l.printLog(ethlog.LvlCrit, msg, ctx...)
}

func (l *Logger) Fatal(msg string, ctx ...interface{}) {
	l.printLog(ethlog.LvlCrit, msg, ctx...)
}

func (l *Logger) Debugf(template string, args ...interface{}) {}

func (l *Logger) Infof(template string, args ...interface{}) {}

func (l *Logger) Warnf(template string, args ...interface{}) {}

func (l *Logger) Errorf(template string, args ...interface{}) {}

func (l *Logger) Critf(template string, args ...interface{}) {}

func (l *Logger) Fatalf(template string, args ...interface{}) {}
