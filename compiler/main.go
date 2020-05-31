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
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/hillu/go-yara"
)

// var compiler *yara.Compiler

func main() {
	// The argument is the path to the rule or a folder containing rules.
	rulesPath := os.Args[1]

	// Check if the path exists.
	if _, err := os.Stat(rulesPath); os.IsNotExist(err) {
		panic(err)
	}

	// Instantiate the Yara compiler.
	compiler, err := yara.NewCompiler()
	if err != nil {
		panic(err)
	}

	rulesStat, _ := os.Stat(rulesPath)
	switch mode := rulesStat.Mode(); {
	// Check if the path is a folder...
	case mode.IsDir():
		fmt.Println("[rules-compiler] The specified path is a folder, looping through files...")
		err = filepath.Walk(rulesPath, func(filePath string, fileObj os.FileInfo, err error) error {
			// Get the file name.
			fileName := fileObj.Name()

			// Check if the file has extension .yar or .yara.
			if (path.Ext(fileName) == ".yar") || (path.Ext(fileName) == ".yara") {
				fmt.Println("[rules-compiler] Adding rule", filePath)

				// Open the rule file and add it to the Yara compiler.
				rulesObj, _ := os.Open(filePath)
				defer rulesObj.Close()
				compiler.AddFile(rulesObj, "")
			}
			return nil
		})
	// Check if it is a file instead...
	case mode.IsRegular():
		fmt.Println("[rules-compiler] Compiling Yara rule", rulesPath)
		// Open the rule file and add it to the Yara compiler.
		rulesObj, _ := os.Open(rulesPath)
		defer rulesObj.Close()
		compiler.AddFile(rulesObj, "")
	}

	// Collect and compile Yara rules.
	rules, _ := compiler.GetRules()

	// Save the compiled rules to a file.
	rules.Save("rules")

	fmt.Println("[rules-compiler] Done!")
}
