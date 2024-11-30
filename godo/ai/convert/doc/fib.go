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
	"errors"

	"github.com/richardlehane/mscfb"
)

var (
	errFibInvalid = errors.New("file information block validation failed")
)

type fib struct {
	base       fibBase
	csw        int
	fibRgW     fibRgW
	cslw       int
	fibRgLw    fibRgLw
	cbRgFcLcb  int
	fibRgFcLcb fibRgFcLcb
}

type fibBase struct {
	fWhichTblStm int
}

type fibRgW struct {
}

type fibRgLw struct {
	ccpText    int
	ccpFtn     int
	ccpHdd     int
	ccpMcr     int
	ccpAtn     int
	ccpEdn     int
	ccpTxbx    int
	ccpHdrTxbx int
	cpLength   int
}

type fibRgFcLcb struct {
	fcPlcfFldMom  int
	lcbPlcfFldMom int
	fcPlcfFldHdr  int
	lcbPlcfFldHdr int
	fcPlcfFldFtn  int
	lcbPlcfFldFtn int
	fcPlcfFldAtn  int
	lcbPlcfFldAtn int
	fcClx         int
	lcbClx        int
}

// parse File Information Block (section 2.5.1)
func getFib(wordDoc *mscfb.File) (*fib, error) {
	if wordDoc == nil {
		return nil, errDocEmpty
	}

	b := make([]byte, 898) // get FIB block up to FibRgFcLcb97
	_, err := wordDoc.ReadAt(b, 0)
	if err != nil {
		return nil, err
	}

	fibBase := getFibBase(b[0:32])

	fibRgW, csw, err := getFibRgW(b, 32)
	if err != nil {
		return nil, err
	}

	fibRgLw, cslw, err := getFibRgLw(b, 34+csw)
	if err != nil {
		return nil, err
	}

	fibRgFcLcb, cbRgFcLcb, err := getFibRgFcLcb(b, 34+csw+2+cslw)

	return &fib{base: *fibBase, csw: csw, cslw: cslw, fibRgW: *fibRgW, fibRgLw: *fibRgLw, fibRgFcLcb: *fibRgFcLcb, cbRgFcLcb: cbRgFcLcb}, err
}

// parse FibBase (section 2.5.2)
func getFibBase(fib []byte) *fibBase {
	byt := fib[11]                    // fWhichTblStm is 2nd highest bit in this byte
	fWhichTblStm := int(byt >> 1 & 1) // set which table (0Table or 1Table) is the table stream
	return &fibBase{fWhichTblStm: fWhichTblStm}
}

func getFibRgW(fib []byte, start int) (*fibRgW, int, error) {
	if start+2 >= len(fib) { // must be big enough for csw
		return &fibRgW{}, 0, errFibInvalid
	}

	csw := int(binary.LittleEndian.Uint16(fib[start:start+2])) * 2 // in bytes
	return &fibRgW{}, csw, nil
}

// parse FibRgLw (section 2.5.4)
func getFibRgLw(fib []byte, start int) (*fibRgLw, int, error) {
	fibRgLwStart := start + 2        // skip cslw
	if fibRgLwStart+88 >= len(fib) { // expect 88 bytes in fibRgLw
		return &fibRgLw{}, 0, errFibInvalid
	}

	cslw := getInt16(fib, start) * 4 // in bytes
	ccpText := getInt(fib, fibRgLwStart+3*4)
	ccpFtn := getInt(fib, fibRgLwStart+4*4)
	ccpHdd := getInt(fib, fibRgLwStart+5*4)
	ccpMcr := getInt(fib, fibRgLwStart+6*4)
	ccpAtn := getInt(fib, fibRgLwStart+7*4)
	ccpEdn := getInt(fib, fibRgLwStart+8*4)
	ccpTxbx := getInt(fib, fibRgLwStart+9*4)
	ccpHdrTxbx := getInt(fib, fibRgLwStart+10*4)

	// calculate cpLength. Used in PlcPcd verification (see section 2.8.35)
	var cpLength int
	if ccpFtn != 0 || ccpHdd != 0 || ccpMcr != 0 || ccpAtn != 0 || ccpEdn != 0 || ccpTxbx != 0 || ccpHdrTxbx != 0 {
		cpLength = ccpFtn + ccpHdd + ccpMcr + ccpAtn + ccpEdn + ccpTxbx + ccpHdrTxbx + ccpText + 1
	} else {
		cpLength = ccpText
	}
	return &fibRgLw{ccpText: ccpText, ccpFtn: ccpFtn, ccpHdd: ccpHdd, ccpMcr: ccpMcr, ccpAtn: ccpAtn,
		ccpEdn: ccpEdn, ccpTxbx: ccpTxbx, ccpHdrTxbx: ccpHdrTxbx, cpLength: cpLength}, cslw, nil
}

// parse FibRgFcLcb (section 2.5.5)
func getFibRgFcLcb(fib []byte, start int) (*fibRgFcLcb, int, error) {
	fibRgFcLcbStart := start + 2          // skip cbRgFcLcb
	if fibRgFcLcbStart+186*4 < len(fib) { // expect 186+ values in FibRgFcLcb
		return &fibRgFcLcb{}, 0, errFibInvalid
	}

	cbRgFcLcb := getInt16(fib, start)
	fcPlcfFldMom := getInt(fib, fibRgFcLcbStart+32*4)
	lcbPlcfFldMom := getInt(fib, fibRgFcLcbStart+33*4)
	fcPlcfFldHdr := getInt(fib, fibRgFcLcbStart+34*4)
	lcbPlcfFldHdr := getInt(fib, fibRgFcLcbStart+35*4)
	fcPlcfFldFtn := getInt(fib, fibRgFcLcbStart+36*4)
	lcbPlcfFldFtn := getInt(fib, fibRgFcLcbStart+37*4)
	fcPlcfFldAtn := getInt(fib, fibRgFcLcbStart+38*4)
	lcbPlcfFldAtn := getInt(fib, fibRgFcLcbStart+39*4)
	fcClx := getInt(fib, fibRgFcLcbStart+66*4)
	lcbClx := getInt(fib, fibRgFcLcbStart+67*4)
	return &fibRgFcLcb{fcPlcfFldMom: fcPlcfFldMom, lcbPlcfFldMom: lcbPlcfFldMom, fcPlcfFldHdr: fcPlcfFldHdr, lcbPlcfFldHdr: lcbPlcfFldHdr,
		fcPlcfFldFtn: fcPlcfFldFtn, lcbPlcfFldFtn: lcbPlcfFldFtn, fcPlcfFldAtn: fcPlcfFldAtn, lcbPlcfFldAtn: lcbPlcfFldAtn,
		fcClx: fcClx, lcbClx: lcbClx}, cbRgFcLcb, nil
}

func getInt16(buf []byte, start int) int {
	return int(binary.LittleEndian.Uint16(buf[start : start+2]))
}
func getInt(buf []byte, start int) int {
	return int(binary.LittleEndian.Uint32(buf[start : start+4]))
}
