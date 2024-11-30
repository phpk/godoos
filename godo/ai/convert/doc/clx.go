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
	errInvalidPrc  = errors.New("Invalid Prc structure")
	errInvalidClx  = errors.New("expected last aCP value to equal fib.cpLength (2.8.35)")
	errInvalidPcdt = errors.New("expected clxt to be equal 0x02")
)

type clx struct {
	pcdt pcdt
}

type pcdt struct {
	lcb    int
	PlcPcd plcPcd
}

type plcPcd struct {
	aCP  []int
	aPcd []pcd
}

type pcd struct {
	fc fcCompressed
}

type fcCompressed struct {
	fc          int
	fCompressed bool
}

// read Clx (section 2.9.38)
func getClx(table *mscfb.File, fib *fib) (*clx, error) {
	if table == nil || fib == nil {
		return nil, errInvalidArgument
	}
	b, err := readClx(table, fib)
	if err != nil {
		return nil, err
	}

	pcdtOffset, err := getPrcArrayEnd(b)
	if err != nil {
		return nil, err
	}

	pcdt, err := getPcdt(b, pcdtOffset)
	if err != nil {
		return nil, err
	}

	if pcdt.PlcPcd.aCP[len(pcdt.PlcPcd.aCP)-1] != fib.fibRgLw.cpLength {
		return nil, errInvalidClx
	}

	return &clx{pcdt: *pcdt}, nil
}

func readClx(table *mscfb.File, fib *fib) ([]byte, error) {
	b := make([]byte, fib.fibRgFcLcb.lcbClx)
	_, err := table.ReadAt(b, int64(fib.fibRgFcLcb.fcClx))
	if err != nil {
		return nil, err
	}
	return b, nil
}

// read Pcdt from Clx (section 2.9.178)
func getPcdt(clx []byte, pcdtOffset int) (*pcdt, error) {
	const pcdSize = 8
	if clx[pcdtOffset] != 0x02 { // clxt must be 0x02 or invalid
		return nil, errInvalidPcdt
	}
	lcb := int(binary.LittleEndian.Uint32(clx[pcdtOffset+1 : pcdtOffset+5])) // skip clxt, get lcb
	plcPcdOffset := pcdtOffset + 5                                           // skip clxt and lcb
	numPcds := (lcb - 4) / (4 + pcdSize)                                     // see 2.2.2 in the spec for equation
	numCps := numPcds + 1                                                    // always 1 more cp than pcds

	cps := make([]int, numCps)
	for i := 0; i < numCps; i++ {
		cpOffset := plcPcdOffset + i*4
		cps[i] = int(binary.LittleEndian.Uint32(clx[cpOffset : cpOffset+4]))
	}

	pcdStart := plcPcdOffset + 4*numCps
	pcds := make([]pcd, numPcds)
	for i := 0; i < numPcds; i++ {
		pcdOffset := pcdStart + i*pcdSize
		pcds[i] = *parsePcd(clx[pcdOffset : pcdOffset+pcdSize])
	}
	return &pcdt{lcb: lcb, PlcPcd: plcPcd{aCP: cps, aPcd: pcds}}, nil
}

// find end of RgPrc array (section 2.9.38)
func getPrcArrayEnd(clx []byte) (int, error) {
	prcOffset := 0
	count := 0
	for {
		clxt := clx[prcOffset]
		if clxt != 0x01 { // this is not a Prc, so exit
			return prcOffset, nil
		}
		prcDataCbGrpprl := binary.LittleEndian.Uint16(clx[prcOffset+1 : prcOffset+3]) // skip the clxt and read 2 bytes
		prcOffset += 1 + 2 + int(prcDataCbGrpprl)                                     // skip clxt, cbGrpprl, and GrpPrl

		if count > 10000 || prcDataCbGrpprl <= 0 || prcOffset+3 > len(clx) { // ensure no infinite loop
			return 0, errInvalidPrc
		}
		count++
	}
}

// parse Pcd (section 2.9.177)
func parsePcd(pcdData []byte) *pcd {
	return &pcd{fc: *parseFcCompressed(pcdData[2:6])}
}

// parse FcCompressed (section 2.9.73)
func parseFcCompressed(fcData []byte) *fcCompressed {
	fCompressed := fcData[3]&64 == 64        // check fcompressed value (second bit from lestmost of the last byte in fcdata)
	fcData[3] = fcData[3] & 63               // clear the fcompressed value from data
	fc := binary.LittleEndian.Uint32(fcData) // word doc generally uses little endian order (1.3.7)
	return &fcCompressed{fc: int(fc), fCompressed: fCompressed}
}
