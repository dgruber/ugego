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

package accounting

import (
	"bytes"
	"reflect"
	"strconv"
	"time"
)

// Entry is one accounting entry line of the accounting file. Multiple
// entries can for a job when it is a parallel job and the accounting is
// configured to have for each task an accounting entry.
type Entry struct {
	Qname                            string
	Hostname                         string
	GroupID                          string
	Owner                            string
	JobName                          string
	JobNumber                        string
	Account                          string
	PosixPriority                    int
	SubmissionTime                   time.Time
	StartTime                        time.Time
	EndTime                          time.Time
	Failed                           int
	ExitStatus                       int
	RuWallclock                      int64
	RuUtime                          int64
	RuStime                          int64
	RuMaxRSS                         int64
	RuIxRSS                          int64
	RuIsmRSS                         int64
	RuIdRSS                          int64
	RuIsRSS                          int64
	RuMinFlt                         int64
	RuMajFlt                         int64
	RunSwap                          int64
	RuInblock                        int64
	RuOublock                        int64
	RuMsgSnd                         int64
	RuMsgRcv                         int64
	RunSignals                       int64
	RuNvCsw                          int64
	RuNivCsw                         int64
	Project                          string
	Department                       string
	GrantedPE                        string
	Slots                            int
	TaskNumber                       int
	CPU                              float64
	MEM                              float64
	IO                               float64
	Category                         string
	IOW                              float64
	PETaskID                         int
	MaxVmem                          int64
	AdvanceReservationID             int
	AdvanceReservationSubmissionTime time.Time
	JobClass                         string
	QdelInfo                         string
	MaxRSS                           int64
	MaxPSS                           int64
	SubmitHost                       string
	CurrentWorkingDirectory          string
	SubmitCommand                    string
	WallClock                        float64 // milliseconds.<>
}

func sEpochTime(s int64) time.Time {
	return time.Unix(s, 0)
}

func msEpochTime(ms int64) time.Time {
	sec := ms / 1000
	nsec := (ms % 1000) * 1000
	return time.Unix(sec, nsec)
}

func parseTime(column []byte) time.Time {
	timeString := string(column)
	timeInt, err := strconv.ParseInt(timeString, 10, 64)
	if err != nil {
		return time.Unix(0, 0)
	}
	// we assume here that all records dates are after 1990.
	if timeInt > 631152000000 {
		// it mus be a ms time-stamp after 1990
		return msEpochTime(timeInt)
	}
	return sEpochTime(timeInt)
}

// ParseLine reads a line of the Grid Engine accounting file
// and returns a parsed Entry structure.
func ParseLine(line []byte) (Entry, error) {
	// split at ":" (strings are utf-8 hence we are on byte level)
	chunks := bytes.Split(line, []byte(":"))
	for i, _ := range chunks {
		// replace special character to ":" within a chunk
		chunks[i] = bytes.Replace(chunks[i], []byte("\xFF"), []byte(":"), 0)
	}
	var e Entry
	v := reflect.ValueOf(&e).Elem()

	// For each field in the Entry struct discover the type
	// and parse the column based on the type of the field of
	// the struct. That means the order in the Element struct
	// definition must be aligned with the column order in
	// the accounting field. Scanning ends when there is no
	// more column or no more field in the struct, so it is
	// save it use it also with older versions of the accounting
	// file. The only thing which needs to be taken care of
	// is that Univa Grid Engine 8.2 introduced ms since epoch
	// instead of seconds since epoch as date fields But this
	// is also recogized by the magnitude of the number.
	for i := 0; i < v.NumField() && i < len(chunks); i++ {
		field := v.Field(i)
		switch field.Kind() {
		case reflect.String:
			field.SetString(string(chunks[i]))
		case reflect.Int64:
			i64val, _ := strconv.ParseInt(string(chunks[i]), 10, 64)
			field.SetInt(i64val)
		case reflect.Int:
			ival, _ := strconv.Atoi(string(chunks[i]))
			field.SetInt(int64(ival))
		case reflect.Float64:
			fl, _ := strconv.ParseFloat(string(chunks[i]), 64)
			field.SetFloat(fl)
		case reflect.Struct:
			// check if it is a Time (we trust it is from time package)
			if field.Type().Name() == "Time" {
				t := parseTime(chunks[i])
				field.Set(reflect.ValueOf(t))
			}
		}
	}
	return e, nil
}
