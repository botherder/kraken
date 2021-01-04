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

package scanner

import (
	"errors"
	"os"

	"github.com/hillu/go-yara/v4"
)

// ScanFile scans a file path with the provided Yara rules.
func (s *Scanner) ScanFile(filePath string) (yara.MatchRules, error) {
	var matches yara.MatchRules

	// Check if the scanner is initialized correctly.
	if s.Available == false {
		return matches, errors.New("The scanner is not initialized")
	}

	// Check if the executable file exists.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return matches, errors.New("Scanned file does not exist")
	}

	// Scan the file.
	err := s.Rules.ScanFile(filePath, 0, 60, &matches)
	if err != nil {
		return matches, err
	}

	// Return any results.
	return matches, nil
}

// ScanProc scans a process memory with the provided Yara rules.
func (s *Scanner) ScanProc(pid int) (yara.MatchRules, error) {
	var matches yara.MatchRules

	// Check if the scanner is initialized correctly.
	if s.Available == false {
		return matches, errors.New("The scanner is not initialized")
	}

	// Scan a process memory.
	err := s.Rules.ScanProc(pid, 0, 60, &matches)
	if err != nil {
		return matches, err
	}

	// Return any results.
	return matches, nil
}
