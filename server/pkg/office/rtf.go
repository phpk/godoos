// Package rtftxt extracts text from .rtf documents
package office

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EndFirstCorp/peekingReader"
)

// ToStr converts a .rtf document file to string
func rtf2txt(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}

	return BytesToStr(content)
}

// BytesToStr converts a []byte representation of a .rtf document file to string
func BytesToStr(data []byte) (string, error) {
	reader := bytes.NewReader(data)
	r, err := Text(reader)
	if err != nil {
		return "", err
	}
	s := r.String()
	return s, nil
}

type stack struct {
	top  *element
	size int
}

type element struct {
	value string
	next  *element
}

func (s *stack) Len() int {
	return s.size
}

func (s *stack) Push(value string) {
	s.top = &element{value, s.top}
	s.size++
}

func (s *stack) Peek() string {
	if s.size == 0 {
		return ""
	}
	return s.top.value
}

func (s *stack) Pop() string {
	if s.size > 0 {
		var v string
		v, s.top = s.top.value, s.top.next
		s.size--
		return v
	}
	return ""
}

// Text is used to convert an io.Reader containing RTF data into
// plain text
func Text(r io.Reader) (*bytes.Buffer, error) {
	pr := peekingReader.NewBufReader(r)

	var text bytes.Buffer
	var symbolStack stack
	for b, err := pr.ReadByte(); err == nil; b, err = pr.ReadByte() {
		switch b {
		case '\\':
			err := readControl(pr, &symbolStack, &text)
			if err != nil {
				return nil, err
			}
		case '{', '}':
		case '\n', '\r': // noop
		default:
			text.WriteByte(b)
		}
	}
	return &text, nil
}

func readControl(r peekingReader.Reader, s *stack, text *bytes.Buffer) error {
	control, num, err := tokenizeControl(r)
	if err != nil {
		return err
	}
	if control == "*" { // this is an extended control sequence
		err = readUntilClosingBrace(r)
		if err != nil {
			return err
		}
		if last := s.Peek(); last != "" {
			val, err := getParams(r) // last control was interrupted, so finish handling Params
			handleParams(control, val, text)
			return err
		}
		return nil
	}
	if isUnicode, u := getUnicode(control); isUnicode {
		text.WriteString(u)
		return nil
	}
	if control == "" {
		p, err := r.Peek(1)
		if err != nil {
			return err
		}
		if p[0] == '\\' || p[0] == '{' || p[0] == '}' { // this is an escaped character
			text.WriteByte(p[0])
			r.ReadByte()
			return nil
		}
		text.WriteByte('\n')
		return nil
	}
	if control == "binN" {
		return handleBinary(r, control, num)
	}

	if symbol, found := convertSymbol(control); found {
		text.WriteString(symbol)
	}

	val, err := getParams(r)
	if err != nil {
		return err
	}
	handleParams(control, val, text)
	s.Push(control)
	return nil
}

func tokenizeControl(r peekingReader.Reader) (string, int, error) {
	var buf bytes.Buffer
	isHex := false
	numStart := -1
	for {
		p, err := r.Peek(1)
		if err != nil {
			return "", -1, err
		}
		b := p[0]
		switch {
		case b == '*' && buf.Len() == 0:
			r.ReadByte() // consume valid digit
			return "*", -1, nil
		case b == '\'' && buf.Len() == 0:
			isHex = true
			buf.WriteByte(b)
			r.ReadByte() // consume valid character
			// read 2 bytes for hex
			for i := 0; i < 2; i++ {
				b, err = r.ReadByte() // consume valid digit
				if err != nil {
					return "", -1, err
				}
				buf.WriteByte(b)
			}
			return buf.String(), -1, nil
		case b >= '0' && b <= '9' || b == '-':
			if numStart == -1 {
				numStart = buf.Len()
			} else if numStart == 0 {
				return "", -1, fmt.Errorf("unexpected control sequence. Cannot begin with digit")
			}
			buf.WriteByte(b)
			r.ReadByte() // consume valid digit
		case b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z':
			if numStart > 0 { // we've already seen alpha character(s) plus digit(s)
				c, num := canonicalize(buf.String(), numStart)
				return c, num, nil
			}
			buf.WriteByte(b)
			r.ReadByte()
		default:
			if isHex {
				return buf.String(), -1, nil
			}
			c, num := canonicalize(buf.String(), numStart)
			return c, num, nil
		}
	}
}

