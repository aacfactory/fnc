package codegen

import "github.com/aacfactory/logs"

var _log logs.Logs

func CreateLog(debug bool) {
	level := logs.LogInfoLevel
	if debug {
		level = logs.LogDebugLevel
	}
	_log = logs.New(logs.LogOption{
		Name:             "FNC",
		Formatter:        logs.LogConsoleFormatter,
		ActiveLevel:      level,
		Colorable:        true,
		EnableStacktrace: false,
	})
}

func Log() logs.Logs {
	return _log
}
