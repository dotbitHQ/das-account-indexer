package mylog

import (
	"fmt"
	"go.uber.org/zap"
	"runtime/debug"
)

const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelPanic
	LevelFatal
)

const (
	colorDebug = 95
	colorInfo  = 94
	colorWarn  = 93
	colorError = 91
	colorPanic = 91
	colorFatal = 91
)

type logger struct {
	log   *zap.SugaredLogger
	name  string
	level int
}

func (l *logger) ErrStack() {
	l.Warn(string(debug.Stack()))
}

func (l *logger) Debugf(format string, a ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorDebug, l.name, fmt.Sprintf(format, a...))
	l.log.Debug(msg)
}

func (l *logger) Infof(format string, a ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorInfo, l.name, fmt.Sprintf(format, a...))
	l.log.Info(msg)
}

func (l *logger) Warnf(format string, a ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorWarn, l.name, fmt.Sprintf(format, a...))
	l.log.Warn(msg)
}

func (l *logger) Errorf(format string, a ...interface{}) {
	if l.level > LevelError {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorError, l.name, fmt.Sprintf(format, a...))
	l.log.Error(msg)
}

func (l *logger) Panicf(format string, a ...interface{}) {
	if l.level > LevelPanic {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorPanic, l.name, fmt.Sprintf(format, a...))
	l.log.Panic(msg)
}

func (l *logger) Fatalf(format string, a ...interface{}) {
	if l.level > LevelFatal {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorFatal, l.name, fmt.Sprintf(format, a...))
	l.log.Fatal(msg)
}

func (l *logger) Debug(a ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorDebug, l.name, fmt.Sprintln(a...))
	l.log.Debug(msg)
}

func (l *logger) Info(a ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorInfo, l.name, fmt.Sprintln(a...))
	l.log.Info(msg)
}

func (l *logger) Warn(a ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorWarn, l.name, fmt.Sprintln(a...))
	l.log.Warn(msg)
}

func (l *logger) Error(a ...interface{}) {
	if l.level > LevelError {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorError, l.name, fmt.Sprintln(a...))
	l.log.Error(msg)
}

func (l *logger) Panic(a ...interface{}) {
	if l.level > LevelPanic {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorPanic, l.name, fmt.Sprintln(a...))
	l.log.Panic(msg)
}

func (l *logger) Fatal(a ...interface{}) {
	if l.level > LevelFatal {
		return
	}
	msg := fmt.Sprintf("\x1b[%dm▶ [%s] %s\x1b[0m", colorFatal, l.name, fmt.Sprintln(a...))
	l.log.Fatal(msg)
}
