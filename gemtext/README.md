# gemtext

Package `gemtext` parses
[text/gemini](https://geminiprotocol.net/docs/specification.gmi) (Gemtext)
documents into a simple, line-oriented abstract syntax tree (AST).

The AST is independent of any output format. It is meant as an input for
converters, for example into other document encodings or HTML.


## AST

`Parse` returns a `Document`, which holds a flat list of lines. Each line is
one of the following types:

| Type           | Source line                                         |
|----------------|-----------------------------------------------------|
| `Text`         | plain text line, possibly empty                     |
| `Link`         | `=> URL [label]`                                    |
| `Heading`      | `#`, `##` or `###` followed by text                 |
| `ListItem`     | `* text`                                            |
| `Quote`        | `>text`                                             |
| `Preformatted` | lines between two \`\`\` lines, with optional alt text |

The tree is deliberately flat. Grouping of consecutive lines, for example into
paragraphs, lists or quote blocks, is left to the consumer, because different
output formats need different groupings.

All strings in the tree are copies. The tree does not refer to the input.

## Parsing rules

Parsing never fails. Input that does not match a more specific line type
becomes a `Text` line. The following decisions go beyond the plain wording of
the specification:

- Lines end with LF or CRLF. A final line terminator does not start an
  additional empty line.
- A line that starts with `=>` but has no URL is a `Text` line.
- At most three leading `#` characters count as the heading level. Further `#`
  characters belong to the heading text, so `####x` is a level 3 heading with
  the text `#x`.
- A `*` is a list item marker only if it is followed by a space.
- Text after the closing `` ``` `` of a preformatted block is ignored.
- A preformatted block that is not closed at the end of the input is closed
  implicitly.
- Whitespace means space and tab.

## Limitations

- URLs are neither validated nor normalized nor resolved against a base URL.
- The input is not checked for valid UTF-8. A leading byte order mark is not
  removed.
- Values of type `Line` must not be compared with `==`, because `Preformatted`
  contains a slice. Use a type switch or `reflect.DeepEqual`.
