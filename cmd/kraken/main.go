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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/botherder/go-autoruns/v2"
	"github.com/botherder/go-savetime/runtime"
	"github.com/botherder/kraken/api"
	"github.com/botherder/kraken/config"
	"github.com/botherder/kraken/detection"
	"github.com/botherder/kraken/scanner"
	"github.com/botherder/kraken/storage"
	"github.com/mattn/go-colorable"
	"github.com/shirou/gopsutil/process"
	flag "github.com/spf13/pflag"
)

var (
	apiClient *api.API
	cfg       *config.Config
	// This is our Yara yaraScanner.
	yaraScanner scanner.Scanner
	// This is a flag to enable debug logging.
	debug *bool
	// This is a flag to determine whether to execute as a permanent agent or not.
	// In case this is false, we just scan active processes and exit.
	daemon *bool
	// This is a flag to enable remote reporting to API server.
	report *bool
	// This is a domain to the backend specified from command-line.
	flagBaseDomain *string
	// This is a folder to be scanned instead of the default.
	flagScanFolder *string
	// This is a file or folder path containing the Yara rules to use.
	flagRulesPath *string
	// This is a flag to disable process scanning.
	flagNoProcessScan *bool
	// This is a flag to disable autorun scanning.
	flagNoAutorunsScan *bool
	// This is a flag to disable filesystem scanning.
	flagNoFilesystemScan *bool
)

func initArguments() {
	debug = flag.Bool("debug", false, "Enable debug logs")
	report = flag.Bool("report", false, "Enable reporting of events to the backend")
	daemon = flag.Bool("daemon", false, "Enable daemon mode (this will also enable the report flag)")
	flagBaseDomain = flag.String("backend", "", "Specify a particular hostname to the backend to connect to (overrides the default)")
	flagScanFolder = flag.String("folder", "", "Specify a particular folder to be scanned (overrides the default full filesystem)")
	flagRulesPath = flag.String("rules", "", "Specify a particular path to a file or folder containing the Yara rules to use")
	flagNoProcessScan = flag.Bool("no-process", false, "Disable scanning of running processes")
	flagNoAutorunsScan = flag.Bool("no-autoruns", false, "Disable scanning of autoruns")
	flagNoFilesystemScan = flag.Bool("no-filesystem", false, "Disable scanning of filesystem")
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
		// This should create storage.StorageBase and storage.StorageFiles.
		if _, err := os.Stat(storage.StorageFiles); os.IsNotExist(err) {
			os.MkdirAll(storage.StorageFiles, 0777)
		}
	}
}

func initScanner() error {
	log.Info("Loading Yara rules...")

	var err error
	yaraScanner = scanner.New()

	// If a flagRulesPath is specified, we compile those rules.
	if *flagRulesPath != "" {
		yaraScanner.Rules, err = scanner.Compile(*flagRulesPath)
		if err != nil {
			yaraScanner.Available = false
			return fmt.Errorf("Unable to compile custom Yara rules at path %s: %s", flagRulesPath, err)
		}

		yaraScanner.Available = true

		// NOTE: In this case we return either way, because if a suer specified
		//       a custom rules path, we should expect they wouldn't want to
		//       continue if an error occurred.
		return nil
	}

	// If no flagRulesPath is specified, we try to locate a locally stored
	// compiled rules file.
	localRulesPaths := []string{
		filepath.Join(runtime.GetExecutableDirectory(), "rules"),
		storage.StorageRules,
	}
	for _, localRulesPath := range localRulesPaths {
		if _, err := os.Stat(localRulesPath); !os.IsNotExist(err) {
			err = yaraScanner.LoadRules(localRulesPath)
			if err != nil {
				log.Error("Unable to load compiled rules file at path %s: %s", localRulesPath, err)
				yaraScanner.Available = false
			} else {
				yaraScanner.Available = true
				// If we successfully loaded some rules, we can stop here.
				return nil
			}
		}
	}

	// If no rules have been selected yet, we try to extract a rules file
	// from the embedded assets.
	log.Debug("Trying to load rules file from embedded assets...")

	// Load embedded rules.
	data, _ := Asset("rules")

	// Create a temporary file for execution.
	// TODO: should this be stored permanently on disk instead?
	tmpRulesFile, err := ioutil.TempFile("", "agent-")
	if err != nil {
		return fmt.Errorf("Unable to temporarily store rules: %s", err)
	}
	defer tmpRulesFile.Close()

	// We store the rules file to the temporary location.
	tmpRulesFile.Write(data)

	err = yaraScanner.LoadRules(tmpRulesFile.Name())
	if err != nil {
		return fmt.Errorf("Unable to load rules from embedded assets: %s", err)
	}

	return nil
}

// This function contains just the preliminary actions.
func init() {
	// Parse arguments.
	initArguments()
	// Initialize debugging.
	initLogging()
	// Initialize storage.
	initStorage()
	// Initialize Yara yaraScanner.
	initScanner()

	cfg = config.New(*flagBaseDomain, DefaultBaseDomain)

	log.Debug("This machine is identified as ", cfg.MachineID)
	log.Debug("The agent is going to communicate to: ", cfg.BaseDomain)

	// We register to the backend only if report is enabled.
	if *report == true {
		apiClient = api.New(cfg.BaseDomain, cfg.MachineID)

		// Register to the API server.
		err := apiClient.Register()
		if err != nil {
			log.Error("API registration failed: ", err)
		}

		// Try to send an heartbeat.
		err = apiClient.Heartbeat()
		if err != nil {
			log.Error("Unable to send heartbeat to server: ", err)
		}
	}
}

func main() {
	// We store here the list of detections.
	var detections []*detection.Detection

	// Empty list of pids.
	var pids []int32
	// We do a first scan of running processes.
	if !*flagNoProcessScan {
		log.Info("Scanning running processes...")
		pids, _ = process.Pids()
		for _, pid := range pids {
			detections = append(detections, processScan(pid)...)
		}
	}

	// We scan the running autoruns.
	if !*flagNoAutorunsScan {
		log.Info("Scanning autoruns...")
		for _, autorun := range autoruns.GetAllAutoruns() {
			detection := autorunScan(autorun)
			if detection != nil {
				detections = append(detections, detection)
			}
		}
	}

	// Now we do a scan of the file system if required.
	if !*flagNoFilesystemScan {
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
		fmt.Scanln()
	}
}
