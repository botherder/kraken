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
	"errors"
	"io/ioutil"
	"os"

	"github.com/hillu/go-yara"
)

// Scanner is an instance of the Yara scanner.
type Scanner struct {
	Available bool
	RulesPath string
	Rules     *yara.Rules
}

// Init instantiates a new Yara scanner.
func (s *Scanner) Init() error {
	var err error

	// Load embedded rules.
	// TODO: do this only if there is no local rules files on the disk.
	data, _ := Asset("rules")
	// Create a temporary file for execution.
	// TODO: should this be stored permanently on disk instead?
	rulesObject, _ := ioutil.TempFile("", "agent-")

	// We store the rules file to the temporary location.
	ioutil.WriteFile(rulesObject.Name(), data, 0644)

	// Check if the file was created successfully.
	if _, err := os.Stat(rulesObject.Name()); os.IsNotExist(err) {
		return errors.New("Unable to extract rules from assets")
	}

	// Instantiate scanner.
	s.RulesPath = rulesObject.Name()

	// Load rules.
	s.Rules, err = yara.LoadRules(s.RulesPath)
	if err != nil {
		return errors.New("Unable to load rules (maybe corrupted?)")
	}

	return nil
}

// Close deletes the temporary Yara rules file.
func (s *Scanner) Close() {
	// Delete rules file.
	os.Remove(s.RulesPath)
}

// ScanFile scans a file path with the provided Yara rules.
func (s *Scanner) ScanFile(filePath string) ([]yara.MatchRule, error) {
	var matches []yara.MatchRule

	// Check if the scanner is initialized correctly.
	if s.Available == false {
		return matches, errors.New("The scanner is not initialized")
	}

	// Check if the executable file exists.
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return matches, errors.New("Scanned file does not exist")
	}

	// Scan the file.
	matches, _ = s.Rules.ScanFile(filePath, 0, 60)

	// Return any results.
	return matches, nil
}

// ScanProc scans a process memory with the provided Yara rules.
func (s *Scanner) ScanProc(pid int) ([]yara.MatchRule, error) {
	var matches []yara.MatchRule

	// Check if the scanner is initialized correctly.
	if s.Available == false {
		return matches, errors.New("The scanner is not initialized")
	}

	// Scan a process memory.
	matches, _ = s.Rules.ScanProc(pid, 0, 60)

	// Return any results.
	return matches, nil
}
