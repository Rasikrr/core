package sentry

type Config struct {
	SampleRate       float64
	TracesSampleRate float64
	Debug            bool
}
