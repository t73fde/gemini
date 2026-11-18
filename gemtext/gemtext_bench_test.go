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
	"fmt"
	"strings"
	"testing"

	"t73f.de/r/gemini/gemtext"
)

const mixedUnit = "# Heading\n" +
	"Some plain text line with a few words in it.\n" +
	"\n" +
	"=> gemini://example.org/page Link label\n" +
	"* item one\n" +
	"* item two\n" +
	"> quoted text\n" +
	"```go\n" +
	"func main() {}\n" +
	"\treturn\n" +
	"```\n"

func BenchmarkParse(b *testing.B) {
	profiles := []struct {
		name string
		unit string
	}{
		{"mixed", mixedUnit},
		{"text", "Some plain text line with a few words in it.\n"},
		{"links", "=> gemini://example.org/page Link label\n"},
		{"blank", "\n"},
		{"pre", "```\nverbatim line one\nverbatim line two\n```\n"},
	}
	for _, p := range profiles {
		for _, n := range []int{1, 100, 10_000} {
			src := []byte(strings.Repeat(p.unit, n))
			b.Run(fmt.Sprintf("%s/%dB", p.name, len(src)), func(b *testing.B) {
				b.ReportAllocs()
				b.SetBytes(int64(len(src)))
				for b.Loop() {
					gemtext.Parse(src)
				}
			})
		}
	}
}
