/* Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License. */

package log

import (
	"flag"
	"fmt"
	"io"
)

type logLevel int

func (ll logLevel) String() string {
	switch ll {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarning:
		return "warning"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	}

	return "unknown"
}

func (ll *logLevel) Set(value string) error {
	switch value {
	case "trace":
		*ll = LevelTrace
	case "debug":
		*ll = LevelDebug
	case "info":
		*ll = LevelInfo
	case "warning":
		*ll = LevelWarning
	case "error":
		*ll = LevelError
	case "fatal":
		*ll = LevelFatal
	default:
		return fmt.Errorf("Unknown logging level %s", value)
	}

	return nil
}

const (
	// LevelTrace is most detailed logging
	LevelTrace logLevel = 1 + iota
	// LevelDebug is level for debugging logs
	LevelDebug
	// LevelInfo is level for info logs
	LevelInfo
	// LevelWarning is level for warning logs
	LevelWarning
	// LevelError is level for error logs
	LevelError
	// LevelFatal is logging only for fatal errors
	LevelFatal
)

var (
	// Level is current log level for logger
	Level = LevelInfo
	// Writer for writing logs to. You can change it for your own writer
	Writer io.Writer = DefaultWriter{}
)

func init() {
	flag.Var(&Level, "log-level", "Log level: trace|debug|info|warning|error|fatal")
}

func printString(s string) {
	_, err := Writer.Write([]byte(s))
	if err != nil {
		fmt.Printf("Error write log: %s\n", err)
	}
}

func lprint(level logLevel, value interface{}) {
	if level >= Level {
		printString(fmt.Sprint(value))
	}
}

func lprintf(level logLevel, format string, params ...interface{}) {
	lprint(level, fmt.Sprintf(format, params...))
}

// Println is unconditional log
func Println(value interface{}) {
	printString(fmt.Sprint(value))
}

// Printf is unconditional formatted log
func Printf(format string, params ...interface{}) {
	printString(fmt.Sprintf(format, params...))
}

// Trace logging. Use it for most detailed logs
func Trace(value interface{}) { lprint(LevelTrace, value) }

// Tracef is formatted trace logging
func Tracef(format string, params ...interface{}) { lprintf(LevelTrace, format, params...) }

// Debug logging
func Debug(value interface{}) { lprint(LevelDebug, value) }

// Debugf is formatted debug logging
func Debugf(format string, params ...interface{}) { lprintf(LevelDebug, format, params...) }

// Info logging
func Info(value interface{}) { lprint(LevelInfo, value) }

// Infof is formatted info logging
func Infof(format string, params ...interface{}) { lprintf(LevelInfo, format, params...) }

// Warning logging
func Warning(value interface{}) { lprint(LevelWarning, value) }

// Warningf is formatted warning logging
func Warningf(format string, params ...interface{}) { lprintf(LevelWarning, format, params...) }

// Error logging
func Error(value interface{}) { lprint(LevelError, value) }

// Errorf is formatted error logging
func Errorf(format string, params ...interface{}) { lprintf(LevelError, format, params...) }

// Fatal logs fatal error and panic
func Fatal(value interface{}) {
	str := fmt.Sprint(value)
	printString(str)
	panic(str)
}

// Fatalf logs fatal error with format and panic
func Fatalf(format string, params ...interface{}) {
	str := fmt.Sprintf(format, params...)
	printString(str)
	panic(str)
}
