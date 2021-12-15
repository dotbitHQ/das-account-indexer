package mylog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"time"
)

var encoderConfig = zapcore.EncoderConfig{
	TimeKey:        "timer",
	LevelKey:       "level",
	NameKey:        "name",
	CallerKey:      "caller",
	MessageKey:     "message",
	StacktraceKey:  "stacktrace",
	LineEnding:     "\n",
	EncodeLevel:    zapcore.CapitalColorLevelEncoder,
	EncodeTime:     encodeTime, //zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

func encodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.999999999"))
}

func initLog() *zap.SugaredLogger {
	zapConfig := zap.Config{
		Level:             zap.NewAtomicLevelAt(zapcore.DebugLevel),
		Development:       false,
		DisableStacktrace: true,
		Encoding:          "console", //"json",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stdout"},
		DisableCaller:     false,
	}
	zapLogger, _ := zapConfig.Build(zap.AddCallerSkip(1))
	return zapLogger.Sugar()
}

func NewLogger(name string, level int) *logger {
	return &logger{
		name:  name,
		level: level,
		log:   initLog(),
	}
}

func initDefaultLog(fileOut *lumberjack.Logger) *zap.SugaredLogger {
	if fileOut == nil {
		fileOut = &lumberjack.Logger{
			Filename:   "./logs/mylog.log", // log path
			MaxSize:    100,                // log file size, M
			MaxBackups: 30,                 // backups num
			MaxAge:     7,                  // log save days
			LocalTime:  true,
			Compress:   false,
		}
	}
	// zap log
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(fileOut),
			zapcore.AddSync(os.Stdout),
		),
		zap.NewAtomicLevelAt(zapcore.DebugLevel),
	)
	// log
	caller := zap.AddCaller()
	zapLogger := zap.New(core, caller, zap.AddCallerSkip(1))
	return zapLogger.Sugar()
}

func NewLoggerDefault(name string, level int, fileOut *lumberjack.Logger) *logger {
	return &logger{
		name:  name,
		level: level,
		log:   initDefaultLog(fileOut),
	}
}
