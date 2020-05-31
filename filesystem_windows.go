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
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	kernel32             = windows.NewLazySystemDLL("kernel32.dll")
	procGetDriveTypeW    = kernel32.NewProc("GetDriveTypeW")
	procGetLogicalDrives = kernel32.NewProc("GetLogicalDrives")
)

func getDrives(bitMap uint32) (drives []string) {
	driveLetters := []string{
		"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M",
		"N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	}

	for _, letter := range driveLetters {
		if bitMap&1 == 1 {
			drives = append(drives, fmt.Sprintf("%s:\\", letter))
		}
		bitMap >>= 1
	}

	return
}

func getDriveType(driveType uint32) string {
	switch driveType {
	case windows.DRIVE_CDROM:
		return "cd-rom"
	case windows.DRIVE_FIXED:
		return "fixed"
	case windows.DRIVE_NO_ROOT_DIR:
		return "no-root-dir"
	case windows.DRIVE_RAMDISK:
		return "ram-disk"
	case windows.DRIVE_REMOTE:
		return "remote"
	case windows.DRIVE_REMOVABLE:
		return "removable"
	case windows.DRIVE_UNKNOWN:
		return "unknown"
	}
	return "unrecognized"
}

func getFileSystemRoots() []string {
	var drives []string
	var toScan []string

	ret, _, _ := procGetLogicalDrives.Call()
	drives = getDrives(uint32(ret))

	for _, drive := range drives {
		dtp, _, _ := procGetDriveTypeW.Call(uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(drive))))
		driveType := getDriveType(uint32(dtp))

		// TODO: Shall we scan also removables?
		if driveType == "fixed" {
			toScan = append(toScan, drive)
		}
	}

	return toScan
}
