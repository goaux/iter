// Package bufioscanner provides a wrapper around [bufio.Scanner] with iteration capabilities.
package bufioscanner

import (
	"bufio"
	"io"
	"iter"
)

// Scanner wraps a [bufio.Scanner] and provides methods to iterate over its content.
type Scanner struct {
	*bufio.Scanner
}

// New creates a new [Scanner] from a [bufio.Scanner].
func New(s *bufio.Scanner) *Scanner {
	return &Scanner{Scanner: s}
}

// NewScanner creates a new [Scanner] from an [io.Reader].
func NewScanner(r io.Reader) *Scanner {
	return New(bufio.NewScanner(r))
}

// Bytes returns an iterator that yields byte slices from the scanner using [bufio.Scanner.Bytes].
func (s *Scanner) Bytes() iter.Seq2[int, []byte] {
	return func(yield func(int, []byte) bool) {
		for i := 0; s.Scan(); i++ {
			if !yield(i, s.Scanner.Bytes()) {
				return
			}
		}
	}
}

// Text returns an iterator that yields strings from the scanner using [bufio.Scanner.Text].
func (s *Scanner) Text() iter.Seq2[int, string] {
	return func(yield func(int, string) bool) {
		for i := 0; s.Scan(); i++ {
			if !yield(i, s.Scanner.Text()) {
				return
			}
		}
	}
}
