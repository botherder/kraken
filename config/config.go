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

package config

import (
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/botherder/kraken/profile"
	"github.com/botherder/kraken/storage"
	"github.com/spf13/viper"
)

// Config contains all the information for reporting back.
type Config struct {
	MachineID  string
	BaseDomain string
}

func New(customBaseDomain, defaultBaseDomain string) *Config {
	var baseDomain string
	if customBaseDomain != "" {
		baseDomain = customBaseDomain
	} else {
		baseDomain = defaultBaseDomain
	}

	// Get folder and file name of standard config file.
	fileName := filepath.Base(storage.StorageConfig)

	log.Info("Looking for configuration file with name ", fileName)

	// Specify configuration file properties.
	viper.SetConfigName(strings.Split(fileName, ".")[0])
	viper.SetConfigType("yaml")
	viper.AddConfigPath(storage.StorageBase)
	viper.AddConfigPath(".")

	// Try to read the config file.
	readError := viper.ReadInConfig()
	if readError != nil {
		log.Info("No configuration file found, generating a default one...")
	}

	// Just in case there is no config file, we set defaults.
	viper.SetDefault("machine_id", profile.GetMachineID())
	viper.SetDefault("base_domain", baseDomain)

	return &Config{
		MachineID:  viper.GetString("machine_id"),
		BaseDomain: viper.GetString("base_domain"),
	}
}

func (c *Config) WriteToFile(configPath string) {
	// err := viper.SafeWriteConfigAs(configPath)
	err := viper.WriteConfigAs(configPath)
	if err != nil {
		log.Fatal("Unable to write a new default configuration to file: ", err)
	}
	log.Info("New default configuration file written to ", configPath)
}
