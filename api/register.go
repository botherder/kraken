// Kraken
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

	"github.com/botherder/kraken/profile"
	"github.com/go-resty/resty/v2"
)

// Registration contains the information to register to the API.
type Registration struct {
	Identifier      string `json:"identifier"`
	UserName        string `json:"user_name"`
	ComputerName    string `json:"computer_name"`
	OperatingSystem string `json:"operating_system"`
	Version         string `json:"version"`
}

// Register to the API server.
func (a *API) Register() error {
	registration := Registration{
		Identifier:      a.Config.MachineID,
		UserName:        profile.GetUsername(),
		ComputerName:    profile.GetComputerName(),
		OperatingSystem: profile.GetOperatingSystem(),
	}

	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(registration).
		Post(a.Config.URLToRegister)

	// Check if request failed.
	if err != nil {
		return fmt.Errorf("Unable to register to API server: %s", err.Error())
	}

	// Check if the response wasn't right.
	if response.StatusCode() != 200 {
		return fmt.Errorf("Unable to register to API server: we received response code %d", response.StatusCode())
	}

	return nil
}
