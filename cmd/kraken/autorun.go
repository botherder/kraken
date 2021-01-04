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
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/botherder/go-autoruns/v2"
	"github.com/botherder/go-savetime/files"
	"github.com/botherder/kraken/api"
	"github.com/botherder/kraken/detection"
	"github.com/botherder/kraken/storage"
)

func autorunDetected(autorun *autoruns.Autorun, signature string) *detection.Detection {
	log.WithFields(log.Fields{
		"type":       autorun.Type,
		"image_path": autorun.ImagePath,
	}).Warning("DETECTION! Malicious autorun detected as ", signature)

	detection := detection.New(detection.TypeAutorun, autorun.ImagePath, autorun.ImageName, signature, 0)

	return detection
}

func autorunStoreInDatabase(autorun *autoruns.Autorun, wasReported bool) {
	// Check if we have already seen this autorun before.
	db := NewDatabase()
	err := db.Open()
	if err != nil {
		log.Error("Failed to store autorun record in local database: ", err)
		return
	}
	defer db.Close()

	isStored, _ := db.IsAutorunStored(autorun)
	if isStored == false {
		// If not, we store it in the local database with the appropriate marking.
		db.StoreAutorun(autorun, wasReported)

		log.WithFields(log.Fields{
			"path":      autorun.ImagePath,
			"arguments": autorun.Arguments,
		}).Debug("New autorun identified and stored in local database!")
	}
}

func autorunScan(autorun *autoruns.Autorun) *detection.Detection {
	// We want to report autorun records even if they were not detected as malicous.
	wasReported := false
	if *flagReport == true {
		client := api.New(cfg.BaseDomain, cfg.MachineID)
		err := client.ReportAutorun(autorun)
		if err != nil {
			log.Error("Failed to report autorun record: ", err)
		} else {
			log.Debug("Autorun record reported to server!")
			wasReported = true
		}
	}

	// We store data only if we run in daemon mode.
	if *flagDaemon == true {
		// Store autorun in database if necessary.
		autorunStoreInDatabase(autorun, wasReported)

		// We backup the autorun file.
		if _, err := os.Stat(autorun.ImagePath); err == nil {
			dstPath := filepath.Join(storage.StorageFiles, autorun.SHA1)
			if _, err := os.Stat(dstPath); os.IsNotExist(err) {
				files.Copy(autorun.ImagePath, dstPath)
			}
		}
	}

	// Lastly, we scan it.
	if _, err := os.Stat(autorun.ImagePath); err == nil {
		log.Debug("Scanning ", autorun.Type, " autorun at path ", autorun.ImagePath)
		matches, _ := yaraScanner.ScanFile(autorun.ImagePath)
		for _, match := range matches {
			return autorunDetected(autorun, match.Rule)
		}
	}

	return nil
}
