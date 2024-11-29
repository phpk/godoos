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
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unicode/utf16"

	"github.com/mattetti/filebuffer"
	"github.com/richardlehane/mscfb"
)

var (
	errTable           = errors.New("cannot find table stream")
	errDocEmpty        = errors.New("WordDocument not found")
	errDocShort        = errors.New("wordDoc block too short")
	errInvalidArgument = errors.New("invalid table and/or fib")
)

type allReader interface {
	io.Closer
	io.ReaderAt
	io.ReadSeeker
}

func wrapError(e error) error {
	return errors.New("Error processing file: " + e.Error())
}

// ParseDoc converts a standard io.Reader from a Microsoft Word
// .doc binary file and returns a reader (actually a bytes.Buffer)
// which will output the plain text found in the .doc file
func ParseDoc(r io.Reader) (io.Reader, error) {
	ra, ok := r.(io.ReaderAt)
	if !ok {
		ra, _, err := toMemoryBuffer(r)
		if err != nil {
			return nil, wrapError(err)
		}
		defer ra.Close()
	}

	d, err := mscfb.New(ra)
	if err != nil {
		return nil, wrapError(err)
	}

	wordDoc, table0, table1 := getWordDocAndTables(d)
	fib, err := getFib(wordDoc)
	if err != nil {
		return nil, wrapError(err)
	}

	table := getActiveTable(table0, table1, fib)
	if table == nil {
		return nil, wrapError(errTable)
	}

	clx, err := getClx(table, fib)
	if err != nil {
		return nil, wrapError(err)
	}

	return getText(wordDoc, clx)
}

func toMemoryBuffer(r io.Reader) (allReader, int64, error) {
	var b bytes.Buffer
	size, err := b.ReadFrom(r)
	if err != nil {
		return nil, 0, err
	}
	fb := filebuffer.New(b.Bytes())
	return fb, size, nil
}

func getText(wordDoc *mscfb.File, clx *clx) (io.Reader, error) {
	//var buf bytes.Buffer
	var buf utf16Buffer
	for i := 0; i < len(clx.pcdt.PlcPcd.aPcd); i++ {
		pcd := clx.pcdt.PlcPcd.aPcd[i]
		cp := clx.pcdt.PlcPcd.aCP[i]
		cpNext := clx.pcdt.PlcPcd.aCP[i+1]

		//var start, end, size int
		var start, end int
		if pcd.fc.fCompressed {
			//size = 1
			start = pcd.fc.fc / 2
			end = start + cpNext - cp
		} else {
			//size = 2
			start = pcd.fc.fc
			end = start + 2*(cpNext-cp)
		}

		b := make([]byte, end-start)
		//_, err := wordDoc.ReadAt(b, int64(start/size)) // read all the characters
		_, err := wordDoc.ReadAt(b, int64(start))
		if err != nil {
			return nil, err
		}
		translateText(b, &buf, pcd.fc.fCompressed)
	}
	//return &buf, nil
	runes := utf16.Decode(buf.Chars())

	var out bytes.Buffer
	out.Grow(len(runes))
	for _, r := range runes {
		if r == 7 { // table column separator
			r = ' '
		} else if r < 32 && r != 9 && r != 10 && r != 13 { // skip non-printable ASCII characters
			continue
		}
		out.WriteRune(r)
	}

	return &out, nil
}

func translateText(b []byte, buf *utf16Buffer, fCompressed bool) {
	fieldLevel := 0
	var isFieldChar bool
	for cIndex := range b {
		// Handle special field characters (section 2.8.25)
		if b[cIndex] == 0x13 {
			isFieldChar = true
			fieldLevel++
			continue
		} else if b[cIndex] == 0x14 {
			isFieldChar = false
			continue
		} else if b[cIndex] == 0x15 {
			isFieldChar = false
			continue
		} else if isFieldChar {
			continue
		}

		// if b[cIndex] == 7 { // table column separator
		// 	buf.WriteByte(' ')
		// 	continue
		// } else if b[cIndex] < 32 && b[cIndex] != 9 && b[cIndex] != 10 && b[cIndex] != 13 { // skip non-printable ASCII characters
		// 	//buf.Write([]byte(fmt.Sprintf("|%#x|", b[cIndex])))
		// 	continue
		// }

		if fCompressed { // compressed, so replace compressed characters
			buf.Write(replaceCompressed(b[cIndex]))
		} else {
			//buf.Write(b)
			buf.WriteByte(b[cIndex])
		}
	}
}

func replaceCompressed(char byte) []byte {
	var v uint16
	switch char {
	case 0x82:
		v = 0x201A
	case 0x83:
		v = 0x0192
	case 0x84:
		v = 0x201E
	case 0x85:
		v = 0x2026
	case 0x86:
		v = 0x2020
	case 0x87:
		v = 0x2021
	case 0x88:
		v = 0x02C6
	case 0x89:
		v = 0x2030
	case 0x8A:
		v = 0x0160
	case 0x8B:
		v = 0x2039
	case 0x8C:
		v = 0x0152
	case 0x91:
		v = 0x2018
	case 0x92:
		v = 0x2019
	case 0x93:
		v = 0x201C
	case 0x94:
		v = 0x201D
	case 0x95:
		v = 0x2022
	case 0x96:
		v = 0x2013
	case 0x97:
		v = 0x2014
	case 0x98:
		v = 0x02DC
	case 0x99:
		v = 0x2122
	case 0x9A:
		v = 0x0161
	case 0x9B:
		v = 0x203A
	case 0x9C:
		v = 0x0153
	case 0x9F:
		v = 0x0178
	default:
		//return []byte{char}
		return []byte{char, 0x00}
	}
	out := make([]byte, 2)
	binary.LittleEndian.PutUint16(out, v)
	return out
}

func getWordDocAndTables(r *mscfb.Reader) (*mscfb.File, *mscfb.File, *mscfb.File) {
	var wordDoc, table0, table1 *mscfb.File
	for i := 0; i < len(r.File); i++ {
		stream := r.File[i]

		switch stream.Name {
		case "WordDocument":
			wordDoc = stream
		case "0Table":
			table0 = stream
		case "1Table":
			table1 = stream
		}
	}
	return wordDoc, table0, table1
}

func getActiveTable(table0 *mscfb.File, table1 *mscfb.File, f *fib) *mscfb.File {
	if f.base.fWhichTblStm == 0 {
		return table0
	}
	return table1
}
