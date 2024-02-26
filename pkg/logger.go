package pkg

import "log/slog"

type NilWriter struct{}

func (NilWriter) Write([]byte) (int, error) { return 0, nil }

func NilLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(NilWriter{}, nil))
}
