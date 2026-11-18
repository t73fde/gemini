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

package gemtext

import "bytes"

const whitespace = " \t"

var fenceBytes = []byte("```")
var utf8bom = []byte("\xEF\xBB\xBF")
var cr = []byte{'\r'}

// Parse parses src as text/gemini and returns the resulting document.
func Parse(src []byte) Document {
	doc := Document{}
	var pre Preformatted
	inPre := false

	src = bytes.TrimPrefix(src, utf8bom)
	for len(src) > 0 {
		var line []byte
		line, src = nextLine(src)

		if inPre {
			if bytes.HasPrefix(line, fenceBytes) {
				doc.Lines = append(doc.Lines, pre)
				inPre = false
			} else {
				pre.Lines = append(pre.Lines, string(line))
			}
			continue
		}
		if bytes.HasPrefix(line, fenceBytes) {
			pre = Preformatted{AltText: string(bytes.Trim(line[len(fenceBytes):], whitespace))}
			inPre = true
			continue
		}
		doc.Lines = append(doc.Lines, parseLine(line))
	}
	if inPre {
		doc.Lines = append(doc.Lines, pre)
	}
	return doc
}

// nextLine splits src at the first line feed and removes a trailing carriage return.
func nextLine(src []byte) (line, rest []byte) {
	if i := bytes.IndexByte(src, '\n'); i >= 0 {
		line, rest = src[:i], src[i+1:]
	} else {
		line = src
	}
	return bytes.TrimSuffix(line, cr), rest
}

func parseLine(line []byte) Line {
	if len(bytes.TrimLeft(line, whitespace)) == 0 {
		return Blank{}
	}
	switch line[0] {
	case '=':
		if len(line) > 1 && line[1] == '>' {
			if l, ok := parseLink(line[2:]); ok {
				return l
			}
		}
	case '#':
		return parseHeading(line)
	case '*':
		if len(line) > 1 && line[1] == ' ' {
			return ListItem{Text: trimLeft(line[2:])}
		}
	case '>':
		return Quote{Text: trimLeft(line[1:])}
	}
	return Text{Text: string(line)}
}

func parseLink(rest []byte) (Link, bool) {
	rest = bytes.TrimLeft(rest, whitespace)
	if len(rest) == 0 {
		return Link{}, false
	}
	url, label := rest, []byte(nil)
	if i := bytes.IndexAny(rest, whitespace); i >= 0 {
		url, label = rest[:i], bytes.Trim(rest[i:], whitespace)
	}
	return Link{URL: string(url), Label: string(label)}, true
}

func parseHeading(line []byte) Heading {
	limit := min(3, len(line))
	level := 1
	for level < limit && line[level] == '#' {
		level++
	}
	return Heading{Level: level, Text: string(bytes.Trim(line[level:], whitespace))}
}

func trimLeft(b []byte) string {
	return string(bytes.TrimLeft(b, whitespace))
}
