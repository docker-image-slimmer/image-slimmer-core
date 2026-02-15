package analyzer

import "time"

type Metrics struct {
	FetchDuration time.Duration
	BuildDuration time.Duration
	TotalDuration time.Duration
}
