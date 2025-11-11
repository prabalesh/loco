package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	Log  *zap.Logger
	once sync.Once
)

func InitLogger(logLevel string) error {
	var err error

	once.Do(func() {
		if err = os.MkdirAll("./logs", 0755); err != nil {
			return
		}

		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   "./logs/app.log",
			MaxSize:    10,
			MaxBackups: 5,
			MaxAge:     30,
			Compress:   true,
		})

		consoleWriter := zapcore.AddSync(os.Stdout)

		loggerLevel := parseLogLevel(logLevel)

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.TimeKey = "timestamp"
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)

		core := zapcore.NewTee(
			zapcore.NewCore(consoleEncoder, consoleWriter, loggerLevel),
			zapcore.NewCore(fileEncoder, fileWriter, loggerLevel),
		)

		Log = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(0),
			zap.AddStacktrace(zapcore.ErrorLevel),
		)

		Log.Info("Logger initialized successfully",
			zap.String("level", loggerLevel.String()),
			zap.String("log_file", "./logs/app.log"),
		)
	})

	return err
}

func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func GetLogger() *zap.Logger {
	if Log == nil {
		panic("logger not initialized, call InitLogger() first")
	}
	return Log
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
