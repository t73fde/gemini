# Gemini - Gemini protocol components for Go applications

This is a collection of packages for building software around the
[Gemini protocol](https://geminiprotocol.net/) in [Go](https://go.dev/).

* [gemtext](/dir?ci=tip&name=gemtext): A parser for text/gemini (Gemtext)
  documents, producing a simple, line-oriented syntax tree.

The name _Gemini_ refers to the Gemini protocol, a lightweight protocol for
serving and retrieving text documents. Its document format is called _Gemtext_
and has the media type `text/gemini`.

## Usage instructions

To import this library into your own [Go](https://go.dev/) software, you need
to run the `go get` command. Since Go does not handle non-standard software and
platforms well, some additional steps are required.

First, install the version control system [Fossil](https://fossil-scm.org),
which is a superior alternative to Git in many use cases. Fossil is just a
single executable, nothing more. Make sure it is included in your system's
command search path.

Then, run the following Go command to retrieve a specific version of
this library:

    GOVCS=t73f.de:fossil go get t73f.de/r/gemini@HASH

Here, `HASH` represents the commit hash of the version you want to use.

Go currently does not seem to support software versioning for projects managed
by Fossil. This is why the hash value is required. However, this method works
reliably.

Individual packages are imported by their path below the library, for example
`t73f.de/r/gemini/gemtext`.