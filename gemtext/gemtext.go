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
var cr = []byte{'\r'}

// Parse parses src as text/gemini and returns the resulting document.
// It does not check src to contain valid utf8 sequences.
func Parse(src []byte) Document {
	doc := Document{}
	var pre Preformatted
	inPre := false

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
			pre = Preformatted{AltText: string(line[len(fenceBytes):])}
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
	line, rest, _ = bytes.Cut(src, []byte{'\n'})
	return bytes.TrimSuffix(line, cr), rest
}

func parseLine(line []byte) Line {
	if len(line) == 0 {
		return Text{}
	}
	switch line[0] {
	case '=':
		if len(line) > 1 && line[1] == '>' {
			return parseLink(line[2:])
		}
	case '#':
		return parseHeading(line[1:])
	case '*':
		if len(line) > 1 && line[1] == ' ' {
			return ListItem{Text: string(line[2:])}
		}
	case '>':
		return Quote{Text: string(line[1:])}
	}
	return Text{Text: string(line)}
}

func parseLink(rest []byte) Link {
	rest = bytes.TrimLeft(rest, whitespace)
	url, label := rest, []byte(nil)
	if i := bytes.IndexAny(rest, whitespace); i >= 0 {
		url, label = rest[:i], bytes.TrimLeft(rest[i:], whitespace)
	}
	return Link{URL: string(url), Label: string(label)}
}

func parseHeading(rest []byte) Heading {
	level := 1
	if len(rest) > 0 && rest[0] == '#' {
		level = 2
		rest = rest[1:]
	}
	if len(rest) > 0 && rest[0] == '#' {
		level = 3
		rest = rest[1:]
	}
	return Heading{Level: level, Text: string(rest)}
}
