package log

import (
	"log/slog"
	"time"
)

type Attr slog.Attr

func String(key, value string) Attr {
	return Attr{Key: key, Value: slog.StringValue(value)}
}

func Bool(key string, value bool) Attr {
	return Attr{Key: key, Value: slog.BoolValue(value)}
}

func Int(key string, value int) Attr {
	return Attr{Key: key, Value: slog.IntValue(value)}
}

func Float(key string, value float64) Attr {
	return Attr{Key: key, Value: slog.Float64Value(value)}
}

func Err(value error) Attr {
	return Attr{Key: "error", Value: slog.AnyValue(value)}
}

func Duration(key string, value time.Duration) Attr {
	return Attr{Key: key, Value: slog.DurationValue(value)}
}

func Time(key string, value time.Time) Attr {
	return Attr{Key: key, Value: slog.TimeValue(value)}
}

func Any(key string, value any) Attr {
	return Attr{Key: key, Value: slog.AnyValue(value)}
}

func convertAttrs(attrs []Attr) []any {
	anyAttrs := make([]any, len(attrs)*2)
	for i, attr := range attrs {
		anyAttrs[i*2] = attr.Key
		anyAttrs[i*2+1] = attr.Value
	}
	return anyAttrs
}
