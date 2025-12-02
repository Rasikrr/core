package sentry

import (
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

func Init(config Config, env enum.Environment) error {
	var initErr error
	once.Do(func() {
		initErr = sentrySDK.Init(sentrySDK.ClientOptions{
			Dsn:              config.DSN,
			Environment:      env.String(),
			Release:          version.GetVersion(),
			SampleRate:       config.SampleRate,
			TracesSampleRate: config.TracesSampleRate,
			AttachStacktrace: true,
			EnableTracing:    config.Tracing,
			Debug:            config.Debug,
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

func Enabled() bool {
	return enabled.Load()
}

func CurrentHub() (*sentrySDK.Hub, error) {
	if Enabled() {
		return sentrySDK.CurrentHub(), nil
	}
	return nil, errors.New("sentry: sentry not enabled")
}

func Flush(timeout time.Duration) bool {
	if !Enabled() {
		return true
	}
	return sentrySDK.Flush(timeout)
}
