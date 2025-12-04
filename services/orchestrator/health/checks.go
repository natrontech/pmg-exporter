package health

import (
	"context"
	"runtime"

	"github.com/kelseyhightower/envconfig"
)

// ChecksConfig holds the configuration for specific health check thresholds
type ChecksConfig struct {
	MemoryThreshold     float64 `envconfig:"HEALTH_MEMORY_THRESHOLD" default:"0.9"`       // 90% of total memory
	GoroutinesThreshold int     `envconfig:"HEALTH_GOROUTINES_THRESHOLD" default:"10000"` // Maximum number of goroutines
}

var checksConfig ChecksConfig

// DefaultChecks returns a map of default health checks
func DefaultChecks() map[string]Check {
	_ = envconfig.Process("", &checksConfig) // Process env vars once when checks are initialized
	return map[string]Check{
		"memory": MemoryCheck,
		// "goroutines": GoroutinesCheck,
	}
}

// MemoryCheck checks memory stats
func MemoryCheck(ctx context.Context) Component {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	memoryUsage := float64(m.Alloc) / float64(m.Sys)

	if memoryUsage > checksConfig.MemoryThreshold {
		return Component{
			Name:   "memory",
			Status: StatusDown,
			Error:  "Memory usage exceeds threshold",
		}
	}

	return Component{
		Name:   "memory",
		Status: StatusUp,
	}
}

// GoroutinesCheck monitors number of goroutines
func GoroutinesCheck(ctx context.Context) Component {
	numGoroutines := runtime.NumGoroutine()

	if numGoroutines > checksConfig.GoroutinesThreshold {
		return Component{
			Name:   "goroutines",
			Status: StatusDown,
			Error:  "Too many goroutines",
		}
	}

	return Component{
		Name:   "goroutines",
		Status: StatusUp,
	}
}
