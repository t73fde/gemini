# sxhtml

Package `sxhtml` converts Gemtext (`text/gemini`) into SxHTML.

The input is the line-oriented AST produced by package
[`gemtext`](../gemtext). The output is a list of SxHTML block elements, built
from the data types of package [`sx`](https://t73f.de/r/sx). It can be
rendered to HTML by an SxHTML encoder.

## Usage

```go
import (
	"t73f.de/r/gemini/gemtext"
	"t73f.de/r/gemini/sxhtml" // adapt to the actual import path
)

lines := gemtext.Parse(data) // adapt to the actual parser function
blocks := sxhtml.Convert(lines)
fmt.Println(blocks)
```

## Mapping

| Gemtext line               | SxHTML                                          |
|----------------------------|-------------------------------------------------|
| Text                       | `(p "text")`                                    |
| Heading, level 1 to 3      | `(h1 "text")` to `(h3 "text")`                  |
| Link                       | `(p (a ((href . "URL")) "label"))`              |
| Quote                      | `(blockquote "text")`                           |
| Preformatted block         | `(pre (code "lines joined by LF"))`             |
| Consecutive list items     | `(ul (li "item") (li "item") ...)`              |

Details:

* A link without a label uses its URL as text. URLs are not validated;
  invalid URLs are left to the browser.
* Every quote line is converted on its own, independent of its neighbors,
  like text lines. Gemtext defines no multi-line quote construct.
* A list ends with the first line that is not a list item. Two lists
  separated by a blank line result in two `ul` elements.
* Headings with a level outside 1 to 3 are dropped.
* The alt text of a preformatted block is currently not output.

## Blank lines

Blank lines, i.e. lines that are empty or consist of white space only, do not
produce elements of their own, because all other elements are rendered as
paragraph-like blocks.

Instead, `n` consecutive blank lines between two elements, or at the start of
the input, produce `n-1` `(br)` elements. A single blank line therefore
produces nothing. Blank lines at the end of the input are dropped.
