package logger

import (
	"log"
	"os"
	"strings"

	"sensor-data-service.backend/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger

func New(cfg config.LogConfig) *zap.SugaredLogger {
	if sugar != nil {
		return sugar
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Level))); err != nil {
		level = zapcore.InfoLevel
	}

	var encoderCfg zapcore.EncoderConfig
	if cfg.Pretty {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var encoder zapcore.Encoder
	if cfg.Pretty {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var output zapcore.WriteSyncer
	if strings.ToLower(cfg.Output) == "file" {
		logFile := "logs/app.log"
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		output = zapcore.AddSync(f)
	} else {
		output = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, output, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // skip wrapper
	sugar = logger.Sugar()
	sugar.Infof("logger initialized [level=%s] [pretty=%v] [output=%s]", cfg.Level, cfg.Pretty, cfg.Output)

	return sugar
}
