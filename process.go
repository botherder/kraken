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
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/process"
)

// This function should terminate processes on all platforms.
func processTerminate(pid int32) bool {
	proc, _ := process.NewProcess(pid)
	err := proc.Terminate()

	return err == nil
}

func processDetected(pid int32, processName, processPath, signature string) *Detection {
	log.WithFields(log.Fields{
		"process": processName,
		"pid":     pid,
	}).Warning("DETECTION! Malicious process detected as ", signature)

	detection := NewDetection("process", processPath, processName, signature, pid)
	detection.ReportAndStore()

	return detection
}

func processScan(pid int32) *Detection {
	// We check we're not scanning ourselves. That'd be awkward.
	if os.Getpid() == int(pid) {
		return nil
	}

	// Fetch details of the process from the pid.
	proc, _ := process.NewProcess(pid)
	procName, _ := proc.Name()
	procExe, err := proc.Exe()

	if err == nil {
		// We check if the process executable exists. If it does, then we scan.
		if _, err := os.Stat(procExe); err == nil {
			log.Debug("Scanning executable for process ", pid, " at path ", procExe)
			matches, _ := scanner.ScanFile(procExe)
			for _, match := range matches {
				return processDetected(pid, procName, procExe, match.Rule)
			}
		}
	}

	// If we got to this point, then it means that the process executable has
	// not been detected. We try now to scan the process memory.
	// NOTE: This might be very prone to false positives. I should probably
	//	   treat this differently.
	log.Debug("Scanning memory of process with PID ", pid)
	matches, _ := scanner.ScanProc(int(pid))
	for _, match := range matches {
		return processDetected(pid, procName, procExe, match.Rule)
	}

	// Nothing found for this process.
	return nil
}

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
