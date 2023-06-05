package logstd

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger represents configuration options for the zap logger. If this is not set in the config, defaults to a
// logger optimized for production use.
// This can be embedded in a viper compatible config struct.
type Logger struct {
	Level       string   `mapstructure:"level"`
	OutputPaths []string `mapstructure:"outputpaths"`
	Encoding    string   `mapstructure:"encoding"`

	hooks       []func(zapcore.Entry) error
	atomicLevel zap.AtomicLevel
}

func (lcfg *Logger) SetAtomicLevel() error {
	atomicLevel, err := zap.ParseAtomicLevel(lcfg.Level)
	if err != nil {
		return err
	}
	lcfg.atomicLevel = atomicLevel
	return nil
}

func (lcfg *Logger) AddHooks(hooks ...func(zapcore.Entry) error) {
	lcfg.hooks = hooks
}

// NewLoggerCfgForTest returns a logging configuration that is optimized for use in a testing environment.
func NewLoggerCfgForTest() *Logger {
	debugAtomicLevel, err := zap.ParseAtomicLevel("debug")
	if err != nil {
		panic(err)
	}

	return &Logger{
		Level:       "debug",
		atomicLevel: debugAtomicLevel,
		Encoding:    "console",
	}
}
