package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(logsPath string, isDebugMode bool, opts ...zap.Option) *zap.SugaredLogger {
	opts = append(opts, zap.AddCaller())

	var writeSyncers = []zapcore.WriteSyncer{
		os.Stderr,
	}

	if len(logsPath) > 0 {
		writeSyncers = append(writeSyncers, getLogWriter(logsPath))
	}

	encoder := getEncoder(isDebugMode)

	var core zapcore.Core
	if isDebugMode {
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncers...), zapcore.DebugLevel)
	} else {
		core = zapcore.NewCore(encoder, zapcore.NewMultiWriteSyncer(writeSyncers...), zapcore.ErrorLevel)
	}

	logger := zap.New(core, opts...)
	return logger.Sugar()
}

func getEncoder(isDebugMode bool) zapcore.Encoder {
	var encoderConfig zapcore.EncoderConfig
	if !isDebugMode {
		encoderConfig = zap.NewProductionEncoderConfig()
	} else {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(logsPath string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logsPath,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}
