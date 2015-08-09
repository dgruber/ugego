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

// MakeLog creates a log object which is used for writing logfiles
// in Grid Engine's master messages log file format.
func MakeLog(component, hostname string, file *os.File) *GELog {
	var log GELog
	log.Component = component
	log.Hostname = hostname
	log.File = file
	return &log
}

func (g *GELog) printMessage(level, component, format string, a ...interface{}) {
	layout := "02/01/2006 15:04:05.999"
	msg := fmt.Sprintf(format, a...)
	t := time.Now()
	fmt.Fprintf(g.File, "%s|%17s|%s|%s|%s\n", t.Format(layout), component, g.Hostname, level, msg)
}

// Info prints an INFO level log message for a given component (like thread).
func (g *GELog) Info(component string, format string, a ...interface{}) {
	g.printMessage("I", component, format, a...)
}

// I prints an INFO level log message using the pre-configured component.
func (g *GELog) I(format string, a ...interface{}) {
	g.Info(g.Component, format, a...)
}

// Warning prints a WARNING level log message for a given component (like thread).
func (g *GELog) Warning(component string, format string, a ...interface{}) {
	g.printMessage("W", component, format, a...)
}

// W prints a WARNING level log message using the pre-configured component.
func (g *GELog) W(format string, a ...interface{}) {
	g.Warning(g.Component, format, a...)
}

// Error prints an ERROR level log message for a given component (like thread).
func (g *GELog) Error(component string, format string, a ...interface{}) {
	g.printMessage("E", component, format, a...)
}

// I prints an INFO level log message using the pre-configured component.
func (g *GELog) E(format string, a ...interface{}) {
	g.Error(g.Component, format, a...)
}

// Critical prints a CRITICAL level log message for a given component (like thread).
func (g *GELog) Critical(component string, format string, a ...interface{}) {
	g.printMessage("C", component, format, a...)
}

// C prints a CRITICAL level log message using the pre-configured component.
func (g *GELog) C(format string, a ...interface{}) {
	g.Critical(g.Component, format, a...)
}

// Profile prints a PROFILE level log message for a given component (like thread).
func (g *GELog) Profile(component string, format string, a ...interface{}) {
	g.printMessage("P", component, format, a...)
}

// P prints a PROFILE level log message using the pre-configured component.
func (g *GELog) P(format string, a ...interface{}) {
	g.Profile(g.Component, format, a...)
}
