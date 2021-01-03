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
	"time"

	log "github.com/Sirupsen/logrus"
)

func heartbeatManager() {
	log.Info("Starting heartbeat...")

	ticker := time.NewTicker(time.Minute * 30).C

	for {
		select {
		case <-ticker:
			log.Info("Sending heartbeat to server...")
			err := apiHeartbeat()
			if err != nil {
				log.Error("Unable to send heartbeat to server: ", err.Error())
			}
		}
	}
}
