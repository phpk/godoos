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
	"fmt"
	"io"
	"net/http"

	"jaytaylor.com/html2text"
)

func resErr(err error) Res {
	return Res{
		Status: 201,
		Data:   fmt.Sprintf("error opening file: %v", err),
	}
}
func ConvertHttp(url string) Res {
	resp, err := http.Get(url)
	if err != nil {
		return resErr(err)
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(resp.Body)
	if errRead != nil {
		return resErr(errRead)
	}
	text, err := html2text.FromString(string(body), html2text.Options{PrettyTables: false})
	if err != nil {
		return resErr(err)
	}
	return Res{
		Status: 0,
		Data:   text,
	}
}
