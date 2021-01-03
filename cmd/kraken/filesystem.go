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

package main

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func fileDetected(filePath, signature string) *Detection {
	log.WithFields(log.Fields{
		"file": filePath,
	}).Warning("DETECTION! Malicious file detected as ", signature)

	detection := NewDetection("filesystem", filePath, "", signature, 0)
	detection.ReportAndStore()

	return detection
}

func filesystemScan() (detections []*Detection) {
	var roots []string
	if *customFileSystemRoot == "" {
		roots = getFileSystemRoots()
	} else {
		roots = []string{*customFileSystemRoot}
	}

	for _, root := range roots {
		if _, err := os.Stat(root); os.IsNotExist(err) {
			log.Error("Cannot scan this folder, it does not appear to exist: ", root)
			continue
		}

		filepath.Walk(root, func(filePath string, fileInfo os.FileInfo, err error) error {
			log.Debug("Scanning file ", filePath)
			matches, _ := yaraScanner.ScanFile(filePath)
			for _, match := range matches {
				detections = append(detections, fileDetected(filePath, match.Rule))
			}

			return nil
		})
	}

	return
}
