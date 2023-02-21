package trace

import (
	"context"
)

type key struct{}

func WithTIDCtx(ctx context.Context) context.Context {
	return context.WithValue(ctx, key{}, Generate())
}

func SetTIDCtx(ctx context.Context, tid string) context.Context {
	return context.WithValue(ctx, key{}, tid)
}

func GetTIDCtx(ctx context.Context) string {
	v := ctx.Value(key{})

	tid, ok := v.(string)
	if !ok {
		return Generate()
	}

	return tid
}
