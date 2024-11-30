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

/* I don't think I'm going to need this
type plcFld struct {
	aCp  []int
	aFld []fld
}

type fld struct {
	fldch     int
	grffld    int
	fieldtype string
	fNested   bool
	fHasSep   bool
}

func getPlcFld(table *mscfb.File, offset, size int) (*plcFld, error) {
	if table == nil {
		return nil, errInvalidArgument
	}
	b := make([]byte, size)
	_, err := table.ReadAt(b, int64(offset))
	if err != nil {
		return nil, err
	}

	f, err := getFld(b)
	if err != nil {
		return nil, err
	}

	return f, nil
}

func getFld(plc []byte) (*plcFld, error) {
	return nil, nil
}

func getFieldType(grffld byte) string {
	switch grffld {
	case 0x01:
		return "Not Named"
	case 0x02:
		return "Not Named"
	case 0x03:
		return "REF"
	case 0x05:
		return "FTNREF"
	case 0x06:
		return "SET"
	case 0x07:
		return "IF"
	case 0x08:
		return "INDEX"
	case 0x0A:
		return "STYLEREF"
	case 0x0C:
		return "SEQ"
	case 0x0D:
		return "TOC"
	case 0x0E:
		return "INFO"
	case 0x0F:
		return "TITLE"
	case 0x10:
		return "SUBJECT"
	case 0x11:
		return "AUTHOR"
	case 0x12:
		return "KEYWORDS"
	case 0x13:
		return "COMMENTS"
	case 0x14:
		return "LASTSAVEDBY"
	case 0x15:
		return "CREATEDATE"
	case 0x16:
		return "SAVEDATE"
	case 0x17:
		return "PRINTDATE"
	case 0x18:
		return "REVNUM"
	case 0x19:
		return "EDITTIME"
	case 0x1A:
		return "NUMPAGES"
	case 0x1B:
		return "NUMWORDS"
	case 0x1C:
		return "NUMCHARS"
	case 0x1D:
		return "FILENAME"
	case 0x1E:
		return "TEMPLATE"
	case 0x1F:
		return "DATE"
	case 0x20:
		return "TIME"
	case 0x21:
		return "PAGE"
	case 0x22:
		return "="
	case 0x23:
		return "QUOTE"
	case 0x24:
		return "INCLUDE"
	case 0x25:
		return "PAGEREF"
	case 0x26:
		return "ASK"
	case 0x27:
		return "FILLIN"
	case 0x28:
		return "DATA"
	case 0x29:
		return "NEXT"
	case 0x2A:
		return "NEXTIF"
	case 0x2B:
		return "SKIPIF"
	case 0x2C:
		return "MERGEREC"
	case 0x2D:
		return "DDE"
	case 0x2E:
		return "DDEAUTO"
	case 0x2F:
		return "GLOSSARY"
	case 0x30:
		return "PRINT"
	case 0x31:
		return "EQ"
	case 0x32:
		return "GOTOBUTTON"
	case 0x33:
		return "MACROBUTTON"
	case 0x34:
		return "AUTONUMOUT"
	case 0x35:
		return "AUTONUMLGL"
	case 0x36:
		return "AUTONUM"
	case 0x37:
		return "IMPORT"
	case 0x38:
		return "LINK"
	case 0x39:
		return "SYMBOL"
	case 0x3A:
		return "EMBED"
	case 0x3B:
		return "MERGEFIELD"
	case 0x3C:
		return "USERNAME"
	case 0x3D:
		return "USERINITIALS"
	case 0x3E:
		return "USERADDRESS"
	case 0x3F:
		return "BARCODE"
	case 0x40:
		return "DOCVARIABLE"
	case 0x41:
		return "SECTION"
	case 0x42:
		return "SECTIONPAGES"
	case 0x43:
		return "INCLUDEPICTURE"
	case 0x44:
		return "INCLUDETEXT"
	case 0x45:
		return "FILESIZE"
	case 0x46:
		return "FORMTEXT"
	case 0x47:
		return "FORMCHECKBOX"
	case 0x48:
		return "NOTEREF"
	case 0x49:
		return "TOA"
	case 0x4B:
		return "MERGESEQ"
	case 0x4F:
		return "AUTOTEXT"
	case 0x50:
		return "COMPARE"
	case 0x51:
		return "ADDIN"
	case 0x53:
		return "FORMDROPDOWN"
	case 0x54:
		return "ADVANCE"
	case 0x55:
		return "DOCPROPERTY"
	case 0x57:
		return "CONTROL"
	case 0x58:
		return "HYPERLINK"
	case 0x59:
		return "AUTOTEXTLIST"
	case 0x5A:
		return "LISTNUM"
	case 0x5B:
		return "HTMLCONTROL"
	case 0x5C:
		return "BIDIOUTLINE"
	case 0x5D:
		return "ADDRESSBLOCK"
	case 0x5E:
		return "GREETINGLINE"
	case 0x5F:
		return "SHAPE"
	default:
		return "UNKNOWN"
	}
}
*/
