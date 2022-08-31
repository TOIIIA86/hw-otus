package logger

import (
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	Zap *zap.Logger
}

func New(level string, development bool) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}

	config := zap.NewProductionConfig()
	config.Level = lvl
	config.Development = development

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &Logger{Zap: logger}, nil
}

func (l Logger) Debug(msg string) {
	l.Zap.Debug(msg)
}

func (l Logger) Info(msg string) {
	l.Zap.Info(msg)
}

func (l Logger) Warn(msg string) {
	l.Zap.Warn(msg)
}

func (l Logger) Error(msg string) {
	l.Zap.Error(msg)
}

func (l Logger) Fatal(msg string) {
	l.Zap.Fatal(msg)
	os.Exit(1)
}
