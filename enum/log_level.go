package enum

import "log/slog"

//go:generate enumer -type=LogLevel -text -json -trimprefix LogLevel -transform=snake -output log_level_enumer.go -comment "log level"

type LogLevel uint8

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

var (
	ToSlogLevelMap = map[LogLevel]slog.Leveler{
		LogLevelDebug: slog.LevelDebug,
		LogLevelInfo:  slog.LevelInfo,
		LogLevelWarn:  slog.LevelWarn,
		LogLevelError: slog.LevelError,
	}
)

func (l LogLevel) ToSlogLevel() slog.Leveler {
	out, ok := ToSlogLevelMap[l]
	if !ok {
		return slog.LevelInfo
	}
	return out
}
