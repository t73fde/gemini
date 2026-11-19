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

// Document is the parsed form of a text/gemini document.
type Document struct {
	Lines []Line
}

// Line is the interface implemented by all line types.
type Line interface {
	isLine()
}

// Text is a plain text, non-empty line.
type Text struct{ Text string }

// Link is a link line: "=> URL [label]".
type Link struct{ URL, Label string }

// Heading is a heading line of level 1 to 3.
type Heading struct {
	Level int
	Text  string
}

// ListItem is a single list item line: "* text".
type ListItem struct{ Text string }

// Quote is a single quote line: "> text".
type Quote struct{ Text string }

// Preformatted is a block of verbatim lines between two "```" lines.
type Preformatted struct {
	AltText string
	Lines   []string
}

func (Text) isLine()         {}
func (Link) isLine()         {}
func (Heading) isLine()      {}
func (ListItem) isLine()     {}
func (Quote) isLine()        {}
func (Preformatted) isLine() {}
