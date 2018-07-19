package logger

// Copyright 2017 Microsoft Corporation
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

// LevelType tells a logger the minimum level to log. When code reports a log entry,
// the LogLevel indicates the level of the log entry. The logger only records entries
// whose level is at least the level it was told to log. See the Log* constants.
// For example, if a logger is configured with LogError, then LogError, LogPanic,
// and LogFatal entries will be logged; lower level entries are ignored.
type LevelType uint32

const (
	// LogNone tells a logger not to log any entries passed to it.
	LogNone LevelType = iota

	// LogFatal tells a logger to log all LogFatal entries passed to it.
	LogFatal

	// LogPanic tells a logger to log all LogPanic and LogFatal entries passed to it.
	LogPanic

	// LogError tells a logger to log all LogError, LogPanic and LogFatal entries passed to it.
	LogError

	// LogWarning tells a logger to log all LogWarning, LogError, LogPanic and LogFatal entries passed to it.
	LogWarning

	// LogInfo tells a logger to log all LogInfo, LogWarning, LogError, LogPanic and LogFatal entries passed to it.
	LogInfo

	// LogDebug tells a logger to log all LogDebug, LogInfo, LogWarning, LogError, LogPanic and LogFatal entries passed to it.
	LogDebug
)

// ParseLevel converts the specified string into the corresponding LevelType.
func ParseLevel(s string) (lt LevelType, err error) {
	switch strings.ToLower(s) {
	case "fatal":
		lt = LogFatal
	case "panic":
		lt = LogPanic
	case "error":
		lt = LogError
	case "warning":
		lt = LogWarning
	case "info":
		lt = LogInfo
	case "debug":
		lt = LogDebug
	default:
		err = fmt.Errorf("bad log level '%s'", s)
	}
	return
}

// String implements the stringer interface for LevelType.
func (lt LevelType) String() string {
	switch lt {
	case LogNone:
		return "NONE"
	case LogFatal:
		return "FATAL"
	case LogPanic:
		return "PANIC"
	case LogError:
		return "ERROR"
	case LogWarning:
		return "WARNING"
	case LogInfo:
		return "INFO"
	case LogDebug:
		return "DEBUG"
	default:
		return "UNKNOWN"
	}
}

const (
	// this format provides a fixed number of digits so the size of the timestamp is constant
	logTimestamp = "2006-01-02T15:04:05.0000000Z07:00"
)

// Writer defines methods for writing to a logging facility.
type Writer interface {
	// Writeln writes the specified message with the standard log entry header and new-line character.
	Writeln(level LevelType, message string)

	// Writef writes the specified format specifier with the standard log entry header and no new-line character.
	Writef(level LevelType, format string, a ...interface{})

	// WriteRequest writes the specified HTTP request to the logger if the log level is greater than
	// or equal to LogInfo.  The request body, if set, is logged at level LogDebug or higher.
	WriteRequest(req *http.Request)

	// WriteResponse writes the specified HTTP response to the logger if the log level is greater than
	// or equal to LogInfo.  The response body, if set, is logged at level LogDebug or higher.
	WriteResponse(resp *http.Response)
}

type fileLogger struct {
	logLevel LevelType
	logFile  *os.File
}

// NewFileLogger creates a file-based logger that logs messages of the specified level.
// A File is used instead of a Logger so the stream can be flushed after every write.
func NewFileLogger(level LevelType, dest *os.File) Writer {
	return fileLogger{
		logLevel: level,
		logFile:  dest,
	}
}

func (fl fileLogger) Writeln(level LevelType, message string) {
	fl.Writef(level, "%s\n", message)
}

func (fl fileLogger) Writef(level LevelType, format string, a ...interface{}) {
	if fl.logLevel >= level {
		fmt.Fprintf(fl.logFile, "%s %s", entryHeader(level), fmt.Sprintf(format, a...))
		fl.logFile.Sync()
	}
}

func (fl fileLogger) WriteRequest(req *http.Request) {
	if fl.logLevel >= LogInfo {
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "%s REQUEST: %s %s\n", entryHeader(LogInfo), req.Method, req.URL.String())
		// dump headers
		for k, v := range req.Header {
			if strings.ToLower(k) == "authorization" {
				v = []string{"**REDACTED**"}
			}
			fmt.Fprintf(b, "%s: %s\n", k, strings.Join(v, ","))
		}
		if fl.shouldLogBody(req.Header, req.Body) {
			// dump body
			body, err := ioutil.ReadAll(req.Body)
			if err == nil {
				fmt.Fprintln(b, string(body))
				if nc, ok := req.Body.(io.Seeker); ok {
					// rewind to the beginning
					nc.Seek(0, io.SeekStart)
				} else {
					// recreate the body
					req.Body = ioutil.NopCloser(bytes.NewReader(body))
				}
			} else {
				fmt.Fprintf(b, "failed to read body: %v", err)
			}
		}
		fmt.Fprintln(fl.logFile, b.String())
		fl.logFile.Sync()
	}
}

func (fl fileLogger) WriteResponse(resp *http.Response) {
	if fl.logLevel > LogInfo {
		b := &bytes.Buffer{}
		fmt.Fprintf(b, "%s RESPONSE: %d %s\n", entryHeader(LogInfo), resp.StatusCode, resp.Request.URL.String())
		// dump headers
		for k, v := range resp.Header {
			fmt.Fprintf(b, "%s: %s\n", k, strings.Join(v, ","))
		}
		if fl.shouldLogBody(resp.Header, resp.Body) {
			// dump body
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				fmt.Fprintln(b, string(body))
				resp.Body = ioutil.NopCloser(bytes.NewReader(body))
			} else {
				fmt.Fprintf(b, "failed to read body: %v", err)
			}
		}
		fmt.Fprintln(fl.logFile, b.String())
		fl.logFile.Sync()
	}
}

// returns true if the provided body should be included in the log
func (fl fileLogger) shouldLogBody(header http.Header, body io.ReadCloser) bool {
	ct := header.Get("Content-Type")
	return fl.logLevel >= LogDebug && body != nil && strings.Index(ct, "application/octet-stream") == -1
}

// creates standard header for log entries, it contains a timestamp and the log level
func entryHeader(level LevelType) string {
	return fmt.Sprintf("(%s) %s:", time.Now().Format(logTimestamp), level.String())
}
