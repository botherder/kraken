// This file is part of Kraken (https://github.com/botherder/kraken)
// Copyright (C) 2016-2021  Claudio Guarnieri
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

package api

import (
	"fmt"
)

type API struct {
	MachineID      string
	BaseDomain     string
	URLToRules     string
	URLToRegister  string
	URLToHeartbeat string
	URLToDetection string
	URLToAutorun   string
}

func New(baseDomain, machineID string) *API {
	return &API{
		MachineID:      machineID,
		BaseDomain:     baseDomain,
		URLToRules:     fmt.Sprintf("https://%s/rules", baseDomain),
		URLToRegister:  fmt.Sprintf("https://%s/api/register/", baseDomain),
		URLToHeartbeat: fmt.Sprintf("https://%s/api/heartbeat/", baseDomain),
		URLToDetection: fmt.Sprintf("https://%s/api/detection/%s/", baseDomain, machineID),
		URLToAutorun:   fmt.Sprintf("https://%s/api/autorun/%s/", baseDomain, machineID),
	}
}
