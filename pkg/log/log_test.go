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
	"fmt"
	"os"
	"testing"
	"time"
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

func TestParseLevel(t *testing.T) {
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

func TestCreateProfile(t *testing.T) {
	start := time.Now()
	log := MakeLogger("test", os.Stdout)
	log.CreateProfile(start, time.Now(), "test1")
	log.CreateProfile(start, time.Now(), "test2")
	log.CreateProfile(start, time.Now(), "test3")
}

func TestCreateChannel(t *testing.T) {
	name := "test_create_channel.tmp"
	// create a file and fill it with a line
	f, errCreate := os.Create(name)
	if errCreate != nil {
		t.Fatalf("Can't create test file: %s", errCreate)
	}
	defer func() {
		f.Close()
		os.Remove(name)
	}()
	log := MakeLogger("TestCreateChannel", f)
	LogLevelFilter = Info
	log.Critical("critical error happend")
	log.Info("critical written")

	// attach a channel to it
	ch, err := CreateChannel(name)
	if err != nil {
		t.Fatalf("Error when creating channel: %s", err)
	}

	// read first line out
	line := <-ch

	if line.Level != Critical {
		t.Fatalf("Loglevel not recognized correctly.")
	}

	// add more entries
	log.Error("error happened")
	log.Info("something interesting happened")
	log.Profile("this took soo long")
	log.Warning("didn't found that but I could revocer easily - no worries")

	// check if channel returned all entries
	amount := 0
	for ent := range ch {
		if err != nil {
			t.Fatalf("error when getting channel entries happened: %s", err)
		}
		fmt.Println(ent)
		amount++
		if amount >= 4 {
			t.Log("Got 4 more entries")
			break
		}
	}

	// close channel?
}
