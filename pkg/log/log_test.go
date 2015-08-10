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
