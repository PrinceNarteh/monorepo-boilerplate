package config

import (
	"errors"
	"fmt"
	"time"
)

const (
	slowQueryThreshold  = 100 * time.Millisecond // Default threshold for slow queries
	healthCheckInterval = 30 * time.Second       // Default interval for health checks
	healthCheckTimeout  = 5 * time.Second        // Default timeout for health checks
)

// ObservabilityConfig holds the configuration for observability features
type ObservabilityConfig struct {
	ServiceName  string             `koanf:"service_name"  validate:"required"`
	Environment  string             `koanf:"environment"   validate:"required"`
	Logging      LoggingConfig      `koanf:"logging"       validate:"required"`
	NewRelic     NewRelicConfig     `koanf:"new_relic"     validate:"required"`
	HealthChecks HealthChecksConfig `koanf:"health_checks" validate:"required"`
}

// LoggingConfig holds the configuration for logging
type LoggingConfig struct {
	Level              string        `koanf:"level"                validate:"required,oneof=debug info warn error fatal"`
	Format             string        `koanf:"format"               validate:"required,oneof=json text"`
	SlowQueryThreshold time.Duration `koanf:"slow_query_threshold" validate:"required,gt=0"`
}

// NewRelicConfig holds the configuration for New Relic integration
type NewRelicConfig struct {
	LicenseKey                string `koanf:"license_key"                 validate:"required"`
	AppLogForwardingEnabled   bool   `koanf:"app_log_forwarding_enabled"  validate:"required"`
	DistributedTracingEnabled bool   `koanf:"distributed_tracing_enabled" validate:"required"`
	DebugLogging              bool   `koanf:"debug_logging"               validate:"required"`
}

// HealthChecksConfig holds the configuration for health checks
type HealthChecksConfig struct {
	Enabled  bool          `koanf:"enabled"`
	Interval time.Duration `koanf:"interval" validate:"min=1s"`
	Timeout  time.Duration `koanf:"timeout"  validate:"min=1s"`
	Checks   []string      `koanf:"checks"`
}

// DefaultObservabilityConfig returns a default configuration for observability features
// with sensible defaults for a production environment.
func DefaultObservabilityConfig() *ObservabilityConfig {
	return &ObservabilityConfig{
		ServiceName: "api",
		Environment: "development",
		Logging: LoggingConfig{
			Level:              "info",
			Format:             "json",
			SlowQueryThreshold: slowQueryThreshold,
		},
		NewRelic: NewRelicConfig{
			LicenseKey:                "",
			AppLogForwardingEnabled:   true,
			DistributedTracingEnabled: true,
			DebugLogging:              false,
		},
		HealthChecks: HealthChecksConfig{
			Enabled:  true,
			Interval: healthCheckInterval,
			Timeout:  healthCheckTimeout,
			Checks:   []string{"database", "redis"},
		},
	}
}

// Validate checks the ObservabilityConfig for required fields and valid values
func (c *ObservabilityConfig) Validate() error {
	if c.ServiceName == "" {
		return errors.New("service_name is required")
	}

	// Validate log level
	validLevels := map[string]bool{
		"debug": true, "info": true, "warn": true, "error": true,
	}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid logging level: %s (must be one of: debug, info, warn, error)", c.Logging.Level)
	}

	// Validate slow query threshold
	if c.Logging.SlowQueryThreshold < 0 {
		return errors.New("logging slow_query_threshold must be non-negative")
	}

	return nil
}

// GetLogLevel returns the appropriate log level based on the environment
// If the level is not set, it defaults to "info" for production and "debug" for development
func (c *ObservabilityConfig) GetLogLevel() string {
	switch c.Environment {
	case "production":
		if c.Logging.Level == "" {
			return "info"
		}
	case "development":
		if c.Logging.Level == "" {
			return "debug"
		}
	}
	return c.Logging.Level
}

// IsProduction checks if the current environment is production
func (c *ObservabilityConfig) IsProduction() bool {
	return c.Environment == "production"
}
