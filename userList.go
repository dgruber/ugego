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

package ugego

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// UserList is a Univa Grid Engine access control list or Department.
type UserList struct {
	Name    string
	Type    string
	FShare  int
	OTicket int
	Entries []string
}

// ParseUserList parses an Univa Grid Engine ACL by the given name into an struct.
// The expected structure of ul is a multiline string in the Univa Grid Engine format:
// name    exclude
// type    ACL
// fshare  0
// oticket 0
// entries daniel,root,%wheel
func ParseUserList(ul string) (userList *UserList, err error) {
	var ol UserList
	el := strings.Split(ul, "\n")
	if len(el) != 5 {
		return nil, errors.New(fmt.Sprintf("User list does not have 5 lines, it has %d.\n", len(el)))
	}
	name := strings.Split(el[0], "name")
	if len(name) != 2 {
		return nil, errors.New("Error during name parsing.")
	}
	ol.Name = strings.TrimSpace(name[1])
	t := strings.Split(el[1], "type")
	if len(t) != 2 {
		return nil, errors.New("Error during type parsing.")
	}
	ol.Type = strings.TrimSpace(t[1])
	fs := strings.Split(el[2], "fshare")
	if len(fs) != 2 {
		return nil, errors.New("Error during fshare parsing.")
	}
	ot := strings.Split(el[3], "oticket")
	if len(ot) != 2 {
		return nil, errors.New("Error during oticket parsing.")
	}
	ent := strings.Split(el[4], "entries")
	if len(ent) != 2 {
		return nil, errors.New("Error during entries parsing.")
	}
	if ol.FShare, err = strconv.Atoi(strings.TrimSpace(fs[1])); err != nil {
		return nil, err
	}
	if ol.OTicket, err = strconv.Atoi(strings.TrimSpace(ot[1])); err != nil {
		return nil, err
	}
	ol.Entries = strings.Split(strings.TrimSpace(ent[1]), ",")
	return &ol, nil
}

// GetUserLists calls qconf -su <listOfUl> and parses the output
// into UserList structs.
func GetUserLists(userlist ...string) ([]UserList, error) {
	rootPath := os.Getenv("SGE_ROOT")
	if rootPath == "" {
		return nil, errors.New("$SGE_ROOT environment variable not set")
	}
	qconf := fmt.Sprintf("%s/bin/lx-amd64/qconf", rootPath)

	// create comma separated list of user list names
	var csvUserList string
	for k, v := range userlist {
		if k == 0 {
			csvUserList = fmt.Sprintf("%s", v)
		} else {
			csvUserList = fmt.Sprintf("%s,%s", csvUserList, v)
		}
	}

	cmd := exec.Command(qconf, "-su", csvUserList)
	if out, errOut := cmd.CombinedOutput(); errOut == nil {
		// empty line is the delimiter
		uls := strings.Split(string(out), "\n\n")
		if len(uls) == 0 {
			return nil, errors.New(fmt.Sprintf("Could not split output: %s", string(out)))
		}
		outputList := make([]UserList, len(uls), len(uls))
		for i, ul := range uls {
			parsedUserList, errParse := ParseUserList(strings.TrimSpace(uls[i]))
			if errParse != nil {
				log.Printf("Error during parsing user list: %s\n%s\n", errParse, ul)
				log.Printf("%s", ul)
				return nil, errParse
			}
			outputList[i] = *parsedUserList
		}
		return outputList, nil
	} else {
		log.Printf("Error: %s\n", out)
		return nil, errOut
	}
}
