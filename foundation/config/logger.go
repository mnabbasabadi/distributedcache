package config

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

const (
	// LogStrKeyModule is for use with the logger as a key to specify the module name.
	LogStrKeyModule  = "module"
	logStrKeyService = "service"
	// logStrKeyRecoveredValue is for use with the logger as a key to specify the value recovered from a Panic().
	logStrKeyRecoveredValue = "recoveredValue"
)

// SetupLogger ...
func SetupLogger(serviceName, module string, loglevel string, isProduction bool) (zerolog.Logger, error) {
	z := zerolog.New(os.Stderr)

	logger := z.With().Timestamp().
		Str(logStrKeyService, serviceName).
		Str(LogStrKeyModule, module).
		Logger()

	level, err := zerolog.ParseLevel(loglevel)
	if err != nil {
		return logger, fmt.Errorf("invalid log level: %s", loglevel)
	}
	logger = logger.Level(level)
	if isProduction {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		logger = logger.Level(zerolog.InfoLevel)
	}

	return logger, nil
}

// RecoverAndLogPanic ...
func RecoverAndLogPanic(log zerolog.Logger) {
	if err := recover(); err != nil {
		log.Panic().Str(logStrKeyRecoveredValue, fmt.Sprintf("%+v", err)).Stack().Msg("recovered in main")
	}
}
