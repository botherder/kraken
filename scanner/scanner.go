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

package scanner

import (
	"github.com/hillu/go-yara/v4"
)

// Scanner is an instance of the Yara scanner.
type Scanner struct {
	Available bool
	Rules     *yara.Rules
}

// LoadRules uses Yara's LoadRules to load compiled rules from the specified
// file path.
func (s *Scanner) LoadRules(rulesPath string) error {
	var err error
	s.Rules, err = yara.LoadRules(rulesPath)
	if err != nil {
		return err
	}

	return nil
}

// New returns a new Scanner instance.
func New() Scanner {
	return Scanner{}
}
