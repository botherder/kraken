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
	"io"
	"os"
	"path"
)

// Just check if a given string exists in a slice.
func stringInSlice(item string, list []string) bool {
	for _, current := range list {
		if current == item {
			return true
		}
	}

	return false
}

// Get current working directory.
func getCwd() string {
	exe, err := os.Executable()
	if err != nil {
		return ""
	}

	return path.Dir(exe)
}

// I can't believe that I need to make a function to copy a file from
// one location to another.
func copyFile(src, dst string) (err error) {
	srcHandle, err := os.Open(src)
	if err != nil {
		return
	}
	defer srcHandle.Close()

	dstHandle, err := os.Create(dst)
	if err != nil {
		return
	}
	defer dstHandle.Close()

	if _, err = io.Copy(dstHandle, srcHandle); err != nil {
		return
	}

	err = dstHandle.Sync()
	return
}
