package log

import (
	"context"

	"go.uber.org/zap"
)

type key struct{}

func SetLogCtx(ctx context.Context, tid string) context.Context {
	// デフォルト値でロギング設定をしているので
	// errorは発生しない
	// -> errorはにぎりつぶしても問題ない
	logger, _ := zap.NewProduction()

	log := logger.With(zap.String("tid", tid))

	return context.WithValue(ctx, key{}, log)
}

func GetLogCtx(ctx context.Context) *zap.Logger {
	v := ctx.Value(key{})

	log, ok := v.(*zap.Logger)

	if !ok {
		// デフォルト値でロギング設定をしているので
		// errorは発生しない
		// -> errorはにぎりつぶしても問題ない
		logger, _ := zap.NewProduction()

		return logger
	}

	return log
}
