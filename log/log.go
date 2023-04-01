package log

import (
	"os"
	"time"

	"golang.org/x/exp/slog"
)

func New() *slog.Logger {
	th := slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove time from the output for predictable test output.
			if a.Key == slog.TimeKey {
				return slog.String("ts", time.Now().Format("15:04:05.000"))
			}
			return a
		},
	}.NewTextHandler(os.Stderr)

	return slog.New(th)
}
