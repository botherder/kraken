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
	"os"
	"path/filepath"
)

// StorageBase is the path to the storage folder for the data of the agent.
var StorageBase = filepath.Join(os.Getenv("HOME"), ".kraken")

// StorageFiles is the path to the folder containing all the copied binaries.
var StorageFiles = filepath.Join(StorageBase, "files")

// StorageConfig is the path to the configuration file.
var StorageConfig = filepath.Join(StorageBase, "config.yaml")

// StorageDatabase is the path to the local SQLite database.
var StorageDatabase = filepath.Join(StorageBase, "database.db")

// StorageRules is the path to a stored rules file.
var StorageRules = filepath.Join(StorageBase, "rules")
