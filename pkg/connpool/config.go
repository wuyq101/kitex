package connpool

import "time"

// IdleConfig contains idle configuration for long-connection pool.
type IdleConfig struct {
	MaxIdlePerAddress int
	MaxIdleGlobal     int
	MaxIdleTimeout    time.Duration
}

const (
	defaultMaxIdleTimeout = 30 * time.Second
	minMaxIdleTimeout     = 3 * time.Second
)

// CheckPoolConfig to check invalid param.
// default MaxIdleTimeout = 30s, min value is 3s
func CheckPoolConfig(config IdleConfig) *IdleConfig {
	if config.MaxIdleTimeout == 0 {
		config.MaxIdleTimeout = defaultMaxIdleTimeout
	} else if config.MaxIdleTimeout < 3*time.Second {
		config.MaxIdleTimeout = minMaxIdleTimeout
	}
	return &config
}
