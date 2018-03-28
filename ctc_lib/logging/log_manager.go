/*
Copyright 2018 Google Inc. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logging

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/GoogleCloudPlatform/runtimes-common/ctc_lib/constants"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	log "github.com/sirupsen/logrus"
)

func NewLogger(dir string, prefix string, level log.Level, enableColors bool) *log.Logger {
	if level == log.DebugLevel {
		// Log to File when verbosity=debug
		logging_dir, err := ioutil.TempDir(dir, prefix)
		if err != nil {
			return handleFileLoggingError(err, level, enableColors)
		}
		path := filepath.Join(logging_dir, constants.LogFileName)
		writer, err := rotatelogs.New(
			path+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(path),
			rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
			rotatelogs.WithRotationTime(time.Duration(86400)*time.Second),
		)
		if err != nil {
			return handleFileLoggingError(err, level, enableColors)
		}
		return &log.Logger{
			Out:       writer,
			Formatter: new(log.JSONFormatter),
			Hooks:     make(log.LevelHooks),
			Level:     log.DebugLevel,
		}
	}
	return getStdOutLogger(level, enableColors)
}

func handleFileLoggingError(err error, level log.Level, enableColors bool) *log.Logger {
	logger := getStdOutLogger(level, enableColors)
	logger.Errorf(`Could not initialize file logger due to %s.
Falling back to StdOut Logging`, err)
	return logger
}
func getStdOutLogger(level log.Level, enableColors bool) *log.Logger {
	return &log.Logger{
		Out:       os.Stderr,
		Formatter: NewCTCLogFormatter(enableColors),
		Hooks:     make(log.LevelHooks),
		Level:     level,
	}
}
func GetCurrentFileName(l *log.Logger) (string, bool) {
	rl, ok := l.Out.(*rotatelogs.RotateLogs)
	return rl.CurrentFileName(), ok
}

// Define Explicit StdOut Loggers which can be used to always print to StdOut.
// This can also add other functionalilty like colored output etc.
var Out = log.New()
