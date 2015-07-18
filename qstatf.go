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
	"encoding/xml"
	"errors"
	"log"
	"os/exec"
)

// QstatQueue represents an entries out of the qstat -f -xml .. command.
type QstatQueue struct {
	Name       string  `xml:"name"`
	QType      string  `xml:"qtype"`
	SlotsUsed  int     `xml:"slots_used"`
	SlotsResv  int     `xml:"slots_resv"`
	SlotsTotal int     `xml:"slots_total"`
	NPLoadAvg  float64 `xml:"np_load_avg"`
	Arch       string  `xml:"arch"`
	// state is only there if it is not available
	State string `xml:"state"`
}

// QstatQueueInfoList is the representation of the qstat -f -xml .. output.
type QstatQueueInfoList struct {
	XMLName xml.Name `xml:"job_info"`
	// contains queue_info (qstat -f -xml)
	QueueList []QstatQueue `xml:"queue_info>Queue-List"`
}

// parseQstatf pares the xml output of qstat -f -xml.
func parseQstatf(xmlOut []byte) ([]QstatQueue, error) {
	var qil QstatQueueInfoList
	if err := xml.Unmarshal(xmlOut, &qil); err != nil {
		return qil.QueueList, errors.New("XML QstatQueueInfo List unmarshall error")
	}
	return qil.QueueList, nil
}

// Qstatf executes a qstat -f -xml -q <queueFilter> and retuns
// an array of qstat queueinstances.
func Qstatf(queueFilter string) ([]QstatQueue, error) {
	cmd := exec.Command("qstat", "-f", "-q", queueFilter, "-xml")
	out, errOut := cmd.Output()
	if errOut != nil {
		log.Printf("Could not execute qstat -f -q %s -xml.", queueFilter)
		return nil, errOut
	}
	ql, err := parseQstatf(out)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ql, nil
}
