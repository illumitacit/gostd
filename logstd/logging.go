package logstd

import (
	"go.uber.org/zap"
)

// GetZapConfig returns the configuration for zap.
func GetZapConfig(lcfg *Logger) zap.Config {
	zapC := zap.NewProductionConfig()

	if lcfg == nil {
		return zapC
	}

	if lcfg.Level != "" {
		zapC.Level = lcfg.atomicLevel
	}
	if lcfg.Encoding != "" {
		zapC.Encoding = lcfg.Encoding
	}
	if len(lcfg.OutputPaths) > 0 {
		zapC.OutputPaths = lcfg.OutputPaths
	}

	return zapC
}

// GetSugardLoggerWithoutConfig returns a generic zap sugared logger without configuration. This is useful for emitting
// log messages before the initialization stage.
func GetSugaredLoggerWithoutConfig() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	if err != nil {
		return nil, err
	}
	sugar := logger.Sugar()
	return sugar, nil
}

// GetLogger returns the zap logger object. This can be used to customize the logger object.
func GetLogger(c *Logger) *zap.Logger {
	zapC := GetZapConfig(c)
	logger := zap.Must(zapC.Build())
	return logger
}

// GetSugaredLogger returns the zap sugared logger which can be used to emit log messages from the app.
func GetSugaredLogger(c *Logger) *zap.SugaredLogger {
	return GetLogger(c).Sugar()
}
