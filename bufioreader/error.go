package bufioreader

import "errors"

// Error is a custom error type that includes a buffer of bytes along with the error.
type Error struct {
	// Err holds the underlying error.
	Err error

	// Buf holds the buffer.
	Buf []byte
}

// NewError creates a new Error instance. If the buffer is empty, it returns the original error.
// Otherwise, it returns a new Error with the provided error and buffer.
func NewError(err error, buf []byte) error {
	return &Error{Err: err, Buf: buf}
}

// Error returns the error message of the underlying error.
func (err *Error) Error() string { return err.Err.Error() }

// Unwrap returns the underlying error.
func (err *Error) Unwrap() error { return err.Err }

// Buffer returns the byte buffer associated with the error.
func (err *Error) Buffer() []byte { return err.Buf }

// BufferString returns the byte buffer as a string.
func (err *Error) BufferString() string { return string(err.Buf) }

// GetErrorBuffer attempts to extract the byte buffer from an error.
// If the error is of type *[Error], it returns the buffer calling [Error.Buffer].
// Otherwise, it returns nil.
//
// see [Reader.Err].
func GetErrorBuffer(err error) []byte {
	var buf *Error
	if errors.As(err, &buf) {
		return buf.Buffer()
	}
	return nil
}

// GetErrorBufferString attempts to extract the byte buffer from an error and return it as a string.
// If the error is of type *[Error], it returns the buffer as a string calling [Error.BufferString].
// Otherwise, it returns an empty string.
//
// see [Reader.Err].
func GetErrorBufferString(err error) string {
	var buf *Error
	if errors.As(err, &buf) {
		return buf.BufferString()
	}
	return ""
}
