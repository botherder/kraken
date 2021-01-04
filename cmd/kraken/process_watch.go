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

package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/process"
)

func processWatch(oldPids []int32) {
	log.Info("Starting process monitor...")

	for {
		// Grab a list of pids.
		pids, _ := process.Pids()

		// Loop through the current pids.
		for _, pid := range pids {
			// Mark if the pid has been observed in the previous iteration.
			var previouslySeen = false
			// Loop through the pids from the previous iteration.
			for _, oldPid := range oldPids {
				// Check if the pid is already known.
				if pid == oldPid {
					previouslySeen = true
				}
			}

			// If the current pid is not known from the previous iteration,
			// then we launch a new scan goroutine.
			if previouslySeen == false {
				log.Info("New process to scan with PID ", pid)
				go processScan(pid)
			}
		}

		// Update the pids last iteration list.
		oldPids = pids

		// Zzz.
		time.Sleep(time.Second)
	}
}
