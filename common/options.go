package gepis_common

import (
	"github.com/pkg/errors"
)

const (
	defaultJSONOutput  = false
	defaultOutputLevel = "info"
	undefinedAppID     = ""
)

type Options struct {
	appID string
	JSONFormatEnabled bool
	OutputLevel string
}

func (o *Options) SetOutputLevel(outputLevel string) error {
	if toLogLevel(outputLevel) == UndefinedLevel {
		return errors.Errorf("undefined Log Output Level: %s", outputLevel)
	}

	o.OutputLevel = outputLevel
	return nil
}

func (o *Options) SetAppID(id string) {
	o.appID = id
}

func (o *Options) AttachCmdFlags(
	stringVar func(p *string, name string, value string, usage string),
	boolVar func(p *bool, name string, value bool, usage string)) {
	if stringVar != nil {
		stringVar(
			&o.OutputLevel,
			"log-level",
			defaultOutputLevel,
			"Options are debug, info, warn, error, or fatal (default info)")
	}

	if boolVar != nil {
		boolVar(
			&o.JSONFormatEnabled,
			"log-as-json",
			defaultJSONOutput,
			"print log as JSON (default false)")
	}
}

func DefaultOptions() Options {
	return Options{
		JSONFormatEnabled: defaultJSONOutput,
		appID:             undefinedAppID,
		OutputLevel:       defaultOutputLevel,
	}
}

func ApplyOptionsToLoggers(options *Options) error {
	internalLoggers := getLoggers()

	for _, v := range internalLoggers {
		v.EnableJSONOutput(options.JSONFormatEnabled)

		if options.appID != undefinedAppID {
			v.SetAppID(options.appID)
		}
	}

	daprLogLevel := toLogLevel(options.OutputLevel)
	if daprLogLevel == UndefinedLevel {
		return errors.Errorf("invalid value for --log-level: %s", options.OutputLevel)
	}

	for _, v := range internalLoggers {
		v.SetOutputLevel(daprLogLevel)
	}

	return nil
}
