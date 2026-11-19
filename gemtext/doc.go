//-----------------------------------------------------------------------------
// Copyright (c) 2026-present Detlef Stern
//
// This file is part of Gemini.
//
// Gemini is licensed under the latest version of the EUPL (European Union
// Public License). Please see file LICENSE.txt for your rights and obligations
// under this license.
//
// SPDX-License-Identifier: EUPL-1.2
// SPDX-FileCopyrightText: 2026-present Detlef Stern
//-----------------------------------------------------------------------------

// Package gemtext parses text/gemini (Gemtext) documents into a
// line-oriented abstract syntax tree.
//
// Gemtext is defined by the Gemini specification, see
// https://geminiprotocol.net/docs/specification.gmi. The package recognizes
// the line types of the specification and represents each of them as one
// element of a [Document]:
//
//   - [Text]: plain text line, possibly empty
//   - [Link]: "=> URL [label]"
//   - [Heading]: "#", "##" or "###" followed by text
//   - [ListItem]: "* text"
//   - [Quote]: ">text"
//   - [Preformatted]: lines between two "```" lines, with optional alt text
//
// Parsing never fails. Input that does not match a more specific line type
// is a [Text] line. In detail:
//
//   - Lines end with LF or CRLF. A final line terminator does not start an
//     additional empty line.
//   - A line that starts with "=>" but has no URL is a [Text] line.
//   - At most three leading '#' characters are counted as the heading level.
//     Further '#' characters belong to the heading text.
//   - A "*" is a list item marker only if it is followed by a space.
//   - A preformatted block that is not closed by the end of the input is
//     closed implicitly.
//   - Whitespace means space and tab. The input is not validated to be UTF-8.
//
// All strings in the tree are copies; the tree does not refer to the input.
//
// Example:
//
//	doc := gemtext.Parse(src)
//	for _, line := range doc.Lines {
//		switch l := line.(type) {
//		case gemtext.Heading:
//			fmt.Println(l.Level, l.Text)
//		case gemtext.Link:
//			fmt.Println(l.URL, l.Label)
//		}
//	}
package gemtext
