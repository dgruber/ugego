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

func TestParseUserList(t *testing.T) {
	ul := `name    exclude
type    ACL
fshare  13
oticket 77
entries daniel,root,%wheel`
	u, err := ParseUserList(ul)
	if err != nil {
		t.Error(err)
		return
	}
	if u.Name != "exclude" {
		t.Errorf("Name is not exclude, it is %s", u.Name)
	}
	if u.Type != "ACL" {
		t.Errorf("Type is not ACL, it is %s", u.Type)
	}
	if u.FShare != 13 {
		t.Errorf("FShare is not 13, it is %d", u.FShare)
	}
	if u.OTicket != 77 {
		t.Errorf("OTicket is not 77, it is %d", u.OTicket)
	}
	if len(u.Entries) != 3 {
		t.Errorf("Entries are not 3 there are %d", u.Entries)
		return
	}
	if u.Entries[0] != "daniel" {
		t.Errorf("First entry needs to be daniel, but it is %s", u.Entries[0])
	}
	if u.Entries[1] != "root" {
		t.Errorf("First entry needs to be root, but it is %s", u.Entries[1])
	}
	if u.Entries[2] != "%wheel" {
		t.Errorf("First entry needs to be %%wheel, but it is %s", u.Entries[2])
	}
}

func TestGetUserLists(t *testing.T) {
	ul, err := GetUserLists("deadlineusers")
	if err != nil {
		t.Errorf("Error during GetUserLists(\"deadlineusers\"): %s", err)
		return
	}
	if ul == nil {
		t.Error("UserList is nil.")
		return
	}
	if len(ul) != 1 {
		t.Errorf("One user list expected, but got %d", len(ul))
		return
	}
	if ul[0].Name != "deadlineusers" {
		t.Errorf("Wrong name in deadlineusers list: %s", ul[0].Name)
	}
}
