/*
 * Copyright 2021 Wang Min Xiang
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * 	http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	if _log == nil {
		_log = logs.New(logs.LogOption{
			Name:             "FNC",
			Formatter:        logs.LogConsoleFormatter,
			ActiveLevel:      logs.LogWarnLevel,
			Colorable:        true,
			EnableStacktrace: false,
		})
	}
	return _log
}
