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

package detection

import (
	"fmt"
	"os"
	"path/filepath"

	// log "github.com/Sirupsen/logrus"
	"github.com/botherder/go-savetime/files"
	"github.com/botherder/go-savetime/hashes"
)

const (
	TypeAutorun    = "autorun"
	TypeProcess    = "process"
	TypeFilesystem = "filesystem"
)

// Detection contains the information to report a Yara detection.
type Detection struct {
	Type      string `json:"type"`
	ImagePath string `json:"image_path"`
	ImageName string `json:"image_name"`
	MD5       string `json:"md5"`
	SHA1      string `json:"sha1"`
	SHA256    string `json:"sha256"`
	ProcessID int32  `json:"process_id"`
	Signature string `json:"signature"`
}

// New instantiates a new Detection.
func New(recordType, imagePath, imageName, signature string, pid int32) *Detection {
	md5, _ := hashes.FileMD5(imagePath)
	sha1, _ := hashes.FileSHA1(imagePath)
	sha256, _ := hashes.FileSHA256(imagePath)

	return &Detection{
		Type:      recordType,
		ImagePath: imagePath,
		ImageName: imageName,
		MD5:       md5,
		SHA1:      sha1,
		SHA256:    sha256,
		ProcessID: pid,
		Signature: signature,
	}
}

// // Report sends information on a detection tot he API server.
// func (d *Detection) Report() error {
// 	// First we try to report the detection automatically, so that don't have
// 	// to wait for the next heartbeat iterations in order to report it.
// 	err := apiDetection(d)
// 	// If the report was successful, we don't need to mark it as pending in the
// 	// local database.
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}

// 	return nil
// }

// // Store stores the Detection record in the local SQLite database.
// func (d *Detection) Store(wasReported bool) error {
// 	db := NewDatabase()
// 	err := db.Open()
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}
// 	defer db.Close()

// 	_, err = db.StoreDetection(d, wasReported)
// 	if err != nil {
// 		log.Error(err)
// 		return err
// 	}

// 	return nil
// }

// Backup will keep a copy
func (d *Detection) Backup(folder string) error {
	if _, err := os.Stat(d.ImagePath); err != nil {
		return fmt.Errorf("Binary for detection does not exist at path: %s", d.ImagePath)
	}

	dstPath := filepath.Join(folder, d.SHA1)
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		err = files.Copy(d.ImagePath, dstPath)
		if err != nil {
			return fmt.Errorf("Failed to copy detection binary from %s to %s", d.ImagePath, dstPath)
		}
	}

	return nil
}

// // ReportAndStore is a helper function that will correctly report to the API
// // server and store an entry in the local database. Created to avoid code reuse.
// func (d *Detection) ReportAndStore() error {
// 	// If the report flag was enabled, we send the detection to the API server.
// 	wasReported := false
// 	if *report == true {
// 		err := d.Report()
// 		if err == nil {
// 			wasReported = true
// 		}
// 	}

// 	// If we're running in daemon mode, we store results locally.
// 	if *daemon == true {
// 		d.Store(wasReported)
// 		d.Backup()
// 	}

// 	return nil
// }
