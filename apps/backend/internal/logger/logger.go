// Package logger provides a structured logging service using New Relic and zerolog.
// It initializes a New Relic application for distributed tracing and performance monitoring.
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	// "github.com/newrelic/go-agent/v3/integrations/logcontext-v2/zerologWriter"
	"github.com/newrelic/go-agent/v3/newrelic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/PrinceNarteh/go-boilerplate/internal/config"
)

// LoggerService provides logging capabilities using New newrelic
// and zerolog for structured logging.
type LoggerService struct {
	nrApp *newrelic.Application
}

// NewLoggerService creates a new instance of LoggerService.
// It initializes the New Relic application with the provided configuration.
// If the configuration is nil, it returns nil.
// If there is an error creating the New Relic application, it panics with an error message.
// The New Relic application is configured with the service name, license key,
// distributed tracing enabled, and a debug logger using zerolog.
// The zerolog logger is configured to output to the console in a human-readable format.
// This service can be used to log application events, errors, and performance metrics.
// It is recommended to use this service for all logging needs in the application
func NewLoggerService(cfg *config.ObservabilityConfig) *LoggerService {
	svc := &LoggerService{}

	if cfg.NewRelic.LicenseKey == "" {
		log.Println("New Relic license key not provided, skipping initialization")
		return svc
	}

	var configOpts []newrelic.ConfigOption
	configOpts = append(configOpts,
		newrelic.ConfigAppName(cfg.ServiceName),
		newrelic.ConfigLicense(cfg.NewRelic.LicenseKey),
		newrelic.ConfigAppLogForwardingEnabled(cfg.NewRelic.AppLogForwardingEnabled),
		newrelic.ConfigDistributedTracerEnabled(cfg.NewRelic.DistributedTracingEnabled),
	)

	if cfg.NewRelic.DebugLogging {
		configOpts = append(configOpts, newrelic.ConfigDebugLogger(os.Stdout))
	}

	app, err := newrelic.NewApplication(configOpts...)
	if err != nil {
		log.Printf("Failed to initialized New Relic:  %v\n", err)
	}

	svc.nrApp = app
	log.Printf("New Relic initialized for app: %s\n", cfg.ServiceName)

	return svc
}

// Shutdown shuts down New Relic
func (ls *LoggerService) Shutdown() {
	if ls.nrApp != nil {
		ls.nrApp.Shutdown(10 * time.Second)
	}
}

// GetApplication returns the New Relic application instance
func (ls *LoggerService) GetApplication() *newrelic.Application {
	return ls.nrApp
}

// NewLoggerWithService creates a logger with full config and logger service
func NewLoggerWithService(cfg *config.ObservabilityConfig, loggerService *LoggerService) zerolog.Logger {
	var logLevel zerolog.Level
	level := cfg.GetLogLevel()

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	// Don't set global level - let each logger have its own level
	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var writer io.Writer

	// Setup base writer
	var baseWriter io.Writer
	if cfg.IsProduction() && cfg.Logging.Format == "json" {
		// In production, write to stdout
		baseWriter = os.Stdout

		// Wrap with New Relic zerologWriter for log forwarding in production
		// TODO: Re-enable zerologWriter integration when dependency issue is resolved
		// if loggerService != nil && loggerService.nrApp != nil {
		// 	nrWriter := zerologWriter.New(baseWriter, loggerService.nrApp)
		// 	writer = nrWriter
		// } else {
			writer = baseWriter
		// }
	} else {
		// Development mode - use console writer
		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
		writer = consoleWriter
	}

	// Note: New Relic log forwarding is now handled automatically by zerologWriter integration

	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", cfg.ServiceName).
		Str("environment", cfg.Environment).
		Logger()

	// Include stack traces for errors in development
	if !cfg.IsProduction() {
		logger = logger.With().Stack().Logger()
	}

	return logger
}

// WithTraceContext adds New Relic transaction context to logger
func WithTraceContext(logger zerolog.Logger, txn *newrelic.Transaction) zerolog.Logger {
	if txn == nil {
		return logger
	}

	// Get trace metadata from transaction
	metadata := txn.GetTraceMetadata()

	return logger.With().
		Str("trace.id", metadata.TraceID).
		Str("span.id", metadata.SpanID).
		Logger()
}

// NewPgxLogger creates a database logger
func NewPgxLogger(level zerolog.Level) zerolog.Logger {
	writer := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: "2006-01-02 15:04:05",
		FormatFieldValue: func(i any) string {
			switch v := i.(type) {
			case string:
				// Clean and format SQL for better readability
				if len(v) > 200 {
					// Truncate very long SQL statements
					return v[:200] + "..."
				}
				return v
			case []byte:
				var obj interface{}
				if err := json.Unmarshal(v, &obj); err == nil {
					pretty, _ := json.MarshalIndent(obj, "", "    ")
					return "\n" + string(pretty)
				}
				return string(v)
			default:
				return fmt.Sprintf("%v", v)
			}
		},
	}

	return zerolog.New(writer).
		Level(level).
		With().
		Timestamp().
		Str("component", "database").
		Logger()
}

// GetPgxTraceLogLevel converts zerolog level to pgx tracelog level
func GetPgxTraceLogLevel(level zerolog.Level) int {
	switch level {
	case zerolog.DebugLevel:
		return 6 // tracelog.LogLevelDebug
	case zerolog.InfoLevel:
		return 4 // tracelog.LogLevelInfo
	case zerolog.WarnLevel:
		return 3 // tracelog.LogLevelWarn
	case zerolog.ErrorLevel:
		return 2 // tracelog.LogLevelError
	default:
		return 0 // tracelog.LogLevelNone
	}
}
