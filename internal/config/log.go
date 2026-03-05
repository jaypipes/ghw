//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package config

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strings"
)

const (
	envKeyLogLevel  = "GHW_LOG_LEVEL"
	envKeyLogLogfmt = "GHW_LOG_LOGFMT"
)

var (
	logLevelKey     = Key("ghw.log.level")
	defaultLogLevel = slog.LevelWarn
	logLevelVar     = new(slog.LevelVar)
	loggerKey       = Key("ghw.logger")
	logfmtLogger    = slog.New(
		slog.NewTextHandler(
			os.Stderr,
			&slog.HandlerOptions{
				Level: logLevelVar,
			},
		),
	)
	defaultLogger = slog.New(
		&simpleHandler{
			Handler: slog.NewTextHandler(
				os.Stderr,
				&slog.HandlerOptions{
					Level: logLevelVar,
				},
			),
			l: log.New(os.Stderr, "", 0),
		},
	)
)

// simpleHandler is a custom log handler meant to emulate the experience ghw had
// before we moved to log/slog. It implements log/slog.Handler.
type simpleHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *simpleHandler) Handle(
	ctx context.Context,
	r slog.Record,
) error {
	level := r.Level.String() + ":"

	h.l.Printf("%-6s %s", level, r.Message)

	return nil
}

// WithLogLevel allows overriding the default log level of WARN.
func WithLogLevel(level slog.Level) Modifier {
	return func(ctx context.Context) context.Context {
		logLevelVar.Set(level)
		return context.WithValue(ctx, logLevelKey, level)
	}
}

// WithLogLogfmt sets the log output to the logfmt standard.
func WithLogLogfmt() Modifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, loggerKey, logfmtLogger)
	}
}

// EnvOrDefaultLogLogfmt return true if ghw should use logfmt standard output
// format based on the GHW_LOG_LOGFMT environs variable.
func EnvOrDefaultLogLogfmt() bool {
	if _, exists := os.LookupEnv(envKeyLogLogfmt); exists {
		return true
	}
	return false
}

// WithLogger allows overriding the default logger
func WithLogger(logger *slog.Logger) Modifier {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, loggerKey, logger)
	}
}

// LogLevel gets a context's logger's log level or the default if none is set.
func LogLevel(ctx context.Context) slog.Level {
	if ctx == nil {
		return defaultLogLevel
	}
	if v := ctx.Value(logLevelKey); v != nil {
		return v.(slog.Level)
	}
	return defaultLogLevel
}

// EnvOrDefaultLogLevel return true if ghw should not output warnings
// based on the GHW_LOG_LEVEL environs variable.
func EnvOrDefaultLogLevel() slog.Level {
	if ll, exists := os.LookupEnv(envKeyLogLevel); exists {
		switch strings.ToLower(ll) {
		case "debug":
			return slog.LevelDebug
		case "info":
			return slog.LevelInfo
		case "warn", "warning":
			return slog.LevelWarn
		case "err", "error":
			return slog.LevelError
		default:
			return defaultLogLevel
		}
	}
	return defaultLogLevel
}

// Logger gets a context's logger override or the default if none is set.
func Logger(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return defaultLogger
	}
	if v := ctx.Value(loggerKey); v != nil {
		return v.(*slog.Logger)
	}
	return defaultLogger
}

// WithDebug enables verbose debugging output.
func WithDebug() Modifier {
	return WithLogLevel(slog.LevelDebug)
}

// WithDisableWarnings tells ghw not to output warning messages.
func WithDisableWarnings() Modifier {
	return WithLogLevel(slog.LevelError)
}

// EnvOrDefaultDisableWarnings return true if ghw should not output warnings
// based on the GHW_DISABLE_WARNINGS environs variable.
func EnvOrDefaultDisableWarnings() bool {
	if _, exists := os.LookupEnv(envKeyDisableWarnings); exists {
		return true
	}
	return false
}
