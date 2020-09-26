package zaplogger

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		l := zap.NewExample()
		l = l.With(zap.Namespace("hometic"), zap.String("I'm", "gopher"))
		l.Info("PairDevice")
		ctx := context.WithValue(r.Context(), "logger", l)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func L(c context.Context) *zap.Logger {
	contextVal := c.Value("logger")
	if contextVal == nil {
		return zap.NewExample()
	}
	logger, ok := contextVal.(*zap.Logger)
	if ok {
		return logger
	}
	return zap.NewExample()
}
