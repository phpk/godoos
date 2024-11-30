/*
Type definitions for markdown elements.
*/
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
package libs

import "fmt"

//go:generate stringer -type=Kind
type Kind int

//go:generate stringer -type=ElementType
type ElementType int

// specific types
const (
	// block types
	Head Kind = iota
	Paragraph
	List
	QuoteBlock
	CodeBlock
	Rule
	// inline types
	Emphasis
	Strong
	Link
	Code
	Image
)

// element types
const (
	Block ElementType = iota
	Inline
)

const _Kind_name = "HeadParagraphListQuoteBlockCodeBlockRuleEmphasisStrongLinkCodeImage"

var _Kind_index = [...]uint8{4, 13, 17, 27, 36, 40, 48, 54, 58, 62, 67}

func (i Kind) String() string {
	if i < 0 || i >= Kind(len(_Kind_index)) {
		return fmt.Sprintf("Kind(%d)", i)
	}
	hi := _Kind_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _Kind_index[i-1]
	}
	return _Kind_name[lo:hi]
}

const _ElementType_name = "BlockInline"

var _ElementType_index = [...]uint8{5, 11}

func (i ElementType) String() string {
	if i < 0 || i >= ElementType(len(_ElementType_index)) {
		return fmt.Sprintf("ElementType(%d)", i)
	}
	hi := _ElementType_index[i]
	lo := uint8(0)
	if i > 0 {
		lo = _ElementType_index[i-1]
	}
	return _ElementType_name[lo:hi]
}
