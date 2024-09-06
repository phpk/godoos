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
