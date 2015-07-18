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

package uge_go_helper

import (
	"testing"
)

// TestQstatf tests the qstat -f -q all.q,all.q -xml command and its
// parsing. The requirements for the test are abvously that qstat needs
// to work in the shell where called and the cluster needs to have
// a queue called all.q defined.
func TestQstatf(t *testing.T) {
	// all.q needs to exist!
	ql, err := Qstatf("all.q,all.q")
	if err != nil {
		t.Errorf("Error during qstat_f: %s", err)
	}
	if len(ql) <= 0 {
		t.Error("Length of parsed output is 0 - probably all.q does not exist")
		return
	}
	// TODO check output
}
