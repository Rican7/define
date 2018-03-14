// Package io TODO
//
// Copyright © 2018 Trevor N. Suarez (Rican7)
package io

import (
	"fmt"
	"io"
)

// PanicWriter is a writer that panics if a write operation causes an error
type PanicWriter struct {
	io.Writer
}

// Write writes a given string to the writer, and returns the number of bytes
// that were written. It'll panic if any error occurs during writing.
func (w *PanicWriter) Write(p []byte) int {
	n, err := w.Writer.Write(p)

	if nil != err {
		panic(err)
	}

	return n
}

// WriteString writes a given string to the writer, and returns the number of
// bytes that were written. It'll panic if any error occurs during writing.
func (w *PanicWriter) WriteString(p string) int {
	return w.Write([]byte(p))
}

// FWrite writes the given args to the writer like fmt.Sprint(), and returns
// the number of bytes that were written. It'll panic if any error occurs
// during writing.
func (w *PanicWriter) FWrite(p ...interface{}) int {
	return w.WriteString(fmt.Sprint(p))
}

// FWritef writes the given args to the writer like fmt.Sprintf(), and returns
// the number of bytes that were written. It'll panic if any error occurs
// during writing.
func (w *PanicWriter) FWritef(format string, p ...interface{}) int {
	return w.WriteString(fmt.Sprintf(format, p))
}

// FWriteln writes the given args to the writer like fmt.Sprintln(), and returns
// the number of bytes that were written. It'll panic if any error occurs
// during writing.
func (w *PanicWriter) FWriteln(p ...interface{}) int {
	return w.WriteString(fmt.Sprintln(p))
}

// WriteNewLine writes a new-line character to the writer, and returns the
// number of bytes that were written. It'll panic if any error occurs during
// writing.
func (w *PanicWriter) WriteNewLine() int {
	return w.WriteString("\n")
}

// WriteStringLine writes a given string to the writer with a new-line
// character after the given string, and returns the number of bytes that were
// written. It'll panic if any error occurs during writing.
func (w *PanicWriter) WriteStringLine(p string) int {
	return w.WriteString(p) + w.WriteNewLine()
}
