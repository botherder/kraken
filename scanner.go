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
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/hillu/go-yara"
)

// Scanner is an instance of the Yara scanner.
type Scanner struct {
	Available     bool
	RulesPath     string
	TempRulesPath string
	Rules         *yara.Rules
}

// Compile will compile the provided Yara rules in a Rules object.
func (s *Scanner) Compile() error {
	compiler, err := yara.NewCompiler()
	if err != nil {
		return err
	}

	rulesStat, _ := os.Stat(s.RulesPath)
	switch mode := rulesStat.Mode(); {
	case mode.IsDir():
		log.Debug("The specified rules path is a folder, looping through files...")
		err = filepath.Walk(s.RulesPath, func(filePath string, fileInfo os.FileInfo, err error) error {
			fileName := fileInfo.Name()

			// Check if the file has extension .yar or .yara.
			if (path.Ext(fileName) == ".yar") || (path.Ext(fileName) == ".yara") {
				log.Debug("Adding rule ", filePath)

				// Open the rule file and add it to the Yara compiler.
				rulesFile, _ := os.Open(filePath)
				defer rulesFile.Close()
				err = compiler.AddFile(rulesFile, "")
				if err != nil {
					panic(err)
				}
			}
			return nil
		})
	case mode.IsRegular():
		log.Debug("Compiling Yara rule ", s.RulesPath)
		rulesFile, _ := os.Open(s.RulesPath)
		defer rulesFile.Close()
		err = compiler.AddFile(rulesFile, "")
		if err != nil {
			panic(err)
		}
	}

	// Collect and compile Yara rules.
	s.Rules, err = compiler.GetRules()
	if err != nil {
		return err
	}

	return nil
}

// Init instantiates a new Yara scanner.
func (s *Scanner) Init() error {
	log.Info("Loading Yara rules...")

	// If a customRulesPath is specified, we compile those rules.
	if *customRulesPath != "" {
		if _, err := os.Stat(*customRulesPath); os.IsNotExist(err) {
			return errors.New("The specified rules path does not exist")
		}
		s.RulesPath = *customRulesPath
		return s.Compile()
	}

	// If no customRulesPath is specified, we try to locate a locally stored
	// compiled rules file.
	localRulesPaths := []string{
		filepath.Join(getCwd(), "rules"),
		StorageRules,
	}
	for _, localRulesPath := range localRulesPaths {
		if _, err := os.Stat(localRulesPath); !os.IsNotExist(err) {
			log.Debug("Loading rules file from ", localRulesPath)
			s.RulesPath = localRulesPath
			break
		}
	}

	// If no RulesPath has been selected yet, we try to extract a rules file
	// from the embedded assets.
	if s.RulesPath == "" {
		log.Debug("Loading rules file from embedded assets...")

		// Load embedded rules.
		data, _ := Asset("rules")
		// Create a temporary file for execution.
		// TODO: should this be stored permanently on disk instead?
		tmpRulesFile, _ := ioutil.TempFile("", "agent-")

		// We store the rules file to the temporary location.
		ioutil.WriteFile(tmpRulesFile.Name(), data, 0644)

		// Check if the file was created successfully.
		if _, err := os.Stat(tmpRulesFile.Name()); os.IsNotExist(err) {
			return errors.New("Unable to extract rules from assets")
		}

		// We point RulesPath to the temporary rules file, which should be
		// deleted if we close cleanly.
		s.TempRulesPath = tmpRulesFile.Name()
		s.RulesPath = s.TempRulesPath
	}

	// Load rules.
	var err error
	s.Rules, err = yara.LoadRules(s.RulesPath)
	if err != nil {
		return errors.New("Unable to load rules (maybe corrupted?)")
	}

	return nil
}

// Close deletes the temporary Yara rules file.
func (s *Scanner) Close() {
	// We only delete the temporary rules file extracted from the embedded
	// assets.
	os.Remove(s.TempRulesPath)
}

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
