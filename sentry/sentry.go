package sentry

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Rasikrr/core/enum"
	"github.com/Rasikrr/core/version"
	"github.com/cockroachdb/errors"
	sentrySDK "github.com/getsentry/sentry-go"
)

var (
	enabled atomic.Bool
	once    sync.Once
)

func Enabled() bool {
	return enabled.Load()
}

func Init(config Config, appName string, env enum.Environment) error {
	var initErr error
	once.Do(func() {
		initErr = sentrySDK.Init(sentrySDK.ClientOptions{
			ServerName:       appName,
			Dsn:              config.DSN,
			Environment:      env.String(),
			Release:          version.GetVersion(),
			SampleRate:       config.SampleRate,
			TracesSampleRate: config.TracesSampleRate,
			AttachStacktrace: true,
			EnableTracing:    config.Tracing,
			Debug:            config.Debug,
			MaxBreadcrumbs:   100,
			EnableLogs:       config.EnableLogs,
			BeforeSend: func(event *sentrySDK.Event, _ *sentrySDK.EventHint) *sentrySDK.Event {
				return event
			},
		})
		if initErr == nil {
			enabled.Store(true)
		}
	})
	if initErr != nil {
		return fmt.Errorf("sentry.Init: %v", initErr)
	}
	return nil
}

func ClearBreadcrumbs() {
	if !Enabled() {
		return
	}
	sentrySDK.CurrentHub().Scope().ClearBreadcrumbs()
}

// ClearBreadcrumbs clears all breadcrumbs from the current Sentry hub.
// Typically called after application startup to prevent startup logs
// from being included in error events.
func CurrentHub() (*sentrySDK.Hub, error) {
	if Enabled() {
		return sentrySDK.CurrentHub(), nil
	}
	return nil, errors.New("sentry: sentry not enabled")
}

func SetHubOnCtx(ctx context.Context) context.Context {
	hubClone := sentrySDK.CurrentHub().Clone()
	hubClone.Scope().ClearBreadcrumbs()
	return sentrySDK.SetHubOnContext(ctx, hubClone)
}

func GetHubFromContext(ctx context.Context) *sentrySDK.Hub {
	if !Enabled() {
		return nil
	}
	return sentrySDK.GetHubFromContext(ctx)
}

func Recover() {
	if !Enabled() {
		return
	}
	sentrySDK.Recover()
}

func RecoverWithContext(ctx context.Context) {
	if !Enabled() {
		return
	}
	sentrySDK.RecoverWithContext(ctx)
}

func Flush(timeout time.Duration) bool {
	if !Enabled() {
		return true
	}
	return sentrySDK.Flush(timeout)
}
