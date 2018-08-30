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
	log "github.com/Sirupsen/logrus"
	"github.com/botherder/go-files"
	"os"
	"path/filepath"
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

// NewDetection instantiates a new Detection.
func NewDetection(recordType, imagePath, imageName, signature string, pid int32) *Detection {
	md5, _ := files.HashFile(imagePath, "md5")
	sha1, _ := files.HashFile(imagePath, "sha1")
	sha256, _ := files.HashFile(imagePath, "sha256")

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

// Report sends information on a detection tot he API server.
func (d *Detection) Report() error {
	// First we try to report the detection automatically, so that don't have
	// to wait for the next heartbeat iterations in order to report it.
	err := apiDetection(d)
	// If the report was successful, we don't need to mark it as pending in the
	// local database.
	if err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}

// Store stores the Detection record in the local SQLite database.
func (d *Detection) Store(wasReported bool) error {
	db := NewDatabase()
	err := db.Open()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer db.Close()

	_, err = db.StoreDetection(d, wasReported)
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Info("Detection stored in database!")

	return nil
}

// Backup will keep a copy
func (d *Detection) Backup() error {
	if _, err := os.Stat(d.ImagePath); err != nil {
		return err
	}

	dstPath := filepath.Join(StorageFiles, d.SHA1)
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		err = files.Copy(d.ImagePath, dstPath)
		if err != nil {
			return err
		}
	}

	return nil
}
