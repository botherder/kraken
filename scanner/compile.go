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
	"errors"
	"os"
	"path"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
	"github.com/hillu/go-yara/v4"
)

// Compile will compile the provided Yara rules in a Rules object.
func Compile(rulesPath string) (*yara.Rules, error) {
	if _, err := os.Stat(rulesPath); os.IsNotExist(err) {
		return nil, errors.New("The specified rules path does not exist")
	}

	compiler, err := yara.NewCompiler()
	if err != nil {
		return nil, err
	}

	rulesStat, _ := os.Stat(rulesPath)
	switch mode := rulesStat.Mode(); {
	case mode.IsDir():
		log.Debug("The specified rules path is a folder, looping through files...")
		err = filepath.Walk(rulesPath, func(filePath string, fileInfo os.FileInfo, err error) error {
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
		log.Debug("Compiling Yara rule ", rulesPath)
		rulesFile, _ := os.Open(rulesPath)
		defer rulesFile.Close()
		err = compiler.AddFile(rulesFile, "")
		if err != nil {
			panic(err)
		}
	}

	// Collect and compile Yara rules.
	return compiler.GetRules()
}
