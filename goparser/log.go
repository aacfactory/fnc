package goparser

import "github.com/aacfactory/logs"

var _log logs.Logs

func EnableLog(debug bool) {
	level := logs.LogInfoLevel
	if debug {
		level = logs.LogDebugLevel
	}
	_log = logs.New(logs.LogOption{
		Name:             "GO PARSER",
		Formatter:        logs.LogConsoleFormatter,
		ActiveLevel:      level,
		Colorable:        true,
		EnableStacktrace: false,
	})
}

func Log() logs.Logs {
	if _log == nil {
		_log = logs.New(logs.LogOption{
			Name:             "GO PARSER",
			Formatter:        logs.LogConsoleFormatter,
			ActiveLevel:      logs.LogWarnLevel,
			Colorable:        true,
			EnableStacktrace: false,
		})
	}
	return _log
}
