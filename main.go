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

	log "github.com/Sirupsen/logrus"
	"github.com/botherder/go-autoruns"
	"github.com/mattn/go-colorable"
	"github.com/shirou/gopsutil/process"
	flag "github.com/spf13/pflag"
)

var (
	// This is our Yara scanner.
	scanner Scanner
	// This is a flag to enable debug logging.
	debug *bool
	// This is a flag to determine whether to execute as a permanent agent or not.
	// In case this is false, we just scan active processes and exit.
	daemon *bool
	// This is a flag to enable remote reporting to API server.
	report *bool
	// This is a domain to the backend specified from command-line.
	customBaseDomain *string
	// This is a folder to be scanned instead of the default.
	customFileSystemRoot *string
	// This is a flag to disable process scanning.
	noProcessScan *bool
	// This is a flag to disable autorun scanning.
	noAutorunsScan *bool
	// This is a flag to disable filesystem scanning.
	noFileSystemScan *bool
)

func initArguments() {
	debug = flag.Bool("debug", false, "Enable debug logs")
	report = flag.Bool("report", false, "Enable reporting of events to the backend")
	daemon = flag.Bool("daemon", false, "Enable daemon mode (this will also enable the report flag)")
	customBaseDomain = flag.String("backend", "", "Specify a particular hostname to the backend to connect to (overrides the default)")
	customFileSystemRoot = flag.String("folder", "", "Specify a particular folder to be scanned (overrides the default full filesystem)")
	noProcessScan = flag.Bool("no-process", false, "Disable scanning of running processes")
	noAutorunsScan = flag.Bool("no-autoruns", false, "Disable scanning of autoruns")
	noFileSystemScan = flag.Bool("no-filesystem", false, "Disable scanning of filesystem")
	flag.Parse()

	// If we're running in daemon mode, we enable the report flag too.
	// TODO: Need to review this choice. We might not necessarily want that.
	if *daemon == true {
		*report = true
	}
}

func initLogging() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())

	if *debug {
		log.SetLevel(log.DebugLevel)
	}
}

func initStorage() {
	// We create the folder only if we're running in daemon mode.
	if *daemon == true {
		// This should create StorageBase and StorageFiles.
		if _, err := os.Stat(StorageFiles); os.IsNotExist(err) {
			os.MkdirAll(StorageFiles, 0777)
		}
	}
}

// This function contains just the preliminary actions.
func init() {
	// Parse arguments.
	initArguments()
	// Initialize debugging.
	initLogging()
	// Initialize storage.
	initStorage()
	// Initialize configuration.
	initConfig()

	log.Debug("This machine is identified as ", config.MachineID)
	log.Debug("URLBaseDomain: ", config.URLBaseDomain)
	log.Debug("URLToRules: ", config.URLToRules)
	log.Debug("URLToRegister: ", config.URLToRegister)
	log.Debug("URLToDetection: ", config.URLToDetection)
	log.Debug("URLToAutorun: ", config.URLToAutorun)

	// We register to the backend only if report is enabled.
	if *report == true {
		// Register to the API server.
		err := apiRegister()
		if err != nil {
			log.Error("API registration failed: ", err.Error())
		}

		// Try to send an heartbeat.
		err = apiHeartbeat()
		if err != nil {
			log.Error("Unable to send heartbeat to server: ", err.Error())
		}
	}
}

func main() {
	// Initialize the Yara scanner.
	err := scanner.Init()
	if err != nil {
		scanner.Available = false
	} else {
		scanner.Available = true
	}
	defer scanner.Close()

	// We store here the list of detections.
	var detections []*Detection
	// Empty list of pids.
	var pids []int32

	// We do a first scan of running processes.
	if *noProcessScan == false {
		log.Info("Scanning running processes...")
		pids, _ = process.Pids()
		for _, pid := range pids {
			detections = append(detections, processScan(pid)...)
		}
	}

	// We scan the running autoruns.
	if *noAutorunsScan == false {
		log.Info("Scanning autoruns...")
		autoruns := autoruns.Autoruns()
		for _, autorun := range autoruns {
			detection := autorunScan(autorun)
			if detection != nil {
				detections = append(detections, detection)
			}
		}
	}

	// Now we do a scan of the file system if required.
	if *noFileSystemScan == false {
		log.Info("Scanning the filesystem (this can take several minutes)...")
		detections = append(detections, filesystemScan()...)
	}

	// Now we tell the results.
	if len(detections) > 0 {
		log.Error("Some malicious artifacts have been detected on this system:")
		for _, detection := range detections {
			log.Error("Found detection for ", detection.Signature)
		}
	} else {
		log.Info("GOOD! Nothing detected!")
	}

	// If by command-line it was instructed to run in daemon mode, then
	// we start the process watch.
	if *daemon == true {
		// Start process monitor.
		go processWatch(pids)
		// Start autoruns monitor.
		go autorunWatch()
		// Start heartbeat manager.
		heartbeatManager()
	} else {
		log.Info("Press Enter to finish ...")
		var b = make([]byte, 1)
		os.Stdin.Read(b)
	}
}
