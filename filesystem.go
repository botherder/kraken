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
	root := "/home/nex/projects/"
	filepath.Walk(root, func(filePath string, fileHandle os.FileInfo, err error) error {
		log.Debug("Scanning file ", filePath)
		matches, _ := scanner.ScanFile(filePath)
		for _, match := range matches {
			detections = append(detections, fileDetected(filePath, match.Rule))
		}

		return nil
	})

	return
}
