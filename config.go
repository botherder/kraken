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
	URLToRules     string
	URLToRegister  string
	URLToHeartbeat string
	URLToDetection string
	URLToAutorun   string
}

// This is our configuration file.
var config Config

func configInit() {
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
	viper.SetDefault("url_rules", URLToRules)
	viper.SetDefault("url_register", URLToRegister)
	viper.SetDefault("url_heartbeat", URLToHeartbeat)
	viper.SetDefault("url_detection", URLToDetection)
	viper.SetDefault("url_autorun", URLToAutorun)

	// Save configuration values to our Config instance.
	config.MachineID = viper.GetString("machine_id")
	config.URLToRules = viper.GetString("url_rules")
	config.URLToRegister = viper.GetString("url_register")
	config.URLToHeartbeat = viper.GetString("url_heartbeat")
	config.URLToDetection = fmt.Sprintf("%s%s/", viper.GetString("url_detection"), config.MachineID)
	config.URLToAutorun = fmt.Sprintf("%s%s/", viper.GetString("url_autorun"), config.MachineID)

	log.Info("This machine is identified as ", config.MachineID)
	log.Debug("URLBaseDomain: ", URLBaseDomain)
	log.Debug("URLToRules: ", config.URLToRules)
	log.Debug("URLToRegister: ", config.URLToRegister)
	log.Debug("URLToDetection: ", config.URLToDetection)
	log.Debug("URLToAutorun: ", config.URLToAutorun)

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
