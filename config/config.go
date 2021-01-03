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

package config

import (
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/botherder/kraken/profile"
	"github.com/botherder/kraken/storage"
	"github.com/spf13/viper"
)

// Config contains all the information for reporting back.
type Config struct {
	MachineID      string
	URLBaseDomain  string
	URLToRules     string
	URLToRegister  string
	URLToHeartbeat string
	URLToDetection string
	URLToAutorun   string
}

func New(customBaseDomain, defaultBaseDomain string) *Config {
	var baseDomain string
	if customBaseDomain != "" {
		log.Debug("I was provided a custom backend: ", customBaseDomain)
		baseDomain = customBaseDomain
	} else {
		log.Debug("I am going to use the default backend: ", defaultBaseDomain)
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

	cfg := Config{}

	// Save configuration values to our Config instance.
	cfg.MachineID = viper.GetString("machine_id")
	cfg.URLBaseDomain = viper.GetString("base_domain")
	cfg.URLToRules = fmt.Sprintf("https://%s/rules", viper.GetString("base_domain"))
	cfg.URLToRegister = fmt.Sprintf("https://%s/api/register/", viper.GetString("base_domain"))
	cfg.URLToHeartbeat = fmt.Sprintf("https://%s/api/heartbeat/", viper.GetString("base_domain"))
	cfg.URLToDetection = fmt.Sprintf("https://%s/api/detection/%s/", viper.GetString("base_domain"), viper.GetString("machine_id"))
	cfg.URLToAutorun = fmt.Sprintf("https://%s/api/autorun/%s/", viper.GetString("base_domain"), viper.GetString("machine_id"))

	return &cfg
}

func (c *Config) Write() {
	// err := viper.SafeWriteConfigAs(storage.StorageConfig)
	err := viper.WriteConfigAs(storage.StorageConfig)
	if err != nil {
		log.Fatal("Unable to write a new default configuration to file: ", err.Error())
	}
	log.Info("New default configuration file written to ", storage.StorageConfig)
}
