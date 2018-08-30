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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
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

// This is our configuration.
var config Config

func initConfig() {
	var baseDomain string
	if *customBaseDomain != "" {
		log.Debug("I was provided a custom backend: ", *customBaseDomain)
		baseDomain = *customBaseDomain
	} else {
		log.Debug("I am going to use the default backend: ", DefaultBaseDomain)
		baseDomain = DefaultBaseDomain
	}

	// Get folder and file name of standard config file.
	fileName := filepath.Base(StorageConfig)

	log.Info("Looking for configuration file with name ", fileName)

	// Specify configuration file properties.
	viper.SetConfigName(strings.Split(fileName, ".")[0])
	viper.SetConfigType("yaml")
	viper.AddConfigPath(StorageBase)
	viper.AddConfigPath(".")

	// Try to read the config file.
	readError := viper.ReadInConfig()
	if readError != nil {
		log.Info("No configuration file found, generating a default one...")
	}

	// Just in case there is no config file, we set defaults.
	viper.SetDefault("machine_id", getMachineID())
	viper.SetDefault("base_domain", baseDomain)

	// Save configuration values to our Config instance.
	config.MachineID = viper.GetString("machine_id")
	config.URLBaseDomain = viper.GetString("base_domain")
	config.URLToRules = fmt.Sprintf("https://%s/rules", viper.GetString("base_domain"))
	config.URLToRegister = fmt.Sprintf("https://%s/api/register/", viper.GetString("base_domain"))
	config.URLToHeartbeat = fmt.Sprintf("https://%s/api/heartbeat/", viper.GetString("base_domain"))
	config.URLToDetection = fmt.Sprintf("https://%s/api/detection/%s/", viper.GetString("base_domain"), viper.GetString("machine_id"))
	config.URLToAutorun = fmt.Sprintf("https://%s/api/autorun/%s/", viper.GetString("base_domain"), viper.GetString("machine_id"))

	// Write a new config file if none was found.
	// We actually write it to disk only if we're running in daemon mode.
	if readError != nil && *daemon == true {
		// err := viper.SafeWriteConfigAs(StorageConfig)
		err := viper.WriteConfigAs(StorageConfig)
		if err != nil {
			log.Fatal("Unable to write a new default configuration to file: ", err.Error())
		}
		log.Info("New default configuration file written to ", StorageConfig)
	}
}
