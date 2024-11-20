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
	"fmt"
	"net/http"
)

// BasicAuth structure holds our credentials
type BasicAuth struct {
	user string
	pw   string
}

// Authorize the current request
func (b *BasicAuth) Authorize(c *http.Client, rq *http.Request, path string) error {
	rq.SetBasicAuth(b.user, b.pw)
	return nil
}

// Verify verifies if the authentication
func (b *BasicAuth) Verify(c *http.Client, rs *http.Response, path string) (redo bool, err error) {
	if rs.StatusCode == 401 {
		err = NewPathError("Authorize", path, rs.StatusCode)
	}
	return
}

// Close cleans up all resources
func (b *BasicAuth) Close() error {
	return nil
}

// Clone creates a Copy of itself
func (b *BasicAuth) Clone() Authenticator {
	// no copy due to read only access
	return b
}

// String toString
func (b *BasicAuth) String() string {
	return fmt.Sprintf("BasicAuth login: %s", b.user)
}
