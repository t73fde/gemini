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

package sxhtml_test

import (
	"testing"

	"t73f.de/r/gemini/gemtext"
	gmi2sxhtml "t73f.de/r/gemini/gemtext/sxhtml"
)

// unknownLine is a line type not known to the converter. Embedding promotes
// the methods of gemtext.Text, so unknownLine implements gemtext.Line, but it
// does not match the case gemtext.Text in a type switch.
type unknownLine struct{ gemtext.Text }

func text(s string) gemtext.Line { return gemtext.Text{Text: s} }
func item(s string) gemtext.Line { return gemtext.ListItem{Text: s} }

func TestConvert(t *testing.T) {
	blank := text("")
	testcases := []struct {
		name  string
		lines []gemtext.Line
		exp   string
	}{
		{"empty input", nil, "()"},
		{"text", []gemtext.Line{text("a")}, `((p "a"))`},

		// Blank lines
		{"one blank", []gemtext.Line{text("a"), blank, text("b")}, `((p "a") (p "b"))`},
		{"two blanks", []gemtext.Line{text("a"), blank, blank, text("b")}, `((p "a") (br) (p "b"))`},
		{"three blanks", []gemtext.Line{text("a"), blank, blank, blank, text("b")}, `((p "a") (br) (br) (p "b"))`},
		{"white space only is blank", []gemtext.Line{text("a"), text(" \t"), text(""), text("b")}, `((p "a") (br) (p "b"))`},
		{"leading one blank", []gemtext.Line{blank, text("a")}, `((p "a"))`},
		{"leading two blanks", []gemtext.Line{blank, blank, text("a")}, `((br) (p "a"))`},
		{"trailing blanks", []gemtext.Line{text("a"), blank, blank, blank}, `((p "a"))`},
		{"only blanks", []gemtext.Line{blank, blank}, "()"},

		// Headings
		{"heading 1", []gemtext.Line{gemtext.Heading{Level: 1, Text: "h"}}, `((h1 "h"))`},
		{"heading 2", []gemtext.Line{gemtext.Heading{Level: 2, Text: "h"}}, `((h2 "h"))`},
		{"heading 3", []gemtext.Line{gemtext.Heading{Level: 3, Text: "h"}}, `((h3 "h"))`},
		{"heading level 0 dropped", []gemtext.Line{text("a"), gemtext.Heading{Level: 0, Text: "h"}, text("b")}, `((p "a") (p "b"))`},
		{"heading level 4 dropped", []gemtext.Line{text("a"), gemtext.Heading{Level: 4, Text: "h"}, text("b")}, `((p "a") (p "b"))`},

		// Links
		{"link with label", []gemtext.Line{gemtext.Link{URL: "u", Label: "l"}}, `((p (a ((href . "u")) "l")))`},
		{"link without label", []gemtext.Line{gemtext.Link{URL: "u"}}, `((p (a ((href . "u")) "u")))`},

		// Quotes
		{"quote", []gemtext.Line{gemtext.Quote{Text: "q"}}, `((blockquote "q"))`},
		{"quotes are independent", []gemtext.Line{gemtext.Quote{Text: "q1"}, gemtext.Quote{Text: "q2"}}, `((blockquote "q1") (blockquote "q2"))`},

		// Preformatted
		{"preformatted", []gemtext.Line{gemtext.Preformatted{Lines: []string{"x"}}}, `((pre (code "x")))`},

		// Lists
		{"list at end of input", []gemtext.Line{item("a"), item("b")}, `((ul (li "a") (li "b")))`},
		{"list followed by text", []gemtext.Line{item("a"), item("b"), text("t")}, `((ul (li "a") (li "b")) (p "t"))`},
		{"list followed by heading", []gemtext.Line{item("a"), gemtext.Heading{Level: 1, Text: "h"}}, `((ul (li "a")) (h1 "h"))`},
		{"list after text", []gemtext.Line{text("t"), item("a")}, `((p "t") (ul (li "a")))`},
		{"list after one blank", []gemtext.Line{text("t"), blank, item("a")}, `((p "t") (ul (li "a")))`},
		{"list after two blanks", []gemtext.Line{text("t"), blank, blank, item("a")}, `((p "t") (br) (ul (li "a")))`},
		{"two lists, one blank", []gemtext.Line{item("a"), blank, item("b")}, `((ul (li "a")) (ul (li "b")))`},
		{"two lists, two blanks", []gemtext.Line{item("a"), blank, blank, item("b")}, `((ul (li "a")) (br) (ul (li "b")))`},
		{"list followed by trailing blanks", []gemtext.Line{item("a"), blank, blank}, `((ul (li "a")))`},

		// Unknown line type
		{"unknown line dropped", []gemtext.Line{text("a"), unknownLine{}, text("b")}, `((p "a") (p "b"))`},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if got := gmi2sxhtml.Convert(tc.lines).String(); got != tc.exp {
				t.Errorf("expected %q, but got %q", tc.exp, got)
			}
		})
	}
}
