// MIT License
//
// Copyright (c) 2024 godoos.com
// Email: xpbb@qq.com
// GitHub: github.com/phpk/godoos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
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
