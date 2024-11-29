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
package convert

import (
	"io"
	"regexp"
	"strings"
)

func ConvertMd(r io.Reader) (string, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	re := regexp.MustCompile(`<[^>]*>`)
	content := re.ReplaceAllString(string(b), "")
	reMarkdown := regexp.MustCompile(`(\*{1,4}|_{1,4}|\#{1,6})`)
	content = reMarkdown.ReplaceAllString(content, "")
	// 移除换行符
	content = strings.ReplaceAll(content, "\r", "")
	content = strings.ReplaceAll(content, "\n", "")

	// 移除多余的空格
	content = strings.TrimSpace(content)
	return content, nil
}
