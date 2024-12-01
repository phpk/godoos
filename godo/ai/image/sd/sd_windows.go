// Copyright (c) seasonjs. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

//go:build windows

package sd

import (
	"golang.org/x/sys/windows"
)

func openLibrary(name string) (uintptr, error) {
	handle, err := windows.LoadLibrary(name)
	return uintptr(handle), err
}

func closeLibrary(handle uintptr) error {
	return windows.FreeLibrary(windows.Handle(handle))
}
