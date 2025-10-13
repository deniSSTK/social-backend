package logger

import (
	"log"
	"os"
	e "social-backend/internal/infrastructure/errors"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func Init() error {
	runMode := os.Getenv("RUN_MODE")
	if runMode == "" {
		log.Fatal(e.EnvironmentVariableNotSet.Error() + "RUN_MODE")
	}

	var err error
	once.Do(func() {
		var cfg zap.Config
		if runMode == "prod" {
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
