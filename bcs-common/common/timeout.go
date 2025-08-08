package common

import "time"

// HTTPClientTimeout returns the default timeout used for HTTP clients across the project.
func HTTPClientTimeout() time.Duration {
	return 30 * time.Second
}