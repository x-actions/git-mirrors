// Copyright 2022 xiexianbin<me@xiexianbin.cn>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	CRITICAL
	FATAL
	PRINT
)

const timeLayout = "2006-01-02 15:04:05"

var (
	logMutex sync.RWMutex
	logLevel LogLevel = INFO
)

func SetLogLevel(level LogLevel) LogLevel {
	logMutex.Lock()
	defer logMutex.Unlock()
	logLevel = level
	return logLevel
}

func printer(level LogLevel, format string, a ...interface{}) {
	builder := strings.Builder{}
	if level != PRINT {
		builder.WriteString(time.Now().Format(timeLayout))
		builder.WriteString(" [")
		switch level {
		case DEBUG:
			builder.WriteString("DEBUG")
		case INFO:
			builder.WriteString(fmt.Sprintf("\x1b[34;1m%s\x1b[0m", "INFO"))
		case WARN:
			builder.WriteString(fmt.Sprintf("\x1b[33;1m%s\x1b[0m", "WARN"))
		case ERROR:
			builder.WriteString(fmt.Sprintf("\x1b[31;1m%s\x1b[0m", "ERROR"))
		case CRITICAL:
			builder.WriteString(fmt.Sprintf("\x1b[31;4m%s\x1b[0m", "CRITICAL"))
		case FATAL:
			builder.WriteString(fmt.Sprintf("\x1b[31;7m%s\x1b[0m", "FATAL"))
		}
		// parent id and pid
		builder.WriteString(fmt.Sprintf("] [%d/%d] ", os.Getppid(), os.Getpid()))
	}
	if format == "" {
		builder.WriteString(fmt.Sprint(a...))
	} else {
		builder.WriteString(fmt.Sprintf(format, a...))
	}

	fmt.Println(builder.String())
}

func Debug(a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= DEBUG {
		printer(DEBUG, "", a...)
	}
}

func Debugf(format string, a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= DEBUG {
		printer(DEBUG, format, a...)
	}
}

func Info(a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= INFO {
		printer(INFO, "", a...)
	}
}

func Infof(format string, a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= INFO {
		printer(INFO, format, a...)
	}
}

func Warn(a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= WARN {
		printer(WARN, "", a...)
	}
}

func Warnf(format string, a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= WARN {
		printer(WARN, format, a...)
	}
}

func Error(a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= ERROR {
		printer(ERROR, "", a...)
	}
}

func Errorf(format string, a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= ERROR {
		printer(ERROR, format, a...)
	}
}

func Critical(a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= CRITICAL {
		printer(CRITICAL, "", a...)
	}
}

func Criticalf(format string, a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= CRITICAL {
		printer(CRITICAL, format, a...)
	}
}

func Fatal(a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= FATAL {
		printer(FATAL, "", a...)
		os.Exit(1)
	}
}

func Fatalf(format string, a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logLevel <= FATAL {
		printer(FATAL, format, a...)
		os.Exit(1)
	}
}

func Print(a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	printer(PRINT, "", a...)
}

func Printf(format string, a ...interface{}) {
	logMutex.Lock()
	defer logMutex.Unlock()
	printer(PRINT, format, a...)
}
