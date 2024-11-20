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
	"bufio"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"
)

func parseLine(s string) (login, pass string) {
	fields := strings.Fields(s)
	for i, f := range fields {
		if f == "login" {
			login = fields[i+1]
		}
		if f == "password" {
			pass = fields[i+1]
		}
	}
	return login, pass
}

// ReadConfig reads login and password configuration from ~/.netrc
// machine foo.com login username password 123456
func ReadConfig(uri, netrc string) (string, string) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", ""
	}

	file, err := os.Open(netrc)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	re := fmt.Sprintf(`^.*machine %s.*$`, u.Host)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()

		matched, err := regexp.MatchString(re, s)
		if err != nil {
			return "", ""
		}
		if matched {
			return parseLine(s)
		}
	}

	return "", ""
}
