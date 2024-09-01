// Package bufioreader provides a wrapper around [bufio.Reader] with additional iteration capabilities.
package bufioreader

import (
	"bufio"
	"errors"
	"io"
	"iter"
)

// Reader wraps a [bufio.Reader] and provides methods to iterate over its content.
type Reader struct {
	*bufio.Reader
	err error
}

// NewReader creates a new [Reader] wraps a [bufio.Reader].
func New(b *bufio.Reader) *Reader {
	return &Reader{Reader: b}
}

// NewReader creates a new [Reader] with a default buffer size.
func NewReader(rd io.Reader) *Reader {
	return New(bufio.NewReader(rd))
}

// NewReaderSize creates a new [Reader] with a buffer size.
func NewReaderSize(rd io.Reader, size int) *Reader {
	return New(bufio.NewReaderSize(rd, size))
}

// Err returns the last error. It returns nil if the last error is [io.EOF].
//
// [bufio.Reader.ReadBytes], [bufio.Reader.ReadSlice] and [bufio.Reader.ReadString]
// returns the data read before the error.
// This data would not be passed to loop body.
// Instead, it can be retrieved by [GetErrorBuffer] or [GetErrorBufferString].
//
// In normal, error-free situations, [bufio.Reader.ReadBytes],
// [bufio.Reader.ReadSlice], and [bufio.Reader.ReadString] will always return
// [io.EOF] on the final call. Since this is not an error that indicates an
// abnormality, Err() will return nil.
//
// In the loop body, one way to know that io.EOF has been returned is to check
// whether the delimiter is included in the result.
func (r *Reader) Err() error { return r.err }

// ReadBytes returns an iterator that yields byte slices delimited by the given byte using [bufio.Reader.ReadBytes].
func (r *Reader) ReadBytes(delim byte) iter.Seq2[int, []byte] {
	return read(r, func() ([]byte, error) { return r.Reader.ReadBytes(delim) })
}

// ReadSlice returns an iterator that yields byte slices delimited by the given byte using [bufio.Reader.ReadSlice].
func (r *Reader) ReadSlice(delim byte) iter.Seq2[int, []byte] {
	return read(r, func() ([]byte, error) { return r.Reader.ReadSlice(delim) })
}

// ReadString returns an iterator that yields strings delimited by the given byte using [bufio.Reader.ReadString].
func (r *Reader) ReadString(delim byte) iter.Seq2[int, string] {
	return read(r, func() (string, error) { return r.Reader.ReadString(delim) })
}

func read[T []byte | string](r *Reader, read func() (T, error)) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		for i := 0; ; i++ {
			v, err := read()
			if err != nil { // err != nil will stop the loop at the end
				if errors.Is(err, io.EOF) { // io.EOF is not an error
					if len(v) > 0 { // call the loop body with remaining data
						yield(i, v) // because the loop will end the result of yield is ignored
					}
				} else {
					if len(v) > 0 { // attach the remaining data with the err
						r.err = NewError(err, []byte(v))
					}
				}
				return // stops the loop
			}
			if !yield(i, v) {
				return // stops the loop
			}
		}
	}
}
