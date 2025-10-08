package logger

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func Init(env string) error {
	var err error
	once.Do(func() {
		var cfg zap.Config
		if env == "prod" {
			cfg = zap.NewProductionConfig()
		} else {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		cfg.OutputPaths = []string{"stdout"}
		cfg.ErrorOutputPaths = []string{"stderr"}

		instance, err = cfg.Build()
	})
	return err
}

func Get() *zap.Logger {
	if instance == nil {
		panic("logger not initialized — call logger.Init() first")
	}
	return instance
}
