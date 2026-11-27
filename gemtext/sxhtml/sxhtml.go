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

// Package sxhtml converts a sequence of Gemtext lines (see package
// t73f.de/r/gemini/gemtext) into SxHTML block elements.
//
// The conversion maps each line type as follows:
//
//   - Text line: (p "text")
//   - Heading level 1 to 3: (h1 "text") to (h3 "text"); other levels are dropped
//   - Link line: (p (a ((href . "URL")) "label")); without a label, the URL
//     is used as the text
//   - Quote line: (blockquote "text"); every quote line is independent of its
//     neighbors
//   - Preformatted block: (pre (code "lines joined by LF")); the alt text is
//     currently not output
//   - Consecutive list items: (ul (li "item") ...); a list ends with the first
//     line that is not a list item
//
// Blank lines (empty or consisting of white space only) do not produce
// elements of their own, because all other elements are rendered as
// paragraph-like blocks. Instead, n >= 2 consecutive blank lines between two
// elements, or at the start of the input, produce n-1 (br) elements. Blank
// lines at the end of the input are dropped.
//
// Link URLs are not validated; invalid URLs are left to the browser.
package sxhtml

import (
	"strings"

	"t73f.de/r/gemini/gemtext"
	"t73f.de/r/sx"
)

// converter implements a two-state machine over the Gemtext lines.
// The state "list pending" is represented by inList. The state "no list
// pending" is represented by !inList, together with count, which holds the
// number of consecutive blank lines seen minus one (-1: none pending).
type converter struct {
	inList bool
	count  int // number of blank lines seen minus one; -1: no blank line pending
	items  sx.ListBuilder
	result sx.ListBuilder
}

// Convert transforms a sequence of Gemtext lines into SxHTML block elements.
func Convert(lines []gemtext.Line) *sx.Pair {
	c := converter{count: -1}
	c.items.Reset()
	c.result.Reset()
	for _, line := range lines {
		c.feed(line)
	}
	// EOD: an open list is emitted, pending blank lines are dropped.
	if c.inList {
		c.emitList()
	}
	return c.result.List()
}

var (
	symA          = sx.MakeSymbol("a")
	symBLOCKQUOTE = sx.MakeSymbol("blockquote")
	symBR         = sx.MakeSymbol("br")
	symCODE       = sx.MakeSymbol("code")
	symH1         = sx.MakeSymbol("h1")
	symH2         = sx.MakeSymbol("h2")
	symH3         = sx.MakeSymbol("h3")
	symHREF       = sx.MakeSymbol("href")
	symLI         = sx.MakeSymbol("li")
	symP          = sx.MakeSymbol("p")
	symPRE        = sx.MakeSymbol("pre")
	symUL         = sx.MakeSymbol("ul")
)

func (c *converter) feed(line gemtext.Line) {
	item, isItem := line.(gemtext.ListItem)

	if c.inList {
		if isItem {
			c.items.Add(sx.MakeList(symLI, sx.MakeString(item.Text)))
			return // AddItem
		}
		c.emitList() // resets inList, count stays -1
	}

	if isBlank(line) {
		c.count++
		return
	}
	c.emitBlank()

	if isItem {
		c.items.Reset() // NewList
		c.items.Add(sx.MakeList(symLI, sx.MakeString(item.Text)))
		c.inList = true
		return
	}
	c.emitBlock(line)
}

func (c *converter) emitBlank() {
	for ; c.count > 0; c.count-- {
		c.result.Add(sx.MakeList(symBR))
	}
	c.count = -1
}

func (c *converter) emitList() {
	c.result.Add(c.items.List().Cons(symUL))
	c.items.Reset()
	c.inList = false
}

func (c *converter) emitBlock(line gemtext.Line) {
	var obj sx.Object
	switch ln := line.(type) {
	case gemtext.Text:
		obj = sx.MakeList(symP, sx.MakeString(ln.Text))
	case gemtext.Heading:
		var sym *sx.Symbol
		switch ln.Level {
		case 1:
			sym = symH1
		case 2:
			sym = symH2
		case 3:
			sym = symH3
		default:
			return
		}
		obj = sx.MakeList(sym, sx.MakeString(ln.Text))
	case gemtext.Link:
		label := ln.Label
		if label == "" {
			label = ln.URL
		}
		obj = sx.MakeList(
			symP,
			sx.MakeList(
				symA,
				sx.MakeList(sx.Cons(symHREF, sx.MakeString(ln.URL))),
				sx.MakeString(label),
			),
		)
	case gemtext.Quote:
		obj = sx.MakeList(symBLOCKQUOTE, sx.MakeString(ln.Text))
	case gemtext.Preformatted:
		// TODO: output ln.Alt as attribute (name still open)
		obj = sx.MakeList(
			symPRE,
			sx.MakeList(symCODE, sx.MakeString(strings.Join(ln.Lines, "\n"))),
		)
	default:
		return
	}
	c.result.Add(obj)
}

// isBlank reports whether line is empty or consists of white space only.
func isBlank(line gemtext.Line) bool {
	t, ok := line.(gemtext.Text)
	return ok && strings.TrimSpace(t.Text) == ""
}
