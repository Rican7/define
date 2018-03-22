// Copyright © 2018 Trevor N. Suarez (Rican7)

package io

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

// Enforce interface contracts
var (
	_ io.Writer = (*PanicWriter)(nil)
)

type writerShouldError bool

func (w writerShouldError) Write(p []byte) (int, error) {
	if bool(w) {
		return 0, fmt.Errorf("error during Write of bytes: %+v", p)
	}

	return len(p), nil
}

func TestNewPanicWriter(t *testing.T) {
	pw := NewPanicWriter(&strings.Builder{})

	if nil == pw {
		t.Errorf("NewPanicWriter returned nil")
	}
}

func TestWrite(t *testing.T) {
	toWrite := []byte("test")
	want := len(toWrite)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	got := pw.WriteBytes(toWrite)

	if got != want || got != w.Len() || w.String() != string(toWrite) {
		t.Errorf(
			"Write didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}
}

func TestWriteWithSpaces(t *testing.T) {
	toWrite := []byte("test")
	numSpaces := 5
	expectedBytes := append(bytes.Repeat([]byte(" "), numSpaces), toWrite...)
	want := len(expectedBytes)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w, spaces: uint(numSpaces)}

	got := pw.WriteBytes(toWrite)

	if got != want || got != w.Len() {
		t.Errorf(
			"Write didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}

	if w.String() != string(expectedBytes) {
		t.Errorf(
			"Writer didn't write the expected bytes. Got %+v. Want %+v.",
			[]byte(w.String()),
			expectedBytes,
		)
	}
}

func TestWritePanicsOnError(t *testing.T) {
	defer func() {
		if nil == recover() {
			t.Errorf("Write with an error did not panic.")
		}
	}()

	pw := &PanicWriter{inner: writerShouldError(true)}

	pw.WriteBytes([]byte(""))
}

func TestWriteString(t *testing.T) {
	toWrite := "test"
	want := len(toWrite)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	got := pw.WriteString(toWrite)

	if got != want || got != w.Len() {
		t.Errorf(
			"WriteString didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}

	if w.String() != toWrite {
		t.Errorf(
			"Writer didn't write the expected string. Got %q. Want %q.",
			w.String(),
			toWrite,
		)
	}
}

func TestPrint(t *testing.T) {
	toWrite := []interface{}{1, true}
	expectedString := fmt.Sprint(toWrite...)
	want := len(expectedString)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	got := pw.Print(toWrite...)

	if got != want || got != w.Len() {
		t.Errorf(
			"Print didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}

	if w.String() != expectedString {
		t.Errorf(
			"Writer didn't write the expected string. Got %q. Want %q.",
			w.String(),
			expectedString,
		)
	}
}

func TestPrintf(t *testing.T) {
	format := "%d %v"
	toWrite := []interface{}{1, true}
	expectedString := fmt.Sprintf(format, toWrite...)
	want := len(expectedString)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	got := pw.Printf(format, toWrite...)

	if got != want || got != w.Len() {
		t.Errorf(
			"Printf didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}

	if w.String() != expectedString {
		t.Errorf(
			"Writer didn't write the expected string. Got %q. Want %q.",
			w.String(),
			expectedString,
		)
	}
}

func TestPrintln(t *testing.T) {
	toWrite := []interface{}{1, true}
	expectedString := fmt.Sprintln(toWrite...)
	want := len(expectedString)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	got := pw.Println(toWrite...)

	if got != want || got != w.Len() {
		t.Errorf(
			"Println didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}

	if w.String() != expectedString {
		t.Errorf(
			"Writer didn't write the expected string. Got %q. Want %q.",
			w.String(),
			expectedString,
		)
	}
}

func TestWriteNewLine(t *testing.T) {
	expectedString := "\n"
	want := len(expectedString)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	got := pw.WriteNewLine()

	if got != want || got != w.Len() {
		t.Errorf(
			"WriteNewLine didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}

	if w.String() != expectedString {
		t.Errorf(
			"Writer didn't write the expected string. Got %q. Want %q.",
			w.String(),
			expectedString,
		)
	}
}

func TestWriteStringLine(t *testing.T) {
	toWrite := "test"
	expectedString := toWrite + "\n"
	want := len(expectedString)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	got := pw.WriteStringLine(toWrite)

	if got != want || got != w.Len() {
		t.Errorf(
			"WriteStringLine didn't write the expected number of bytes. Got %d. Want %d.",
			got,
			want,
		)
	}

	if w.String() != expectedString {
		t.Errorf(
			"Writer didn't write the expected string. Got %q. Want %q.",
			w.String(),
			expectedString,
		)
	}
}

func TestIndentWrites(t *testing.T) {
	indentSize := uint(2)

	w := &strings.Builder{}
	pw := &PanicWriter{inner: w}

	if 0 != pw.spaces {
		t.Errorf(
			"Writer has incorrect indent size. Got %d. Want %d.",
			pw.spaces,
			0,
		)
	}

	pw.IndentWrites(indentSize, func(pw *PanicWriter) {
		if indentSize != pw.spaces {
			t.Errorf(
				"Writer has incorrect indent size. Got %d. Want %d.",
				pw.spaces,
				indentSize,
			)
		}

		// Test multi-level (nested) indenting
		pw.IndentWrites(indentSize, func(pw *PanicWriter) {
			if indentSize+indentSize != pw.spaces {
				t.Errorf(
					"Writer has incorrect indent size. Got %d. Want %d.",
					pw.spaces,
					indentSize,
				)
			}
		})
	})
}

func TestIndented(t *testing.T) {
	indentSize := uint(2)

	pw := &PanicWriter{inner: &strings.Builder{}}

	if 0 != pw.spaces {
		t.Errorf(
			"Writer has incorrect indent size. Got %d. Want %d.",
			pw.spaces,
			0,
		)
	}

	w := pw.indented(indentSize)

	if indentSize != w.spaces {
		t.Errorf(
			"Writer has incorrect indent size. Got %d. Want %d.",
			w.spaces,
			indentSize,
		)
	}
}
