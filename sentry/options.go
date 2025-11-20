package sentry

// Option функция для настройки конфига
type Option func(*Config)

// WithSampleRate устанавливает sample rate для событий
func WithSampleRate(rate float64) Option {
	return func(c *Config) {
		c.SampleRate = rate
	}
}

// WithTracesSampleRate устанавливает sample rate для трейсов
func WithTracesSampleRate(rate float64) Option {
	return func(c *Config) {
		c.TracesSampleRate = rate
	}
}

// WithDebug включает debug режим
func WithDebug(debug bool) Option {
	return func(c *Config) {
		c.Debug = debug
	}
}
