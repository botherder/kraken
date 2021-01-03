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

package api

import (
	"fmt"

	"github.com/botherder/go-autoruns/v2"
	"github.com/go-resty/resty/v2"
)

// Report an autorun.
func (a *API) ReportAutorun(record *autoruns.Autorun) error {
	client := resty.New()
	response, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(record).
		Post(a.Config.URLToAutorun)

	// Check if the request failed.
	if err != nil {
		return fmt.Errorf("Unable to send autorun record to API server: %s", err)
	}

	// Check if the response wasn't right.
	if response.StatusCode() != 200 {
		return fmt.Errorf("Unable to send autorun record to API server: we received response code %d", response.StatusCode())
	}

	return nil
}
