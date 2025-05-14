package office

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"unicode/utf16"
	"unicode/utf8"

	"github.com/mattetti/filebuffer"
	"github.com/richardlehane/mscfb"
)

// ---- file doc.go ----

var (
	errTable    = errors.New("cannot find table stream")
	errDocEmpty = errors.New("WordDocument not found")
	// errDocShort        = errors.New("wordDoc block too short")
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

// DOC2Text converts a standard io.Reader from a Microsoft Word .doc binary file and returns a reader (actually a bytes.Buffer) which will output the plain text found in the .doc file
func DOC2Text(r io.Reader) (io.Reader, error) {
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
	var buf bytes.Buffer
	for i := 0; i < len(clx.pcdt.PlcPcd.aPcd); i++ {
		pcd := clx.pcdt.PlcPcd.aPcd[i]
		cp := clx.pcdt.PlcPcd.aCP[i]
		cpNext := clx.pcdt.PlcPcd.aCP[i+1]

		var start, end int
		// https://msdn.microsoft.com/ko-kr/library/office/gg615596(v=office.14).aspx
		// Read the value of the Pcd.Fc.fCompressed field at bit 46 of the current Pcd structure. If 0, the Pcd structure refers to a 16-bit Unicode character. If 1, it refers to an 8-bit ANSI character.
		if pcd.fc.fCompressed {
			start = pcd.fc.fc / 2
			end = start + cpNext - cp
		} else {
			// -> 16-bit Unicode characters
			start = pcd.fc.fc
			end = start + 2*(cpNext-cp)
		}

		b := make([]byte, end-start)
		_, err := wordDoc.ReadAt(b, int64(start)) // read all the characters
		if err != nil {
			return nil, err
		}
		translateText(b, &buf, pcd.fc.fCompressed)
	}
	return &buf, nil
}

// translateText translates the buffer into text. fCompressed = 0 for 16-bit Unicode, 1 = 8-bit ANSI characters.
func translateText(b []byte, buf *bytes.Buffer, fCompressed bool) {
	u16s := make([]uint16, 1)
	b8buf := make([]byte, 4)

	fieldLevel := 0
	var isFieldChar bool
	for cIndex := range b {
		// Convert to rune
		var char rune
		if fCompressed {
			// ANSI, 1 byte
			char = rune(b[cIndex])
		} else {
			// 16-bit Unicode: skip every second byte
			if cIndex%2 != 0 {
				continue
			} else if (cIndex + 1) >= len(b) { // make sure there are at least 2 bytes for Unicode decoding
				continue
			}

			// convert from UTF16 to UTF8
			u16s[0] = uint16(b[cIndex]) + (uint16(b[cIndex+1]) << 8)
			r := utf16.Decode(u16s)
			if len(r) != 1 {
				//fmt.Printf("Invalid rune %v\n", r)
				continue
			}
			char = r[0]
		}

		// Handle special field characters (section 2.8.25)
		if char == 0x13 {
			isFieldChar = true
			fieldLevel++
			continue
		} else if char == 0x14 {
			isFieldChar = false
			continue
		} else if char == 0x15 {
			isFieldChar = false
			continue
		} else if isFieldChar {
			continue
		}

		if char == 7 { // table column separator
			buf.WriteByte(' ')
			continue
		} else if char < 32 && char != 9 && char != 10 && char != 13 { // skip non-printable ASCII characters
			//buf.Write([]byte(fmt.Sprintf("|%#x|", char)))
			continue
		}

		if fCompressed { // compressed, so replace compressed characters
			buf.Write(replaceCompressed(byte(char)))
		} else {
			// encode the rune to UTF-8
			n := utf8.EncodeRune(b8buf, char)
			buf.Write(b8buf[:n])
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
		return []byte{char}
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

// ---- file fib.go ----

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

// ---- file clx.go ----

var (
	errInvalidPrc  = errors.New("invalid Prc structure")
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
	if pcdtOffset < 0 || pcdtOffset+5 >= len(clx) {
		return nil, errInvalidPcdt
	}
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
		if cpOffset < 0 || cpOffset+4 >= len(clx) {
			return nil, errInvalidPcdt
		}
		cps[i] = int(binary.LittleEndian.Uint32(clx[cpOffset : cpOffset+4]))
	}

	pcdStart := plcPcdOffset + 4*numCps
	pcds := make([]pcd, numPcds)
	for i := 0; i < numPcds; i++ {
		pcdOffset := pcdStart + i*pcdSize
		if pcdOffset < 0 || pcdOffset+pcdSize > len(clx) {
			return nil, errInvalidPcdt
		}
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

// IsFileDOC checks if the data indicates a DOC file
// DOC has multiple signature according to https://filesignatures.net/index.php?search=doc&mode=EXT, D0 CF 11 E0 A1 B1 1A E1
func IsFileDOC(data []byte) bool {
	return bytes.HasPrefix(data, []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1})
}
