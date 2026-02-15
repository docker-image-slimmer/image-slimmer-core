package analyzer

import "time"

type options struct {
	timeout      time.Duration
	retries      int
	metadataOnly bool
}

type Option func(*options)

func defaultOptions() *options {
	return &options{
		timeout:      30 * time.Second,
		retries:      2,
		metadataOnly: false,
	}
}

func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.timeout = d
	}
}

func WithRetries(r int) Option {
	return func(o *options) {
		o.retries = r
	}
}

func WithMetadataOnly() Option {
	return func(o *options) {
		o.metadataOnly = true
	}
}
