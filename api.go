// Kraken
// Copyright (C) 2016-2020  Claudio Guarnieri
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

	"github.com/botherder/go-autoruns/v2"
	"gopkg.in/resty.v0"
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
func apiRegister() error {
	registration := Registration{
		Identifier:      config.MachineID,
		UserName:        getUserName(),
		ComputerName:    getComputerName(),
		OperatingSystem: getOperatingSystem(),
		Version:         AgentVersion,
	}

	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(registration).
		Post(config.URLToRegister)

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

// Sends an heartbeat to the API server.
func apiHeartbeat() error {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(fmt.Sprintf(`{"identifier":"%s"}`, config.MachineID)).
		Post(config.URLToHeartbeat)

	// Check if request failed.
	if err != nil {
		return fmt.Errorf("Unable to send heartbeat to API server: %s", err.Error())
	}

	// Check if the response wasn't right.
	if response.StatusCode() != 200 {
		return fmt.Errorf("Unable to send heartbeat to API server: we received response code %d", response.StatusCode())
	}

	return nil
}

// Report a detection.
func apiDetection(record *Detection) error {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(record).
		Post(config.URLToDetection)

	// Check if the request failed.
	if err != nil {
		return fmt.Errorf("Unable to send detection to API server: %s", err.Error())
	}

	// Check if the response wasn't right.
	if response.StatusCode() != 200 {
		return fmt.Errorf("Unable to send detection to API server: we received response code %d", response.StatusCode())
	}

	return nil
}

// Report an autorun.
func apiAutorun(record *autoruns.Autorun) error {
	response, err := resty.R().
		SetHeader("Content-Type", "application/json").
		SetBody(record).
		Post(config.URLToAutorun)

	// Check if the request failed.
	if err != nil {
		return fmt.Errorf("Unable to send autorun record to API server: %s", err.Error())
	}

	// Check if the response wasn't right.
	if response.StatusCode() != 200 {
		return fmt.Errorf("Unable to send autorun record to API server: we received response code %d", response.StatusCode())
	}

	return nil
}
