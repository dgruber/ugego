/*
   Copyright 2015 Daniel Gruber, dgruber@univa.com

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

// Package log provides a set of simple methods for handling logging
// in a similar way (format) like Grid Engine.
package log

import (
	"errors"
	"fmt"
	"github.com/ActiveState/tail"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

// Logger which logs in the style of the Grid Engine's qmaster.

// LogLevel represents the importance of the logging message
type LogLevel int

const (
	Info LogLevel = iota
	Warning
	Error
	Critical
	Profile
)

// LogLevelFilter determines on which level logs should be reported.
// Default initialization is Profile meaning all log messages are
// written. When setting to Warning only Warning, Error, and Critical
// messages are written.
var LogLevelFilter LogLevel

// Profiling determines if profiling information should be printed out
// in the log (P) or not. Per default it is on.
var Profiling bool

// Example: Qmaster messages log
// 07/08/2015 06:07:28.662|          worker|u1010|I|removing trigger to terminate job 3000000278.657

// GELog is a struct which contains the configuration of the logging details.
type GELog struct {
	// Component is the counterpart of the thread name in the Grid Engine scheduler ("worker" above).
	Component string
	// Hostname is the name on which the logging appears.
	Hostname string
	// File points to a file to which to log.
	File *os.File
}

// datelayout defines in which format the date is printed (UGE messages date style)
var datelayout string = "02/01/2006 15:04:05.000"

func init() {
	// default loglevel is profile
	LogLevelFilter = Profile
	Profiling = true
}

// MakeLoggerHostname creates a log object which is used for writing logfiles
// in Grid Engine's master messages log file format. The difference to MakeLog
// is that the hostname is not autoamtically derived. (this is required when
// the hostname in the log output should be different to the current host)
func MakeLoggerHostname(component, hostname string, file *os.File) *GELog {
	var log GELog
	log.Component = component
	log.Hostname = hostname
	log.File = file
	return &log
}

// MakeLogger creates a log object which is used for writing logfiles
// in Grid Engine's master messages log file format.
func MakeLogger(component string, file *os.File) *GELog {
	var log GELog
	log.Component = component
	log.Hostname, _ = os.Hostname()
	log.File = file
	return &log
}

// printMessage prints a logging message for a given logging level but
// only when the logging level is smaller than a certain number
func (g *GELog) printMessage(level, component, format string, a ...interface{}) {
	if level == "P" {
		if Profiling == false {
			return
		}
	}
	if LogLevelFilter > Info {
		switch LogLevelFilter {
		case Warning:
			if level == "I" {
				return
			}
		case Error:
			if level == "I" || level == "W" {
				return
			}
		case Critical:
			if level == "I" || level == "W" || level == "E" {
				return
			}
		}
	}
	msg := fmt.Sprintf(format, a...)
	t := time.Now()
	fmt.Fprintf(g.File, "%s|%17s|%s|%s|%s\n", t.Format(datelayout), component, g.Hostname, level, msg)
}

// InfoC prints an INFO level log message for a given component (like thread).
func (g *GELog) InfoC(component string, format string, a ...interface{}) {
	g.printMessage("I", component, format, a...)
}

// Info prints an INFO level log message using the pre-configured component.
func (g *GELog) Info(format string, a ...interface{}) {
	g.InfoC(g.Component, format, a...)
}

// WarningC prints a WARNING level log message for a given component (like thread).
func (g *GELog) WarningC(component string, format string, a ...interface{}) {
	g.printMessage("W", component, format, a...)
}

// Warning prints a WARNING level log message using the pre-configured component.
func (g *GELog) Warning(format string, a ...interface{}) {
	g.WarningC(g.Component, format, a...)
}

// ErrorC prints an ERROR level log message for a given component (like thread).
func (g *GELog) ErrorC(component string, format string, a ...interface{}) {
	g.printMessage("E", component, format, a...)
}

// Info prints an INFO level log message using the pre-configured component.
func (g *GELog) Error(format string, a ...interface{}) {
	g.ErrorC(g.Component, format, a...)
}

// CriticalC prints a CRITICAL level log message for a given component (like thread).
func (g *GELog) CriticalC(component string, format string, a ...interface{}) {
	g.printMessage("C", component, format, a...)
}

// Crictical prints a CRITICAL level log message using the pre-configured component.
func (g *GELog) Critical(format string, a ...interface{}) {
	g.CriticalC(g.Component, format, a...)
}

// ProfileC prints a PROFILE level log message for a given component (like thread).
func (g *GELog) ProfileC(component string, format string, a ...interface{}) {
	g.printMessage("P", component, format, a...)
}

// Profile prints a PROFILE level log message using the pre-configured component.
func (g *GELog) Profile(format string, a ...interface{}) {
	g.ProfileC(g.Component, format, a...)
}

// CreateProfile writes a profiling entry in the log file for a given event.
func (g *GELog) CreateProfile(start, stop time.Time, event string) {
	duration := stop.Sub(start)
	g.Profile("%s took %s", event, duration)
}

// LogEntry writes a given log entry.
func (g *GELog) LogEntry(entry Entry) {
	g.printMessage(entry.Level, entry.Component, entry.Message)
}

// Entry represents one logging entry, i.e. one line in the logging output.
type Entry struct {
	Time      time.Time
	Component string
	Host      string
	Level     LogLevel
	Message   string
}

// ParseLine parses a string assumed to be in Grid Engine like logging
// and returns a log Entry representing a line. If the line is not
// parsable an error is returned.
func ParseLine(line string) (le Entry, err error) {
	parts := strings.Split(line, "|")
	if len(parts) != 5 {
		return le, errors.New("empty line")
	}
	le.Time, err = time.Parse(datelayout, parts[0])
	le.Component = strings.TrimSpace(parts[1])
	le.Host = strings.TrimSpace(parts[2])
	switch parts[3] {
	case "I":
		le.Level = Info
	case "W":
		le.Level = Warning
	case "E":
		le.Level = Error
	case "C":
		le.Level = Critical
	case "P":
		le.Level = Profile
	default:
		return le, errors.New("Could not parse loglevel.")
	}
	le.Message = strings.TrimSpace(parts[4])
	return le, err
}

// ParseFile parses a given file and converts it into an array of
// logging Entry elements. Note it reads  the complete file in
// memory.
func ParseFile(file *os.File) ([]Entry, error) {
	var last error
	entries := make([]Entry, 0, 0)
	if data, err := ioutil.ReadAll(file); err != nil {
		return nil, err
	} else {
		var entry Entry
		for _, line := range strings.Split(string(data), "\n") {
			if line == "" {
				continue
			}
			if entry, last = ParseLine(line); last == nil {
				entries = append(entries, entry)
			}
		}
	}
	return entries, last
}

// CreateChannel opens a logfile and creates a channel of log
// Entry structs. It keeps the file open and scans it for new
// log file entries forever.
func CreateChannel(file string) (chan Entry, error) {
	lineCh := make(chan Entry, 1024)

	if t, err := tail.TailFile(file, tail.Config{Follow: true}); err == nil {
		// start tailing the logfile in background
		go func() {
			// writing to a closed channel panics - hence
			// we need to recover here since closing is our
			// "official" way to stop tailing
			defer func() {
				t.Stop()
				recover()
			}()
			for line := range t.Lines {
				if entry, errParse := ParseLine(line.Text); errParse == nil {
					// blocking when buffer is full
					lineCh <- entry
				}
			}
		}()
	} else {
		return nil, err
	}
	return lineCh, nil
}

// ParseLevel parses a string and interprets it as a LogLevel.
func ParseLevel(level string) (LogLevel, error) {
	if level == "info" || level == "i" || level == "INFO" || level == "I" {
		return Info, nil
	}
	if level == "warning" || level == "w" || level == "WARNING" || level == "W" {
		return Warning, nil
	}
	if level == "error" || level == "e" || level == "ERROR" || level == "E" {
		return Error, nil
	}
	if level == "critical" || level == "c" || level == "CRITICAL" || level == "C" {
		return Critical, nil
	}
	if level == "profile" || level == "p" || level == "PROFILE" || level == "P" {
		return Profile, nil
	}
	return Info, errors.New("loglevel not parsable")
}
