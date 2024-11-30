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
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
)

// ConvertXML converts an XML file to text.
func ConvertXML(r io.Reader) (string, error) {
	cleanXML, err := Tidy(r)
	if err != nil {
		return "", fmt.Errorf("tidy error: %v", err)
	}
	result, err := XMLToText(bytes.NewReader(cleanXML), []string{}, []string{}, true)
	if err != nil {
		return "", fmt.Errorf("error from XMLToText: %v", err)
	}
	// 直接将原始 io.Reader 传递给 XMLToText
	// result, err := XMLToText(r, []string{}, []string{}, true)
	// if err != nil {
	// 	return "", nil, fmt.Errorf("error from XMLToText: %w", err)
	// }
	return result, nil
}

// XMLToText converts XML to plain text given how to treat elements.
func XMLToText(r io.Reader, breaks []string, skip []string, strict bool) (string, error) {
	var result string

	dec := xml.NewDecoder(io.LimitReader(r, maxBytes))
	dec.Strict = strict
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}

		switch v := t.(type) {
		case xml.CharData:
			result += string(v)
		case xml.StartElement:
			for _, breakElement := range breaks {
				if v.Name.Local == breakElement {
					result += "\n"
				}
			}
			for _, skipElement := range skip {
				if v.Name.Local == skipElement {
					depth := 1
					for {
						t, err := dec.Token()
						if err != nil {
							// An io.EOF here is actually an error.
							return "", err
						}

						switch t.(type) {
						case xml.StartElement:
							depth++
						case xml.EndElement:
							depth--
						}

						if depth == 0 {
							break
						}
					}
				}
			}
		}
	}
	return result, nil
}

// XMLToMap converts XML to a nested string map.
func XMLToMap(r io.Reader) (map[string]string, error) {
	m := make(map[string]string)
	dec := xml.NewDecoder(io.LimitReader(r, maxBytes))
	var tagName string
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		switch v := t.(type) {
		case xml.StartElement:
			tagName = string(v.Name.Local)
		case xml.CharData:
			m[tagName] = string(v)
		}
	}
	return m, nil
}
