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
	log := MakeLog("scheduler", "myhostname", os.Stdout)
	log.I("my test %s", "info")
	log.Info("test", "my test %s", "info")
	log.W("my test %s", "warning")
	log.Warning("test", "my test %s", "warning")
	log.E("my test %s", "error")
	log.Error("test", "my test %s", "error")
	log.C("my test %s", "critical")
	log.Critical("test", "my test %s", "critical")
	log.P("my test %s", "profile")
	log.Profile("test", "my test %s", "profile")
}
