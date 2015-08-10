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
	"fmt"
	"os"
	"time"
)

// Logger which logs in the style of the Grid Engine's qmaster.

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

// MakeLoggerHostname creates a log object which is used for writing logfiles
// in Grid Engine's master messages log file format.
func MakeLogger(component string, file *os.File) *GELog {
	var log GELog
	log.Component = component
	log.Hostname, _ = os.Hostname()
	log.File = file
	return &log
}

func (g *GELog) printMessage(level, component, format string, a ...interface{}) {
	layout := "02/01/2006 15:04:05.000"
	msg := fmt.Sprintf(format, a...)
	t := time.Now()
	fmt.Fprintf(g.File, "%s|%17s|%s|%s|%s\n", t.Format(layout), component, g.Hostname, level, msg)
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
