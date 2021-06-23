package cmd

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogger() *zap.Logger {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(encoderConfig), zapcore.AddSync(colorable.NewColorableStdout()), zapcore.DebugLevel), zap.WithCaller(true))
}