func canonicalize(control string, numStart int) (string, int) {
	if numStart == -1 || numStart >= len(control) {
		return control, -1
	}
	num, err := strconv.Atoi(control[numStart:])
	if err != nil {
		return control, -1
	}
	return control[:numStart] + "N", num
}

func getUnicode(control string) (bool, string) {
	if len(control) < 2 || control[0] != '\'' {
		return false, ""
	}

	var buf bytes.Buffer
	for i := 1; i < len(control); i++ {
		b := control[i]
		if b >= '0' && b <= '9' || b >= 'a' && b <= 'f' || b >= 'A' && b <= 'F' {
			buf.WriteByte(b)
		} else {
			break
		}
	}
	after := control[buf.Len()+1:]
	num, _ := strconv.ParseInt(buf.String(), 16, 16)
	return true, fmt.Sprintf("%c%s", num, after)
}

func getParams(r peekingReader.Reader) (string, error) {
	data, err := peekingReader.ReadUntilAny(r, []byte{'\\', '{', '}', '\n', '\r', ';'})
	if err != nil {
		return "", err
	}
	p, err := r.Peek(1)
	if err != nil {
		return "", err
	}
	if p[0] == ';' { // skip next if it is a semicolon
		r.ReadByte()
	}

	return string(data), nil
}

func handleBinary(r peekingReader.Reader, control string, size int) error {
	if control != "binN" { // wrong control type
		return nil
	}

	_, err := r.ReadBytes(size)
	if err != nil {
		return err
	}
	return nil
}

func readUntilClosingBrace(r peekingReader.Reader) error {
	count := 1
	var b byte
	var err error
	for b, err = r.ReadByte(); err == nil; b, err = r.ReadByte() {
		switch b {
		case '{':
			count++
		case '}':
			count--
		}
		if count == 0 {
			return nil
		}
	}
	return err
}

