package usecase

import "log/slog"

func logWithSource(log slog.Logger, methodName string) *slog.Logger {
	return log.With(slog.String("method", methodName))
}
