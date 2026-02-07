//
// Use and distribution licensed under the Apache license version 2.
//
// See the COPYING file in the root project directory for full text.
//

package log

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"

	"github.com/jaypipes/ghw/internal/config"
)

// Info outputs an INFO-level log message to the logger configured in the
// supplied context.
func Info(ctx context.Context, format string, args ...any) {
	logger := config.Logger(ctx)
	if logger == nil || !logger.Enabled(ctx, slog.LevelInfo) {
		return
	}
	var stack [1]uintptr
	runtime.Callers(2, stack[:]) // skip [Callers, Info]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, strings.TrimSpace(fmt.Sprintf(format, args...)), stack[0])
	_ = logger.Handler().Handle(ctx, r)
}

// Warn outputs an WARN-level log message to the logger configured in the
// supplied context.
func Warn(ctx context.Context, format string, args ...any) {
	logger := config.Logger(ctx)
	if logger == nil || !logger.Enabled(ctx, slog.LevelWarn) {
		return
	}
	var stack [1]uintptr
	runtime.Callers(2, stack[:]) // skip [Callers, Warn]
	r := slog.NewRecord(time.Now(), slog.LevelWarn, strings.TrimSpace(fmt.Sprintf(format, args...)), stack[0])
	_ = logger.Handler().Handle(ctx, r)
}

// Debug outputs an DEBUG-level log message to the logger configured in the
// supplied context.
func Debug(ctx context.Context, format string, args ...any) {
	logger := config.Logger(ctx)
	if logger == nil || !logger.Enabled(ctx, slog.LevelDebug) {
		return
	}
	var stack [1]uintptr
	runtime.Callers(2, stack[:]) // skip [Callers, Debug]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, strings.TrimSpace(fmt.Sprintf(format, args...)), stack[0])
	_ = logger.Handler().Handle(ctx, r)
}
