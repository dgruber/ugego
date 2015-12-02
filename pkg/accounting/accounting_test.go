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
	"testing"
)

var testline1 string = "vcontrol.q:maui:sgegroup:sgetest:vcontrol:1:sge:0:1448376371581:1448376371857:1448438703533:100:137:62331.676:0.121:0.061:5456:0:0:0:0:23373:0:0:16:168:0:0:0:352:84:NONE:defaultdepartment:NONE:1:0:85.400:233.209213:0.796364:-jc vcontrol.default -q vcontrol.q -l h_stack=16M,s_stack=16M -binding no_job_binding 0 0 0 0 no_explicit_binding:0.000000:NONE:6099869696:0:0:vcontrol.default:sgetest@maui:0:0:maui:NONE:NONE:62331.836000"

func TestParseLine(t *testing.T) {
	entry, err := ParseLine([]byte(testline1))
	if err != nil {
		t.Fatalf("Following error: %s", err)
	}
	t.Logf("%v", entry)
	if entry.Qname != "vcontrol.q" {
		t.Errorf("Qname is not vcontrol.q: %s", entry.Qname)
	}
	if entry.Hostname != "maui" {
		t.Errorf("Host name is not maui: %s", entry.Hostname)
	}
	if entry.SubmitHost != "maui" {
		t.Errorf("Submit host is not maui: %s", entry.SubmitCommand)
	}
	if entry.GroupID != "sgegroup" {
		t.Errorf("GroupID is not sgegroup: %s", entry.GroupID)
	}
	// ..
	if entry.Failed != 100 {
		t.Errorf("Failed is not 100 it is %d", entry.Failed)
	}
	// ...
	if entry.PosixPriority != 0 {
		t.Errorf("POSIX priority is not 0 it is %d", entry.PosixPriority)
	}
	if entry.ExitStatus != 137 {
		t.Errorf("Exit status is not 0 it is %d", entry.ExitStatus)
	}
	// ...
	if entry.Slots != 1 {
		t.Errorf("Slots is not 1 it is %d", entry.Slots)
	}
	if entry.TaskNumber != 0 {
		t.Errorf("Task number is not 0 it is %d", entry.TaskNumber)
	}
	// ...
	if entry.Category != "-jc vcontrol.default -q vcontrol.q -l h_stack=16M,s_stack=16M -binding no_job_binding 0 0 0 0 no_explicit_binding" {
		t.Errorf("Entry category is wrong: %s", entry.Category)
	}
}

func BenchmarkParseLine(b *testing.B) {
	line := []byte(testline1)
	for i := 0; i < b.N; i++ {
		ParseLine(line)
	}
}
