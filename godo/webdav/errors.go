/*
 * godoos - A lightweight cloud desktop
 * Copyright (C) 2024 godoos.com
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
	"errors"
	"fmt"
	"os"
)

// ErrAuthChanged must be returned from the Verify method as an error
// to trigger a re-authentication / negotiation with a new authenticator.
var ErrAuthChanged = errors.New("authentication failed, change algorithm")

// ErrTooManyRedirects will be used as return error if a request exceeds 10 redirects.
var ErrTooManyRedirects = errors.New("stopped after 10 redirects")

// StatusError implements error and wraps
// an erroneous status code.
type StatusError struct {
	Status int
}

func (se StatusError) Error() string {
	return fmt.Sprintf("%d", se.Status)
}

// IsErrCode returns true if the given error
// is an os.PathError wrapping a StatusError
// with the given status code.
func IsErrCode(err error, code int) bool {
	if pe, ok := err.(*os.PathError); ok {
		se, ok := pe.Err.(StatusError)
		return ok && se.Status == code
	}
	return false
}

// IsErrNotFound is shorthand for IsErrCode
// for status 404.
func IsErrNotFound(err error) bool {
	return IsErrCode(err, 404)
}

func NewPathError(op string, path string, statusCode int) error {
	return &os.PathError{
		Op:   op,
		Path: path,
		Err:  StatusError{statusCode},
	}
}

func NewPathErrorErr(op string, path string, err error) error {
	return &os.PathError{
		Op:   op,
		Path: path,
		Err:  err,
	}
}