func handleParams(control, param string, text io.StringWriter) {
	param = strings.TrimPrefix(param, " ")
	if param == "" {
		return
	}
	switch control {
	case "fldrslt":
		text.WriteString(param)
	case "acccircle", "acccomma", "accdot", "accnone", "accunderdot",
		"animtextN", "b", "caps", "cbN", "cchsN ", "cfN", "charscalexN",
		"csN", "dnN", "embo", "expndN", "expndtwN ", "fittextN", "fN",
		"fsN", "i", "kerningN ", "langfeN", "langfenpN", "langN", "langnpN",
		"ltrch", "noproof", "nosupersub ", "outl", "plain", "rtlch", "scaps",
		"shad", "strike", "sub ", "super ", "ul", "ulcN", "uld", "uldash",
		"uldashd", "uldashdd", "uldb", "ulhwave", "ulldash", "ulnone", "ulth",
		"ulthd", "ulthdash", "ulthdashd", "ulthdashdd", "ulthldash", "ululdbwave", "ulw", "ulwave", "upN", "v", "webhidden":
		text.WriteString(param)

	// Paragraph Formatting Properties
	case "aspalpha", "aspnum", "collapsed", "contextualspace",
		"cufiN", "culiN", "curiN", "faauto", "facenter",
		"fafixed", "fahang", "faroman", "favar", "fiN", "hyphpar ",
		"indmirror", "intbl", "itapN", "keep", "keepn", "levelN", "liN",
		"linN", "lisaN", "lisbN", "ltrpar", "nocwrap", "noline", "nooverflow",
		"nosnaplinegrid", "nowidctlpar ", "nowwrap", "outlinelevelN ", "pagebb",
		"pard", "prauthN", "prdateN", "qc", "qd", "qj", "qkN", "ql", "qr", "qt",
		"riN", "rinN", "rtlpar", "saautoN", "saN", "sbautoN", "sbN", "sbys",
		"slmultN", "slN", "sN", "spv", "subdocumentN ", "tscbandhorzeven",
		"tscbandhorzodd", "tscbandverteven", "tscbandvertodd", "tscfirstcol",
		"tscfirstrow", "tsclastcol", "tsclastrow", "tscnecell", "tscnwcell",
		"tscsecell", "tscswcell", "txbxtwalways", "txbxtwfirst", "txbxtwfirstlast",
		"txbxtwlast", "txbxtwno", "widctlpar", "ytsN":
		text.WriteString(param)

	// Section Formatting Properties
	case "adjustright", "binfsxnN", "binsxnN", "colnoN ", "colsN", "colsrN ", "colsxN", "colwN ", "dsN", "endnhere", "footeryN", "guttersxnN", "headeryN", "horzsect", "linebetcol", "linecont", "linemodN", "lineppage", "linerestart", "linestartsN", "linexN", "lndscpsxn", "ltrsect", "margbsxnN", "marglsxnN", "margmirsxn", "margrsxnN", "margtsxnN", "pghsxnN", "pgnbidia", "pgnbidib", "pgnchosung", "pgncnum", "pgncont", "pgndbnum", "pgndbnumd", "pgndbnumk", "pgndbnumt", "pgndec", "pgndecd", "pgnganada", "pgngbnum", "pgngbnumd", "pgngbnumk", "pgngbnuml", "pgnhindia", "pgnhindib", "pgnhindic", "pgnhindid", "pgnhnN ", "pgnhnsc ", "pgnhnsh ", "pgnhnsm ", "pgnhnsn ", "pgnhnsp ", "pgnid", "pgnlcltr", "pgnlcrm", "pgnrestart", "pgnstartsN", "pgnthaia", "pgnthaib", "pgnthaic", "pgnucltr", "pgnucrm", "pgnvieta", "pgnxN", "pgnyN", "pgnzodiac", "pgnzodiacd", "pgnzodiacl", "pgwsxnN", "pnseclvlN", "rtlsect", "saftnnalc", "saftnnar", "saftnnauc", "saftnnchi", "saftnnchosung", "saftnncnum", "saftnndbar", "saftnndbnum", "saftnndbnumd", "saftnndbnumk", "saftnndbnumt", "saftnnganada", "saftnngbnum", "saftnngbnumd", "saftnngbnumk", "saftnngbnuml", "saftnnrlc", "saftnnruc", "saftnnzodiac", "saftnnzodiacd", "saftnnzodiacl", "saftnrestart", "saftnrstcont", "saftnstartN", "sbkcol", "sbkeven", "sbknone", "sbkodd", "sbkpage", "sectd", "sectdefaultcl", "sectexpandN", "sectlinegridN", "sectspecifycl", "sectspecifygenN", "sectspecifyl", "sectunlocked", "sftnbj", "sftnnalc", "sftnnar", "sftnnauc", "sftnnchi", "sftnnchosung", "sftnncnum", "sftnndbar", "sftnndbnum", "sftnndbnumd", "sftnndbnumk", "sftnndbnumt", "sftnnganada", "sftnngbnum", "sftnngbnumd", "sftnngbnumk", "sftnngbnuml", "sftnnrlc", "sftnnruc", "sftnnzodiac", "sftnnzodiacd", "sftnnzodiacl", "sftnrestart", "sftnrstcont", "sftnrstpg", "sftnstartN", "sftntj", "srauthN", "srdateN", "titlepg", "vertal", "vertalb", "vertalc", "vertalj", "vertalt", "vertsect":
		text.WriteString(param)

	// Section Text
	case "stextflowN":
		text.WriteString(param)

	// Special Characters
	case "-", ":", "_", "{", "|", "}", "~", "bullet", "chatn", "chdate", "chdpa", "chdpl", "chftn", "chftnsep", "chftnsepc", "chpgn", "chtime", "column", "emdash", "emspace ",
		"endash", "enspace ", "lbrN", "ldblquote", "line", "lquote", "ltrmark", "page", "par", "qmspace", "rdblquote", "row", "rquote", "rtlmark", "sect", "sectnum", "softcol ", "softlheightN ", "softline ", "softpage ", "tab", "zwbo", "zwj", "zwnbo", "zwnj":
		text.WriteString(param)

	// Table Definitions
	case "cell", "cellxN", "clbgbdiag", "clbgcross", "clbgdcross", "clbgdkbdiag", "clbgdkcross", "clbgdkdcross", "clbgdkfdiag", "clbgdkhor", "clbgdkvert", "clbgfdiag", "clbghoriz", "clbgvert", "clbrdrb", "clbrdrl", "clbrdrr", "clbrdrt", "clcbpatN", "clcbpatrawN", "clcfpatN", "clcfpatrawN", "cldel2007", "cldelauthN", "cldeldttmN", "cldgll", "cldglu", "clFitText", "clftsWidthN", "clhidemark", "clins", "clinsauthN", "clinsdttmN", "clmgf", "clmrg", "clmrgd", "clmrgdauthN", "clmrgddttmN", "clmrgdr", "clNoWrap", "clpadbN", "clpadfbN", "clpadflN", "clpadfrN", "clpadftN", "clpadlN", "clpadrN", "clpadtN", "clshdngN", "clshdngrawN", "clshdrawnil", "clspbN", "clspfbN", "clspflN", "clspfrN", "clspftN", "clsplit", "clsplitr", "clsplN", "clsprN", "clsptN", "cltxbtlr", "cltxlrtb", "cltxlrtbv", "cltxtbrl", "cltxtbrlv", "clvertalb", "clvertalc", "clvertalt", "clvmgf", "clvmrg", "clwWidthN", "irowbandN", "irowN", "lastrow", "ltrrow", "nestcell", "nestrow", "nesttableprops", "nonesttables", "rawclbgbdiag", "rawclbgcross", "rawclbgdcross", "rawclbgdkbdiag", "rawclbgdkcross", "rawclbgdkdcross", "rawclbgdkfdiag", "rawclbgdkhor", "rawclbgdkvert", "rawclbgfdiag", "rawclbghoriz", "rawclbgvert", "rtlrow", "tabsnoovrlp", "taprtl", "tblindN", "tblindtypeN", "tbllkbestfit", "tbllkborder", "tbllkcolor", "tbllkfont", "tbllkhdrcols", "tbllkhdrrows", "tbllklastcol", "tbllklastrow", "tbllknocolband", "tbllknorowband", "tbllkshading", "tcelld", "tdfrmtxtBottomN", "tdfrmtxtLeftN", "tdfrmtxtRightN", "tdfrmtxtTopN", "tphcol", "tphmrg", "tphpg", "tposnegxN", "tposnegyN", "tposxc", "tposxi", "tposxl", "tposxN", "tposxo", "tposxr", "tposyb", "tposyc", "tposyil", "tposyin", "tposyN", "tposyout", "tposyt", "tpvmrg", "tpvpara", "tpvpg", "trauthN", "trautofitN", "trbgbdiag", "trbgcross", "trbgdcross", "trbgdkbdiag", "trbgdkcross", "trbgdkdcross", "trbgdkfdiag", "trbgdkhor", "trbgdkvert", "trbgfdiag", "trbghoriz", "trbgvert", "trbrdrb ", "trbrdrh ", "trbrdrl ", "trbrdrr ", "trbrdrt ", "trbrdrv ", "trcbpatN", "trcfpatN", "trdateN", "trftsWidthAN", "trftsWidthBN", "trftsWidthN", "trgaphN", "trhdr ", "trkeep ", "trkeepfollow", "trleftN", "trowd", "trpaddbN", "trpaddfbN", "trpaddflN", "trpaddfrN", "trpaddftN", "trpaddlN", "trpaddrN", "trpaddtN", "trpadobN", "trpadofbN", "trpadoflN", "trpadofrN", "trpadoftN", "trpadolN", "trpadorN", "trpadotN", "trpatN", "trqc", "trql", "trqr", "trrhN", "trshdngN", "trspdbN", "trspdfbN", "trspdflN", "trspdfrN", "trspdftN", "trspdlN", "trspdrN", "trspdtN", "trspobN", "trspofbN", "trspoflN", "trspofrN", "trspoftN", "trspolN", "trsporN", "trspotN", "trwWidthAN", "trwWidthBN", "trwWidthN":
		text.WriteString(param)

	// Table of Contents Entries
	case "tc", "tcfN", "tclN", "tcn ":
		text.WriteString(param)

	// Tabs
	case "tbN", "tldot", "tleq", "tlhyph", "tlmdot", "tlth", "tlul", "tqc", "tqdec", "tqr", "txN":
		text.WriteString(param)
	}
}

func convertSymbol(symbol string) (string, bool) {
	switch symbol {
	case "bullet":
		return "*", true
	case "chdate", "chdpa", "chdpl":
		return time.Now().Format("2005-01-02"), true
	case "chtime":
		return time.Now().Format("4:56 pm"), true
	case "emdash", "endash":
		return "-", true
	case "lquote", "rquote":
		return "'", true
	case "ldblquote", "rdblquote":
		return "\"", true
	case "line", "lbrN":
		return "\n", true
	case "cell", "column", "emspace", "enspace", "qmspace", "nestcell", "nestrow", "page", "par", "row", "sect", "tab":
		return " ", true
	case "|", "~", "-", "_", ":":
		return symbol, true
	case "chatn", "chftn", "chftnsep", "chftnsepc", "chpgn", "sectnum", "ltrmark", "rtlmark", "zwbo", "zwj", "zwnbo", "zwnj", "softcol",
		"softline", "softpage":
		return "", true
	default:
		return "", false
	}
}
