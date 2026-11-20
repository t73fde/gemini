// -----------------------------------------------------------------------------
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
// -----------------------------------------------------------------------------

package gemtext_test

import (
	"reflect"
	"testing"

	"t73f.de/r/gemini/gemtext"
)

func TestParse(t *testing.T) {
	t.Parallel()
	type lines = []gemtext.Line
	tests := []struct {
		name string
		src  string
		want lines
	}{
		// Line endings and blank lines
		{"empty input", "", nil},
		{"single newline", "\n", lines{gemtext.Text{}}},
		{"crlf blank", "\r\n", lines{gemtext.Text{}}},
		{"no final newline", "a", lines{gemtext.Text{Text: "a"}}},
		{"final newline adds no line", "a\n", lines{gemtext.Text{Text: "a"}}},
		{"blank between", "a\n\nb", lines{gemtext.Text{Text: "a"}, gemtext.Text{}, gemtext.Text{Text: "b"}}},
		{"trailing blank", "a\n\n", lines{gemtext.Text{Text: "a"}, gemtext.Text{}}},
		{"crlf text", "a\r\nb\r\n", lines{gemtext.Text{Text: "a"}, gemtext.Text{Text: "b"}}},
		{"whitespace only", " \t \n", lines{gemtext.Text{Text: " \t "}}},
		{"leading whitespace kept", "  a", lines{gemtext.Text{Text: "  a"}}},

		// Links
		{"link", "=> gemini://x.org", lines{gemtext.Link{URL: "gemini://x.org"}}},
		{"link label", "=> /p Label text", lines{gemtext.Link{URL: "/p", Label: "Label text"}}},
		{"link no space", "=>/p L", lines{gemtext.Link{URL: "/p", Label: "L"}}},
		{"link tabs", "=>\t/p \t L \t", lines{gemtext.Link{URL: "/p", Label: "L \t"}}},
		{"link tabs first", "=>\t/p\t L \t", lines{gemtext.Link{URL: "/p", Label: "L \t"}}},
		{"link without url", "=>", lines{gemtext.Link{URL: "", Label: ""}}},
		{"link whitespace only", "=>  ", lines{gemtext.Link{URL: "", Label: ""}}},
		{"leading space no link", " => /p", lines{gemtext.Text{Text: " => /p"}}},
		{"link two spaces", "=> /p  L", lines{gemtext.Link{URL: "/p", Label: "L"}}},
		{"link no label", "=> /p  ", lines{gemtext.Link{URL: "/p", Label: ""}}},

		// Headings
		{"h1", "# a", lines{gemtext.Heading{Level: 1, Text: " a"}}},
		{"h2", "## a", lines{gemtext.Heading{Level: 2, Text: " a"}}},
		{"h3", "### a", lines{gemtext.Heading{Level: 3, Text: " a"}}},
		{"heading no space", "#a", lines{gemtext.Heading{Level: 1, Text: "a"}}},
		{"heading empty", "#", lines{gemtext.Heading{Level: 1}}},
		{"heading trimmed", "##  \t a b \t", lines{gemtext.Heading{Level: 2, Text: "  \t a b \t"}}},
		{"four hashes", "####x", lines{gemtext.Heading{Level: 3, Text: "#x"}}},
		{"four hashes space", "#### x", lines{gemtext.Heading{Level: 3, Text: "# x"}}},

		// List items
		{"list", "* a", lines{gemtext.ListItem{Text: "a"}}},
		{"list extra space", "*  a", lines{gemtext.ListItem{Text: " a"}}},
		{"list empty", "* ", lines{gemtext.ListItem{}}},
		{"star without space", "*a", lines{gemtext.Text{Text: "*a"}}},
		{"star alone", "*", lines{gemtext.Text{Text: "*"}}},

		// Quotes
		{"quote", "> a", lines{gemtext.Quote{Text: " a"}}},
		{"quote no space", ">a", lines{gemtext.Quote{Text: "a"}}},
		{"quote empty", ">", lines{gemtext.Quote{}}},

		// Preformatted
		{"pre", "```\na\n```", lines{gemtext.Preformatted{Lines: []string{"a"}}}},
		{"pre alt", "``` go \na\n```", lines{gemtext.Preformatted{AltText: " go ", Lines: []string{"a"}}}},
		{"pre verbatim", "```\n# h\n=> /p\n* l\n> q\n```", lines{gemtext.Preformatted{Lines: []string{"# h", "=> /p", "* l", "> q"}}}},
		{"pre blank and spaces", "```\n\n  a \n```", lines{gemtext.Preformatted{Lines: []string{"", "  a "}}}},
		{"pre crlf", "```\r\na\r\n```\r\n", lines{gemtext.Preformatted{Lines: []string{"a"}}}},
		{"pre empty", "```\n```", lines{gemtext.Preformatted{}}},
		{"pre unclosed", "```\nopen", lines{gemtext.Preformatted{Lines: []string{"open"}}}},
		{"pre unclosed empty", "```", lines{gemtext.Preformatted{}}},
		{"pre closing with text", "```\na\n```ignored\nb", lines{gemtext.Preformatted{Lines: []string{"a"}}, gemtext.Text{Text: "b"}}},
		{"pre indented fence", "```\n ```\n```", lines{gemtext.Preformatted{Lines: []string{" ```"}}}},
		{"indented fence is text", " ```", lines{gemtext.Text{Text: " ```"}}},
		{"two pre blocks", "```\na\n```\n```x\nb\n```", lines{
			gemtext.Preformatted{Lines: []string{"a"}},
			gemtext.Preformatted{AltText: "x", Lines: []string{"b"}},
		}},

		// Mixed document
		{"mixed", "# T\r\ntext\n\n=> /p L\n* i\n> q\n```go\nx\n```\n", lines{
			gemtext.Heading{Level: 1, Text: " T"},
			gemtext.Text{Text: "text"},
			gemtext.Text{},
			gemtext.Link{URL: "/p", Label: "L"},
			gemtext.ListItem{Text: "i"},
			gemtext.Quote{Text: " q"},
			gemtext.Preformatted{AltText: "go", Lines: []string{"x"}},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			doc := gemtext.Parse([]byte(tt.src))
			got := doc.Lines
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse(%q)\n got: %#v\nwant: %#v", tt.src, got, tt.want)
			}
		})
	}
}

func TestParseCopiesInput(t *testing.T) {
	t.Parallel()
	src := []byte("# head\n=> /p label\ntext\n```alt\npre\n```\n")
	doc := gemtext.Parse(src)
	for i := range src {
		src[i] = 'X'
	}
	want := []gemtext.Line{
		gemtext.Heading{Level: 1, Text: " head"},
		gemtext.Link{URL: "/p", Label: "label"},
		gemtext.Text{Text: "text"},
		gemtext.Preformatted{AltText: "alt", Lines: []string{"pre"}},
	}
	if !reflect.DeepEqual(doc.Lines, want) {
		t.Errorf("tree changed after modifying input: %#v", doc.Lines)
	}
}

func FuzzParse(f *testing.F) {
	for _, s := range []string{"", "\n", "# a", "=> /p x", "```\na", "* a\r\n> b", "####", "=>"} {
		f.Add([]byte(s))
	}
	f.Fuzz(func(t *testing.T, src []byte) {
		for _, line := range gemtext.Parse(src).Lines {
			if h, ok := line.(gemtext.Heading); ok && (h.Level < 1 || h.Level > 3) {
				t.Fatalf("invalid heading level %d for input %q", h.Level, src)
			}
		}
	})
}
