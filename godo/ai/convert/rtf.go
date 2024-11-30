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
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/EndFirstCorp/peekingReader"
)

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
func ConvertRTF(r io.Reader) (string, error) {
	pr := peekingReader.NewBufReader(r)

	var text bytes.Buffer
	var symbolStack stack
	for b, err := pr.ReadByte(); err == nil; b, err = pr.ReadByte() {
		switch b {
		case '\\':
			err := ReadRtfControl(pr, &symbolStack, &text)
			if err != nil {
				return "", err
			}
		case '{', '}':
		case '\n', '\r': // noop
		default:
			text.WriteByte(b)
		}
	}
	return string(text.Bytes()), nil
}

func ReadRtfControl(r peekingReader.Reader, s *stack, text *bytes.Buffer) error {
	control, num, err := tokenizeControl(r)
	if err != nil {
		return err
	}
	if control == "*" { // this is an extended control sequence
		err := readUntilClosingBrace(r)
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
			r.ReadByte()
		case b >= '0' && b <= '9' || b == '-':
			if numStart == -1 {
				numStart = buf.Len()
			} else if numStart == 0 {
				return "", -1, errors.New("Unexpected control sequence. Cannot begin with digit")
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

func handleParams(control, param string, text *bytes.Buffer) {
	if strings.HasPrefix(param, " ") {
		param = param[1:]
	}
	if param == "" {
		return
	}
	switch control {
	// Absolution Position Tabs
	// case "pindtabqc", "pindtabql", "pindtabqr", "pmartabqc", "pmartabql", "pmartabqr", "ptabldot", "ptablmdot", "ptablminus", "ptablnone", "ptabluscore":

	// Associated Character Properties
	// case "ab","acaps","acfN","adnN","aexpndN","afN","afsN","ai","alangN","aoutl","ascaps","ashad","astrike","aul","auld","auldb","aulnone","aulw","aupN","dbch","fcsN","hich","loch":

	// Bookmarks
	// case "bkmkcolfN","bkmkcollN","bkmkend","bkmkstart":

	// Bullets and Numbering
	// case "ilvlN","listtext","pn ","pnacross ","pnaiu","pnaiud","pnaiueo","pnaiueod","pnb ","pnbidia","pnbidib","pncaps ","pncard ","pncfN ","pnchosung","pncnum","pndbnum","pndbnumd","pndbnumk","pndbnuml","pndbnumt","pndec ","pndecd","pnfN ","pnfsN ","pnganada","pngbnum","pngbnumd","pngbnumk","pngbnuml","pnhang ","pni ","pnindentN ","pniroha","pnirohad","pnlcltr ","pnlcrm ","pnlvlblt ","pnlvlbody ","pnlvlcont ","pnlvlN ","pnnumonce ","pnord ","pnordt ","pnprev ","pnqc ","pnql ","pnqr ","pnrestart ","pnscaps ","pnspN ","pnstartN ","pnstrike ","pntext ","pntxta ","pntxtb ","pnucltr ","pnucrm ","pnul ","pnuld ","pnuldash","pnuldashd","pnuldashdd","pnuldb ","pnulhair","pnulnone ","pnulth","pnulw ","pnulwave","pnzodiac","pnzodiacd","pnzodiacl":

	// Character Borders and Shading
	// case "chbgbdiag","chbgcross","chbgdcross","chbgdkbdiag","chbgdkcross","chbgdkdcross","chbgdkfdiag","chbgdkhoriz","chbgdkvert","chbgfdiag","chbghoriz","chbgvert","chbrdr","chcbpatN","chcfpatN","chshdngN":

	// Character Revision Mark Properties
	// case "crauthN","crdateN","deleted","mvauthN ","mvdateN ","mvf","mvt","revauthdelN","revauthN ","revdttmdelN","revdttmN ","revised":

	// Character Set
	// case "ansi","ansicpgN","fbidis","mac","pc","pca","impr","striked1":

	// Code Page Support
	// case "cpgN":

	// Color Scheme Mapping
	// case "colorschememapping":

	// Color Table
	// case "blueN","caccentfive","caccentfour","caccentone","caccentsix","caccentthree","caccenttwo","cbackgroundone","cbackgroundtwo","cfollowedhyperlink","chyperlink","cmaindarkone","cmaindarktwo","cmainlightone","cmainlighttwo","colortbl","cshadeN","ctextone","ctexttwo","ctintN","greenN","redN":

	// Comments (Annotations)
	// case "annotation","atnauthor","atndate ","atnicn","atnid","atnparent","atnref ","atntime","atrfend ","atrfstart ":

	// Control Words Introduced by Other Microsoft Products
	// case "disabled","htmlbase ","htmlrtf","htmltag","mhtmltag","protect","pwdN","urtfN":

	// Custom XML Data Properties
	// case "datastore":

	// Custom XML Tags
	// case "xmlattr","xmlattrname","xmlattrnsN","xmlattrvalue","xmlclose","xmlname","xmlnstbl","xmlopen","xmlsdttcell","xmlsdttpara","xmlsdttregular","xmlsdttrow","xmlsdttunknown","xmlnsN":

	// Default Fonts
	// case "adeffN","adeflangN","deffN","deflangfeN","deflangN","stshfbiN","stshfdbchN","stshfhichN","stshflochN":

	// Default Properties
	// case "defchp","defpap":

	// Document Formatting Properties
	// case "aenddoc","aendnotes","afelev","aftnbj","aftncn","aftnnalc","aftnnar","aftnnauc","aftnnchi","aftnnchosung","aftnncnum","aftnndbar","aftnndbnum","aftnndbnumd","aftnndbnumk","aftnndbnumt","aftnnganada","aftnngbnum","aftnngbnumd","aftnngbnumk","aftnngbnuml","aftnnrlc ","aftnnruc ","aftnnzodiac","aftnnzodiacd","aftnnzodiacl","aftnrestart ","aftnrstcont ","aftnsep ","aftnsepc ","aftnstartN","aftntj ","allowfieldendsel","allprot ","alntblind","annotprot ","ApplyBrkRules","asianbrkrule","autofmtoverride","background","bdbfhdr","bdrrlswsix","bookfold","bookfoldrev","bookfoldsheetsN","brdrartN","brkfrm ","cachedcolbal","ctsN","cvmme ","defformat","deftabN","dghoriginN","dghshowN","dghspaceN","dgmargin","dgsnap","dgvoriginN","dgvshowN","dgvspaceN","dntblnsbdb","doctemp","doctypeN","donotembedlingdataN","donotembedsysfontN","donotshowcomments","donotshowinsdel","donotshowmarkup","donotshowprops","enddoc","endnotes","enforceprotN","expshrtn","facingp","fchars","felnbrelev","fetN ","forceupgrade","formdisp ","formprot ","formshade ","fracwidth","fromhtmlN","fromtext","ftnalt ","ftnbj","ftncn","ftnlytwnine","ftnnalc ","ftnnar ","ftnnauc ","ftnnchi ","ftnnchosung","ftnncnum","ftnndbar","ftnndbnum","ftnndbnumd","ftnndbnumk","ftnndbnumt","ftnnganada","ftnngbnum","ftnngbnumd","ftnngbnumk","ftnngbnuml","ftnnrlc ","ftnnruc ","ftnnzodiac","ftnnzodiacd","ftnnzodiacl","ftnrestart","ftnrstcont ","ftnrstpg ","ftnsep","ftnsepc","ftnstartN","ftntj","grfdoceventsN","gutterN","gutterprl","horzdoc","htmautsp","hwelev2007","hyphauto ","hyphcaps ","hyphconsecN ","hyphhotzN","ignoremixedcontentN","ilfomacatclnupN","indrlsweleven","jcompress","jexpand","jsksu","krnprsnet","ksulangN","landscape","lchars","linestartN","linkstyles ","lnbrkrule","lnongrid","ltrdoc","lytcalctblwd","lytexcttp","lytprtmet","lyttblrtgr","makebackup","margbN","marglN","margmirror","margrN","margtN","msmcap","muser","newtblstyruls","nextfile","noafcnsttbl","nobrkwrptbl","nocolbal ","nocompatoptions","nocxsptable","noextrasprl ","nofeaturethrottle","nogrowautofit","noindnmbrts","nojkernpunct","nolead","nolnhtadjtbl","nospaceforul","notabind ","notbrkcnstfrctbl","notcvasp","notvatxbx","nouicompat","noultrlspc","noxlattoyen","ogutterN","oldas","oldlinewrap","otblrul ","paperhN","paperwN","pgbrdrb","pgbrdrfoot","pgbrdrhead","pgbrdrl","pgbrdroptN","pgbrdrr","pgbrdrsnap","pgbrdrt","pgnstartN","prcolbl ","printdata ","private","protlevelN","psover","pszN ","readonlyrecommended","readprot","relyonvmlN","remdttm","rempersonalinfo","revbarN","revisions","revpropN","revprot ","rtldoc","rtlgutter","saveinvalidxml","saveprevpict","showplaceholdtextN","showxmlerrorsN","snaptogridincell","spltpgpar","splytwnine","sprsbsp","sprslnsp","sprsspbf ","sprstsm","sprstsp ","stylelock","stylelockbackcomp","stylelockenforced","stylelockqfset","stylelocktheme","stylesortmethodN","subfontbysize","swpbdr ","template","themelangcsN","themelangfeN","themelangN","toplinepunct","trackformattingN","trackmovesN","transmf ","truncatefontheight","truncex","tsd","twoonone","useltbaln","usenormstyforlist","usexform","utinl","validatexmlN","vertdoc","viewbkspN","viewkindN","viewnobound","viewscaleN","viewzkN","wgrffmtfilter","widowctrl","windowcaption","wpjst","wpsp","wraptrsp ","writereservation","writereservhash","wrppunct","xform":

	// Document Variables
	// case "docvar":

	// Drawing Object Properties
	// case "hl","hlfr","hlloc","hlsrc","hrule","hsv":

	// Drawing Objects
	// case "do ","dobxcolumn ","dobxmargin ","dobxpage ","dobymargin ","dobypage ","dobypara ","dodhgtN ","dolock ","dpaendhol ","dpaendlN ","dpaendsol ","dpaendwN ","dparc ","dparcflipx ","dparcflipy ","dpastarthol ","dpastartlN ","dpastartsol ","dpastartwN ","dpcallout ","dpcoaccent ","dpcoaN ","dpcobestfit ","dpcoborder ","dpcodabs","dpcodbottom ","dpcodcenter ","dpcodescentN","dpcodtop ","dpcolengthN ","dpcominusx ","dpcominusy ","dpcooffsetN ","dpcosmarta ","dpcotdouble ","dpcotright ","dpcotsingle ","dpcottriple ","dpcountN ","dpellipse ","dpendgroup ","dpfillbgcbN ","dpfillbgcgN ","dpfillbgcrN ","dpfillbggrayN ","dpfillbgpal ","dpfillfgcbN ","dpfillfgcgN ","dpfillfgcrN ","dpfillfggrayN ","dpfillfgpal ","dpfillpatN ","dpgroup ","dpline ","dplinecobN ","dplinecogN ","dplinecorN ","dplinedado ","dplinedadodo ","dplinedash ","dplinedot ","dplinegrayN ","dplinehollow ","dplinepal ","dplinesolid ","dplinewN ","dppolycountN ","dppolygon ","dppolyline ","dpptxN ","dpptyN ","dprect ","dproundr ","dpshadow ","dpshadxN ","dpshadyN ","dptxbtlr","dptxbx ","dptxbxmarN ","dptxbxtext ","dptxlrtb","dptxlrtbv","dptxtbrl","dptxtbrlv","dpxN ","dpxsizeN ","dpyN ","dpysizeN ":

	// East Asian Control Words
	// case "cgridN","g","gcwN","gridtbl","nosectexpand","ulhair":

	// Fields
	// case "datafield ","date","field","fldalt ","flddirty","fldedit","fldinst","fldlock","fldpriv","fldrslt","fldtype","time","wpeqn":
	case "fldrslt":
		text.WriteString(param)

	// File Table
	// case "fidN ","file ","filetbl ","fnetwork ","fnonfilesys","fosnumN ","frelativeN ","fvaliddos ","fvalidhpfs ","fvalidmac ","fvalidntfs ":

	// Font (Character) Formatting Properties
	case "acccircle", "acccomma", "accdot", "accnone", "accunderdot", "animtextN", "b", "caps", "cbN", "cchsN ", "cfN", "charscalexN", "csN", "dnN", "embo", "expndN", "expndtwN ", "fittextN", "fN", "fsN", "i", "kerningN ", "langfeN", "langfenpN", "langN", "langnpN", "ltrch", "noproof", "nosupersub ", "outl", "plain", "rtlch", "scaps", "shad", "strike", "sub ", "super ", "ul", "ulcN", "uld", "uldash", "uldashd", "uldashdd", "uldb", "ulhwave", "ulldash", "ulnone", "ulth", "ulthd", "ulthdash", "ulthdashd", "ulthdashdd", "ulthldash", "ululdbwave", "ulw", "ulwave", "upN", "v", "webhidden":
		text.WriteString(param)

	// Font Family
	// case "fjgothic","fjminchou","jis","falt ","fbiasN","fbidi","fcharsetN","fdecor","fetch","fmodern","fname","fnil","fontemb","fontfile","fonttbl","fprqN ","froman","fscript","fswiss","ftech","ftnil","fttruetype","panose":

	// Footnotes
	// case "footnote":

	// Form Fields
	// case "ffdefresN","ffdeftext","ffentrymcr","ffexitmcr","ffformat","ffhaslistboxN","ffhelptext","ffhpsN","ffl","ffmaxlenN","ffname","ffownhelpN","ffownstatN","ffprotN","ffrecalcN","ffresN","ffsizeN","ffstattext","fftypeN","fftypetxtN","formfield":

	// Generator
	// case "generator":

	// Headers and Footers
	// case "footer","footerf","footerl","footerr","header","headerf","headerl","headerr":

	// Highlighting
	// case "highlightN":

	// Hyphenation Information
	// case "chhresN","hresN":

	// Index Entries
	// case "bxe","ixe","pxe","rxe","txe","xe","xefN","yxe":

	// Information Group
	// case "author","buptim","category","comment","company","creatim","doccomm","dyN","edminsN","hlinkbase","hrN","idN","info","keywords","linkval","manager","minN","moN","nofcharsN","nofcharswsN","nofpagesN","nofwordsN","operator","printim","propname","proptypeN","revtim","secN","staticval","subject","title","userprops","vernN","versionN","yrN":

	// List Levels
	// case "lvltentative":

	// List Table
	// case "jclisttab","levelfollowN","levelindentN","leveljcN","leveljcnN","levellegalN","levelnfcN","levelnfcnN","levelnorestartN","levelnumbers","leveloldN","levelpictureN","levelpicturenosize","levelprevN","levelprevspaceN","levelspaceN","levelstartatN","leveltemplateidN","leveltext","lfolevel","list","listhybrid","listidN","listlevel","listname","listoverride","listoverridecountN","listoverrideformatN","listoverridestartat","listoverridetable","listpicture","listrestarthdnN","listsimpleN","liststyleidN","liststylename","listtable","listtemplateidN","lsN":

	// Macintosh Edition Manager Publisher Objects
	// case "bkmkpub","pubauto":

	// Mail Merge
	// case "mailmerge","mmaddfieldname","mmattach","mmblanklines","mmconnectstr","mmconnectstrdata","mmdatasource","mmdatatypeaccess","mmdatatypeexcel","mmdatatypefile","mmdatatypeodbc","mmdatatypeodso","mmdatatypeqt","mmdefaultsql","mmdestemail","mmdestfax","mmdestnewdoc 2 007","mmdestprinter","mmerrorsN","mmfttypeaddress","mmfttypebarcode","mmfttypedbcolumn","mmfttypemapped","mmfttypenull","mmfttypesalutation","mmheadersource","mmjdsotypeN","mmlinktoquery","mmmailsubject","mmmaintypecatalog","mmmaintypeemail","mmmaintypeenvelopes","mmmaintypefax","mmmaintypelabels","mmmaintypeletters","mmodso","mmodsoactiveN","mmodsocoldelimN","mmodsocolumnN","mmodsodynaddrN","mmodsofhdrN","mmodsofilter","mmodsofldmpdata","mmodsofmcolumnN","mmodsohashN","mmodsolidN","mmodsomappedname","mmodsoname","mmodsorecipdata","mmodsosort","mmodsosrc ","mmodsotable","mmodsoudl","mmodsoudldata 200 7","mmodsouniquetag","mmquery","mmreccurN","mmshowdata":

	// Math
	// case "macc","maccPr","maln","malnScr","margPr","margSzN","mbar","mbarPr","mbaseJc","mbegChr","mborderBox","mborderBoxPr","mbox","mboxPr","mbrkBinN","mbrkBinSubN","mbrkN","mcGpN","mcGpRuleN","mchr","mcount","mcSpN","mctrlPr","md","mdefJcN","mdeg","mdegHide","mden","mdiff","mdiffStyN","mdispdefN","mdPr","me","mendChr","meqArr","meqArrPr","mf","mfName","mfPr","mfunc","mfuncPr","mgroupChr","mgroupChrPr","mgrow","mhideBot","mhideLeft","mhideRight","mhideTop","minterSpN","mintLimN","mintraSpN","mjcN","mlim","mlimloc","mlimlow","mlimlowPr","mlimupp","mlimuppPr","mlit","mlMarginN","mm","mmath","mmathFontN","mmathPict","mmathPr","mmaxdist","mmc","mmcJc","mmcPr","mmcs","mmPr","mmr","mnary","mnaryLimN","mnaryPr","mnoBreak","mnor","mnum","mobjDist","moMath","moMathPara","moMathParaPr","mopEmu","mphant","mphantPr","mplcHide","mpos","mpostSpN","mpreSpN","mr","mrad","mradPr","mrMarginN","mrPr","mrSpN","mrSpRuleN","mscrN","msepChr","mshow","mshp","msmallFracN","msPre","msPrePr","msSub","msSubPr","msSubSup","msSubSupPr","msSup","msSupPr","mstrikeBLTR","mstrikeH","mstrikeTLBR","mstrikeV","mstyN","msub","msubHide","msup","msupHide","mtransp","mtype","mvertJc","mwrapIndentN","mwrapRightN","mzeroAsc","mzeroDesc","mzeroWid":

	// Microsoft Office Outlook
	// case "ebcstart","ebcend":

	// Move Bookmarks
	// case "mvfmf","mvfml","mvtof","mvtol":

	// New Asia Control Words Created by Word
	// case "horzvertN","twoinoneN":

	// Objects
	// case "linkself","objalias","objalignN","objattph","objautlink","objclass","objcropbN","objcroplN","objcroprN","objcroptN","objdata","object","objemb","objhN","objhtml","objicemb","objlink","objlock","objname","objocx","objpub","objscalexN","objscaleyN","objsect","objsetsize","objsub","objtime","objtransyN","objupdate ","objwN","oleclsid","result","rsltbmp","rslthtml","rsltmerge","rsltpict","rsltrtf","rslttxt":

	// Paragraph Borders
	// case "box","brdrb","brdrbar","brdrbtw","brdrcfN","brdrdash ","brdrdashd","brdrdashdd","brdrdashdot","brdrdashdotdot","brdrdashdotstr","brdrdashsm","brdrdb","brdrdot","brdremboss","brdrengrave","brdrframe","brdrhair","brdrinset","brdrl","brdrnil","brdrnone","brdroutset","brdrr","brdrs","brdrsh","brdrt","brdrtbl","brdrth","brdrthtnlg","brdrthtnmg","brdrthtnsg","brdrtnthlg","brdrtnthmg","brdrtnthsg","brdrtnthtnlg","brdrtnthtnmg","brdrtnthtnsg","brdrtriple","brdrwavy","brdrwavydb","brdrwN","brspN":

	// Paragraph Formatting Properties
	case "aspalpha", "aspnum", "collapsed", "contextualspace", "cufiN", "culiN", "curiN", "faauto", "facenter", "fafixed", "fahang", "faroman", "favar", "fiN", "hyphpar ", "indmirror", "intbl", "itapN", "keep", "keepn", "levelN", "liN", "linN", "lisaN", "lisbN", "ltrpar", "nocwrap", "noline", "nooverflow", "nosnaplinegrid", "nowidctlpar ", "nowwrap", "outlinelevelN ", "pagebb", "pard", "prauthN", "prdateN", "qc", "qd", "qj", "qkN", "ql", "qr", "qt", "riN", "rinN", "rtlpar", "saautoN", "saN", "sbautoN", "sbN", "sbys", "slmultN", "slN", "sN", "spv", "subdocumentN ", "tscbandhorzeven", "tscbandhorzodd", "tscbandverteven", "tscbandvertodd", "tscfirstcol", "tscfirstrow", "tsclastcol", "tsclastrow", "tscnecell", "tscnwcell", "tscsecell", "tscswcell", "txbxtwalways", "txbxtwfirst", "txbxtwfirstlast", "txbxtwlast", "txbxtwno", "widctlpar", "ytsN":
		text.WriteString(param)

	// Paragraph Group Properties
	// case "pgp","pgptbl","ipgpN":

	// Paragraph Revision Mark Properties
	// case "dfrauthN","dfrdateN","dfrstart","dfrstop","dfrxst":

	// Paragraph Shading
	// case "bgbdiag","bgcross","bgdcross","bgdkbdiag","bgdkcross","bgdkdcross","bgdkfdiag","bgdkhoriz","bgdkvert","bgfdiag","bghoriz","bgvert","cbpatN","cfpatN","shadingN":

	// Pictures
	// case "binN","bliptagN","blipuid","blipupiN","defshp","dibitmapN","emfblip","jpegblip","macpict","nonshppict","picbmp ","picbppN ","piccropbN","piccroplN","piccroprN","piccroptN","pichgoalN","pichN","picprop","picscaled","picscalexN","picscaleyN","pict","picwgoalN","picwN","pmmetafileN","pngblip","shppict","wbitmapN","wbmbitspixelN","wbmplanesN","wbmwidthbyteN","wmetafileN":

	// Positioned Objects and Frames
	// case "abshN","abslock","absnoovrlpN","abswN","dfrmtxtxN","dfrmtxtyN","dropcapliN ","dropcaptN ","dxfrtextN","frmtxbtlr","frmtxlrtb","frmtxlrtbv","frmtxtbrl","frmtxtbrlv","nowrap","overlay","phcol","phmrg","phpg","posnegxN ","posnegyN ","posxc","posxi","posxl","posxN","posxo","posxr","posyb","posyc","posyil","posyin","posyN","posyout","posyt","pvmrg","pvpara","pvpg","wraparound","wrapdefault","wrapthrough","wraptight":

	// Protection Exceptions
	// case "protend","protstart":

	// Quick Styles
	// case "noqfpromote":

	// Read-Only Password Protection
	// case "password","passwordhash":

	// Revision Marks for Paragraph Numbers and ListNum Fields
	// case "pnrauthN","pnrdateN","pnrnfcN","pnrnot","pnrpnbrN","pnrrgbN","pnrstartN","pnrstopN","pnrxstN":

	// RTF Version
	// case "rtfN":

	// Section Formatting Properties
	case "adjustright", "binfsxnN", "binsxnN", "colnoN ", "colsN", "colsrN ", "colsxN", "colwN ", "dsN", "endnhere", "footeryN", "guttersxnN", "headeryN", "horzsect", "linebetcol", "linecont", "linemodN", "lineppage", "linerestart", "linestartsN", "linexN", "lndscpsxn", "ltrsect", "margbsxnN", "marglsxnN", "margmirsxn", "margrsxnN", "margtsxnN", "pghsxnN", "pgnbidia", "pgnbidib", "pgnchosung", "pgncnum", "pgncont", "pgndbnum", "pgndbnumd", "pgndbnumk", "pgndbnumt", "pgndec", "pgndecd", "pgnganada", "pgngbnum", "pgngbnumd", "pgngbnumk", "pgngbnuml", "pgnhindia", "pgnhindib", "pgnhindic", "pgnhindid", "pgnhnN ", "pgnhnsc ", "pgnhnsh ", "pgnhnsm ", "pgnhnsn ", "pgnhnsp ", "pgnid", "pgnlcltr", "pgnlcrm", "pgnrestart", "pgnstartsN", "pgnthaia", "pgnthaib", "pgnthaic", "pgnucltr", "pgnucrm", "pgnvieta", "pgnxN", "pgnyN", "pgnzodiac", "pgnzodiacd", "pgnzodiacl", "pgwsxnN", "pnseclvlN", "rtlsect", "saftnnalc", "saftnnar", "saftnnauc", "saftnnchi", "saftnnchosung", "saftnncnum", "saftnndbar", "saftnndbnum", "saftnndbnumd", "saftnndbnumk", "saftnndbnumt", "saftnnganada", "saftnngbnum", "saftnngbnumd", "saftnngbnumk", "saftnngbnuml", "saftnnrlc", "saftnnruc", "saftnnzodiac", "saftnnzodiacd", "saftnnzodiacl", "saftnrestart", "saftnrstcont", "saftnstartN", "sbkcol", "sbkeven", "sbknone", "sbkodd", "sbkpage", "sectd", "sectdefaultcl", "sectexpandN", "sectlinegridN", "sectspecifycl", "sectspecifygenN", "sectspecifyl", "sectunlocked", "sftnbj", "sftnnalc", "sftnnar", "sftnnauc", "sftnnchi", "sftnnchosung", "sftnncnum", "sftnndbar", "sftnndbnum", "sftnndbnumd", "sftnndbnumk", "sftnndbnumt", "sftnnganada", "sftnngbnum", "sftnngbnumd", "sftnngbnumk", "sftnngbnuml", "sftnnrlc", "sftnnruc", "sftnnzodiac", "sftnnzodiacd", "sftnnzodiacl", "sftnrestart", "sftnrstcont", "sftnrstpg", "sftnstartN", "sftntj", "srauthN", "srdateN", "titlepg", "vertal", "vertalb", "vertalc", "vertalj", "vertalt", "vertsect":
		text.WriteString(param)

	// Section Text
	case "stextflowN":
		text.WriteString(param)

	// SmartTag Data
	// case "factoidname":

	// Special Characters
	case "-", ":", "_", "{", "|", "}", "~", "bullet", "chatn", "chdate", "chdpa", "chdpl", "chftn", "chftnsep", "chftnsepc", "chpgn", "chtime", "column", "emdash", "emspace ", "endash", "enspace ", "lbrN", "ldblquote", "line", "lquote", "ltrmark", "page", "par", "qmspace", "rdblquote", "row", "rquote", "rtlmark", "sect", "sectnum", "softcol ", "softlheightN ", "softline ", "softpage ", "tab", "zwbo", "zwj", "zwnbo", "zwnj":
		text.WriteString(param)

	// Style and Formatting Restrictions
	// case "latentstyles","lsdlockeddefN","lsdlockedexcept","lsdlockedN","lsdprioritydefN","lsdpriorityN","lsdqformatdefN","lsdqformatN","lsdsemihiddendefN","lsdsemihiddenN","lsdstimaxN","lsdunhideuseddefN","lsdunhideusedN":

	// Style Sheet
	// case "additive","alt","ctrl","fnN","keycode","sautoupd","sbasedonN","scompose","shidden","shift","slinkN","slocked","snextN","spersonal","spriorityN","sqformat","sreply","ssemihiddenN","stylesheet","styrsidN","sunhideusedN","tsN","tsrowd":

	// Table Definitions
	case "cell", "cellxN", "clbgbdiag", "clbgcross", "clbgdcross", "clbgdkbdiag", "clbgdkcross", "clbgdkdcross", "clbgdkfdiag", "clbgdkhor", "clbgdkvert", "clbgfdiag", "clbghoriz", "clbgvert", "clbrdrb", "clbrdrl", "clbrdrr", "clbrdrt", "clcbpatN", "clcbpatrawN", "clcfpatN", "clcfpatrawN", "cldel2007", "cldelauthN", "cldeldttmN", "cldgll", "cldglu", "clFitText", "clftsWidthN", "clhidemark", "clins", "clinsauthN", "clinsdttmN", "clmgf", "clmrg", "clmrgd", "clmrgdauthN", "clmrgddttmN", "clmrgdr", "clNoWrap", "clpadbN", "clpadfbN", "clpadflN", "clpadfrN", "clpadftN", "clpadlN", "clpadrN", "clpadtN", "clshdngN", "clshdngrawN", "clshdrawnil", "clspbN", "clspfbN", "clspflN", "clspfrN", "clspftN", "clsplit", "clsplitr", "clsplN", "clsprN", "clsptN", "cltxbtlr", "cltxlrtb", "cltxlrtbv", "cltxtbrl", "cltxtbrlv", "clvertalb", "clvertalc", "clvertalt", "clvmgf", "clvmrg", "clwWidthN", "irowbandN", "irowN", "lastrow", "ltrrow", "nestcell", "nestrow", "nesttableprops", "nonesttables", "rawclbgbdiag", "rawclbgcross", "rawclbgdcross", "rawclbgdkbdiag", "rawclbgdkcross", "rawclbgdkdcross", "rawclbgdkfdiag", "rawclbgdkhor", "rawclbgdkvert", "rawclbgfdiag", "rawclbghoriz", "rawclbgvert", "rtlrow", "tabsnoovrlp", "taprtl", "tblindN", "tblindtypeN", "tbllkbestfit", "tbllkborder", "tbllkcolor", "tbllkfont", "tbllkhdrcols", "tbllkhdrrows", "tbllklastcol", "tbllklastrow", "tbllknocolband", "tbllknorowband", "tbllkshading", "tcelld", "tdfrmtxtBottomN", "tdfrmtxtLeftN", "tdfrmtxtRightN", "tdfrmtxtTopN", "tphcol", "tphmrg", "tphpg", "tposnegxN", "tposnegyN", "tposxc", "tposxi", "tposxl", "tposxN", "tposxo", "tposxr", "tposyb", "tposyc", "tposyil", "tposyin", "tposyN", "tposyout", "tposyt", "tpvmrg", "tpvpara", "tpvpg", "trauthN", "trautofitN", "trbgbdiag", "trbgcross", "trbgdcross", "trbgdkbdiag", "trbgdkcross", "trbgdkdcross", "trbgdkfdiag", "trbgdkhor", "trbgdkvert", "trbgfdiag", "trbghoriz", "trbgvert", "trbrdrb ", "trbrdrh ", "trbrdrl ", "trbrdrr ", "trbrdrt ", "trbrdrv ", "trcbpatN", "trcfpatN", "trdateN", "trftsWidthAN", "trftsWidthBN", "trftsWidthN", "trgaphN", "trhdr ", "trkeep ", "trkeepfollow", "trleftN", "trowd", "trpaddbN", "trpaddfbN", "trpaddflN", "trpaddfrN", "trpaddftN", "trpaddlN", "trpaddrN", "trpaddtN", "trpadobN", "trpadofbN", "trpadoflN", "trpadofrN", "trpadoftN", "trpadolN", "trpadorN", "trpadotN", "trpatN", "trqc", "trql", "trqr", "trrhN", "trshdngN", "trspdbN", "trspdfbN", "trspdflN", "trspdfrN", "trspdftN", "trspdlN", "trspdrN", "trspdtN", "trspobN", "trspofbN", "trspoflN", "trspofrN", "trspoftN", "trspolN", "trsporN", "trspotN", "trwWidthAN", "trwWidthBN", "trwWidthN":
		text.WriteString(param)

	// Table of Contents Entries
	case "tc", "tcfN", "tclN", "tcn ":
		text.WriteString(param)

	// Table Styles
	// case "tsbgbdiag","tsbgcross","tsbgdcross","tsbgdkbdiag","tsbgdkcross","tsbgdkdcross","tsbgdkfdiag","tsbgdkhor","tsbgdkvert","tsbgfdiag","tsbghoriz","tsbgvert","tsbrdrb","tsbrdrdgl","tsbrdrdgr","tsbrdrh","tsbrdrl","tsbrdrr","tsbrdrr","tsbrdrt","tsbrdrv","tscbandshN","tscbandsvN","tscellcbpatN","tscellcfpatN","tscellpaddbN","tscellpaddfbN","tscellpaddflN","tscellpaddfrN","tscellpaddftN","tscellpaddlN","tscellpaddrN","tscellpaddtN","tscellpctN","tscellwidthftsN","tscellwidthN","tsnowrap","tsvertalb","tsvertalc","tsvertalt":

	// Tabs
	case "tbN", "tldot", "tleq", "tlhyph", "tlmdot", "tlth", "tlul", "tqc", "tqdec", "tqr", "txN":
		text.WriteString(param)

	// Theme Data
	// case "themedata":

	// Theme Font Information
	// case "fbimajor","fbiminor","fdbmajor","fdbminor","fhimajor","fhiminor","flomajor","flominor":

	// Track Changes
	// case "revtbl ":

	// Track Changes (Revision Marks)
	// case "charrsidN","delrsidN","insrsidN","oldcprops","oldpprops","oldsprops","oldtprops","pararsidN","rsidN","rsidrootN","rsidtbl","sectrsidN","tblrsidN":

	// Unicode RTF
	// case "ucN","ud","uN","upr":

	// User Protection Information
	// case "protusertbl":

	// Word through Word RTF for Drawing Objects (Shapes)
	// case "shp","shpbottomN","shpbxcolumn","shpbxignore","shpbxmargin","shpbxpage","shpbyignore","shpbymargin","shpbypage","shpbypara","shpfblwtxtN","shpfhdrN","shpgrp","shpinst","shpleftN","shplidN","shplockanchor","shprightN","shprslt","shptopN","shptxt","shpwrkN","shpwrN","shpzN","sn","sp","sv","svb":
	default:
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
