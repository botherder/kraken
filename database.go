// Kraken
// Copyright (C) 2016-2020  Claudio Guarnieri
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
	"database/sql"
	"fmt"

	"github.com/botherder/go-autoruns"
	_ "github.com/mattn/go-sqlite3"
)

const tableAutoruns = `CREATE TABLE IF NOT EXISTS autoruns (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	record_type VARCHAR,
	image_path VARCHAR NOT NULL,
	image_name VARCHAR,
	arguments VARCHAR,
	sha1 VARCHAR NOT NULL,
	sha256 VARCHAR NOT NULL,
	added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	reported INT DEFAULT 0
)`

const tableDetections = `CREATE TABLE IF NOT EXISTS detections (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	detection_type VARCHAR NOT NULL,
	image_path VARCHAR,
	image_name VARCHAR,
	sha1 VARCHAR,
	sha256 VARCHAR,
	process_id INT,
	signature VARCHAR NOT NULL,
	detected_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	reported INT DEFAULT 0
)`

// Database provides access to the internal SQLite database.
type Database struct {
	db *sql.DB
}

// NewDatabase instantiates a new Database to the local SQLite.
func NewDatabase() *Database {
	return &Database{}
}

// Open initializes the database, if it does not exist yet it creates the necessary tables.
func (d *Database) Open() error {
	var err error
	d.db, err = sql.Open("sqlite3", StorageDatabase)
	if err != nil {
		return fmt.Errorf("Unable to open database: %s", err.Error())
	}
	// defer d.db.Close()

	d.db.Exec(tableAutoruns)
	d.db.Exec(tableDetections)

	return nil
}

// Close closes the connection to the SQLite database.
func (d *Database) Close() {
	d.db.Close()
}

// IsAutorunStored checks if a given Autorun record was already stored in the database.
func (d *Database) IsAutorunStored(record *autoruns.Autorun) (bool, error) {
	statement, err := d.db.Prepare("SELECT COUNT(*) as count FROM autoruns WHERE record_type = ? AND image_path = ? AND arguments = ?;")
	if err != nil {
		return false, fmt.Errorf("Unable to prepare isAutorunStored query: %s", err.Error())
	}
	defer statement.Close()

	rows, err := statement.Query(record.Type, record.ImagePath, record.Arguments)
	var count int
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return false, fmt.Errorf("Unable to check if record exists in IsAutorunStored: %s", err.Error())
		}
	}

	if count > 0 {
		return true, nil
	}

	return false, nil
}

// StoreAutorun creates a record in the SQLite database for a given Autorun.
func (d *Database) StoreAutorun(record *autoruns.Autorun, reported bool) (int64, error) {
	var reportedValue int
	if reported {
		reportedValue = 1
	}

	statement, err := d.db.Prepare("INSERT INTO autoruns (record_type, image_path, image_name, arguments, sha1, sha256, reported) VALUES (?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return 0, fmt.Errorf("Unable to prepare StoreAutorun query: %s", err.Error())
	}

	result, err := statement.Exec(record.Type, record.ImagePath, record.ImageName, record.Arguments, record.SHA1, record.SHA256, reportedValue)
	if err != nil {
		return 0, fmt.Errorf("Unable to insert new record: %s", err.Error())
	}

	return result.LastInsertId()
}

// StoreDetection creates a record in the SQLite database for a malware detection.
func (d *Database) StoreDetection(record *Detection, reported bool) (int64, error) {
	var reportedValue int
	if reported {
		reportedValue = 1
	}

	statement, err := d.db.Prepare("INSERT INTO detections (detection_type, image_path, image_name, sha1, sha256, process_id, signature, reported) VALUES (?, ?, ?, ?, ?, ?, ?, ?);")
	if err != nil {
		return 0, fmt.Errorf("Unable to prepare StoreDetection query: %s", err.Error())
	}

	result, err := statement.Exec(record.Type, record.ImagePath, record.ImageName, record.SHA1, record.SHA256, record.ProcessID, record.Signature, reportedValue)
	if err != nil {
		return 0, fmt.Errorf("Unable to insert new record: %s", err.Error())
	}

	return result.LastInsertId()
}
