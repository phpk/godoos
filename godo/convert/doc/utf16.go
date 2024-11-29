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
package doc

import (
	"encoding/binary"
)

type utf16Buffer struct {
	haveReadLowerByte bool
	char              [2]byte
	data              []uint16
}

func (buf *utf16Buffer) Write(p []byte) (n int, err error) {
	for i := range p {
		buf.WriteByte(p[i])
	}
	return len(p), nil
}

func (buf *utf16Buffer) WriteByte(b byte) error {
	if buf.haveReadLowerByte {
		buf.char[1] = b
		buf.data = append(buf.data, binary.LittleEndian.Uint16(buf.char[:]))
	} else {
		buf.char[0] = b
	}
	buf.haveReadLowerByte = !buf.haveReadLowerByte
	return nil
}

func (buf *utf16Buffer) Chars() []uint16 {
	if buf.haveReadLowerByte {
		return append(buf.data, uint16(buf.char[0]))
	}
	return buf.data
}
