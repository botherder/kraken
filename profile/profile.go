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

package profile

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/user"

	"github.com/botherder/go-savetime/hashes"
	"github.com/jaypipes/ghw"
	"github.com/matishsiao/goInfo"
)

// GetMachineID attempts to generate a serial number, first by trying to
// get a disk serial number. If that fails, we use the network card.
func GetMachineID() string {
	serial, err := GetDiskSerialNumber()
	if err == nil {
		id, _ := hashes.StringSHA1(serial)
		return id
	}

	mac, err := GetMacAddress()
	if err == nil {
		id, _ := hashes.StringSHA1(mac)
		return id
	}

	return ""
}

// GetDiskSerialNumber returns the first
func GetDiskSerialNumber() (string, error) {
	block, err := ghw.Block()
	if err != nil {
		return "", err
	}

	for _, disk := range block.Disks {
		if disk.SerialNumber == "unknown" {
			continue
		}

		return disk.SerialNumber, nil
	}

	return "", errors.New("No available disk found")
}

// GetMacAddress returns the mac address of the first available network
// interface.
func GetMacAddress() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		mac := iface.HardwareAddr.String()
		return mac, nil
	}

	return "", errors.New("No available network interface found")
}

// Get current username.
func GetUsername() string {
	userObject, err := user.Current()
	if err != nil {
		return ""
	}

	return userObject.Username
}

// Get computer name.
func GetComputerName() string {
	hostname, _ := os.Hostname()
	return hostname
}

// Get some accurate version of the operating system.
func GetOperatingSystem() string {
	gi := goInfo.GetInfo()
	return fmt.Sprintf("%s %s", gi.OS, gi.Core)
}
