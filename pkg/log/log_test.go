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

package log

import (
	"os"
	"testing"
)

func TestMakeLog(t *testing.T) {
	log := MakeLogger("scheduler", os.Stdout)
	log.Info("my test %s", "info")
	log.InfoC("test", "my test %s", "info")
	log.Warning("my test %s", "warning")
	log.WarningC("test", "my test %s", "warning")
	log.Error("my test %s", "error")
	log.ErrorC("test", "my test %s", "error")
	log.Critical("my test %s", "critical")
	log.CriticalC("test", "my test %s", "critical")
	log.Profile("my test %s", "profile")
	log.ProfileC("test", "my test %s", "profile")
}

func TestLogLevelFilter(t *testing.T) {
	LogLevelFilter = Error
	fs, _ := os.Create("test_file_output_logging")
	log := MakeLogger("scheduler", fs)
	log.Info("my test %s", "info")
	log.InfoC("test", "my test %s", "info")
	log.Warning("my test %s", "warning")
	log.WarningC("test", "my test %s", "warning")
	log.Error("my test %s", "error")
	log.ErrorC("test", "my test %s", "error")
	log.Critical("my test %s", "critical")
	log.CriticalC("test", "my test %s", "critical")
	log.Profile("my test %s", "profile")
	log.ProfileC("test", "my test %s", "profile")
	// TODO read file and check output
	file, err := os.Open("test_file_output_logging")
	if err != nil {
		t.Error("error opening file")
	}
	entries, errFile := ParseFile(file)
	if errFile != nil {
		t.Error("error parsing file: ", errFile)
	}
	if len(entries) != 6 {
		t.Log(entries)
		t.Log(len(entries))
		for k, v := range entries {
			t.Log(k, v)
		}
		t.Error("couldn't parse 6 entries")
	}
	file.Close()
	if os.Remove("test_file_output_logging") != nil {
		t.Error("could not delete test file")
	}
}

func TestParseLine(t *testing.T) {
	line, err := ParseLine("07/08/2015 06:07:28.662|          worker|u1010|I|removing trigger to terminate job 3000000278.657")
	if err != nil {
		panic(err)
	}
	if line.Component != "worker" {
		t.Error("Component is not recognized")
	}
	if line.Host != "u1010" {
		t.Error("Host is not recognized")
	}
	layout := "02/01/2006 15:04:05.000"
	if tm := line.Time.Format(layout); tm != "07/08/2015 06:07:28.662" {
		t.Error("Time is not correct")
	}
}

func TestParseLogLevel(t *testing.T) {
	testwords := map[string]LogLevel{"I": Info, "warning": Warning, "C": Critical, "profile": Profile, "ERROR": Error}
	for k, v := range testwords {
		if res, err := ParseLevel(k); err == nil {
			if res != v {
				t.Errorf("loglevel %s was not parsed correctly.", k)
			}
		} else {
			t.Log(k)
			t.Error(err)
		}
	}
}
