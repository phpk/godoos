/*
 * GodoOS - A lightweight cloud desktop
 * Copyright (C) 2024 https://godoos.com
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package webdav

import (
	"fmt"
	"os"
	"time"
)

// File is our structure for a given file
type File struct {
	path        string
	name        string
	contentType string
	size        int64
	modified    time.Time
	etag        string
	isdir       bool
}

// Path returns the full path of a file
func (f File) Path() string {
	return f.path
}

// Name returns the name of a file
func (f File) Name() string {
	return f.name
}

// ContentType returns the content type of a file
func (f File) ContentType() string {
	return f.contentType
}

// Size returns the size of a file
func (f File) Size() int64 {
	return f.size
}

// Mode will return the mode of a given file
func (f File) Mode() os.FileMode {
	// TODO check webdav perms
	if f.isdir {
		return 0775 | os.ModeDir
	}

	return 0664
}

// ModTime returns the modified time of a file
func (f File) ModTime() time.Time {
	return f.modified
}

// ETag returns the ETag of a file
func (f File) ETag() string {
	return f.etag
}

// IsDir let us see if a given file is a directory or not
func (f File) IsDir() bool {
	return f.isdir
}

// Sys ????
func (f File) Sys() interface{} {
	return nil
}

// String lets us see file information
func (f File) String() string {
	if f.isdir {
		return fmt.Sprintf("Dir : '%s' - '%s'", f.path, f.name)
	}

	return fmt.Sprintf("File: '%s' SIZE: %d MODIFIED: %s ETAG: %s CTYPE: %s", f.path, f.size, f.modified.String(), f.etag, f.contentType)
}
