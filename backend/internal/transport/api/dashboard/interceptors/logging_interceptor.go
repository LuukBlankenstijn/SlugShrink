package interceptors

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/LuukBlankenstijn/slugshrink/internal/logging"
)

func LoggingInterceptor() connect.Interceptor {
	return connect.UnaryInterceptorFunc(func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			res, err := next(ctx, req)
			if err != nil {
				code := connect.CodeOf(err)
				level := slog.LevelError
				switch code {
				case connect.CodeInvalidArgument, connect.CodeFailedPrecondition, connect.CodeNotFound, connect.CodeUnauthenticated, connect.CodePermissionDenied:
					level = slog.LevelWarn
				}
				logging.Logger().Log(ctx, level, "dashboard request failed",
					slog.String("procedure", req.Spec().Procedure),
					slog.Any("peer", req.Peer()),
					slog.Any("code", code),
					slog.Any("error", err),
				)
			}
			return res, err
		}
	})
}
