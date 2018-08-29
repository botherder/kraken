// Kraken
// Copyright (C) 2016-2018  Claudio Guarnieri
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package main

import (
	"fmt"
)

// AgentVersion is the... current version of the agent.
const AgentVersion = "1.0"

// URLBaseDomain points to the domain name used to compose all backend URLs.
var URLBaseDomain string = ""

// URLToRules is the URL where the agent will download the rules file.
var URLToRules string = fmt.Sprintf("https://%s/rules", URLBaseDomain)

// URLToRegister is the URL where the agent will register.
var URLToRegister = fmt.Sprintf("https://%s/api/register/", URLBaseDomain)

// URLToHeartbeat is the URL where the agent will beacon to.
var URLToHeartbeat = fmt.Sprintf("https://%s/api/heartbeat/", URLBaseDomain)

// URLToDetection is the URL where the agent reports any detections.
var URLToDetection = fmt.Sprintf("https://%s/api/detection/", URLBaseDomain)

// URLToAutorun is the URL where the agent reports any autoruns.
var URLToAutorun = fmt.Sprintf("https://%s/api/autorun/", URLBaseDomain)
