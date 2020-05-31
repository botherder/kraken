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
	"net"
	"os"
	"os/user"

	"github.com/botherder/go-files"
	"github.com/matishsiao/goInfo"
)

// For a machine ID we use a SHA1 of the first Mac address we find.
// TODO: should rather use the disk serial number.
func getMachineID() string {
	ifaces, _ := net.Interfaces()
	for _, iface := range ifaces {
		mac := iface.HardwareAddr.String()
		if len(mac) == 17 {
			hash, _ := files.HashString(mac, "sha1")
			return hash
		}
	}

	return ""
}

// Get current username.
func getUserName() string {
	userObject, err := user.Current()
	if err != nil {
		return ""
	}

	return userObject.Username
}

// Get computer name.
func getComputerName() string {
	hostname, _ := os.Hostname()
	return hostname
}

// Get some accurate version of the operating system.
func getOperatingSystem() string {
	gi := goInfo.GetInfo()
	return fmt.Sprintf("%s %s", gi.OS, gi.Core)
}
