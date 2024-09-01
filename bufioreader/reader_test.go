package bufioreader_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/goaux/iter/bufioreader"
)

func Example() {
	r := bufioreader.NewReader(io.MultiReader(
		strings.NewReader("hello\nworld\nexample"),
		iotest.ErrReader(errors.New("io error?")), // simulate some error
	))
	for i, s := range r.ReadString('\n') { // ReadString returns an iterator
		fmt.Printf("[%d] %q\n", i, s)
	}
	if err := r.Err(); err != nil {
		fmt.Printf("error: %v (remain:%q)\n", err, bufioreader.GetErrorBufferString(err))
	}
	// Output:
	// [0] "hello\n"
	// [1] "world\n"
	// error: io error? (remain:"example")
}

func ExampleReader_ReadBytes() {
	r := bufioreader.NewReader(strings.NewReader("hello\nworld"))
	for i, b := range r.ReadBytes('\n') {
		fmt.Printf("[%d] % 02x\n", i, b)
	}
	if err := r.Err(); err != nil { // r.Err never returns io.EOF
		fmt.Printf("error: %v (remain:%q)\n", err, bufioreader.GetErrorBufferString(err))
	}
	// Output:
	// [0] 68 65 6c 6c 6f 0a
	// [1] 77 6f 72 6c 64
}

func ExampleReader_ReadSlice() {
	r := bufioreader.NewReader(strings.NewReader("hello\nworld"))
	for i, b := range r.ReadSlice('\n') {
		fmt.Printf("[%d] % 02x\n", i, b)
	}
	if err := r.Err(); err != nil { // r.Err never returns io.EOF
		fmt.Printf("error: %v (remain:%q)\n", err, bufioreader.GetErrorBufferString(err))
	}
	// Output:
	// [0] 68 65 6c 6c 6f 0a
	// [1] 77 6f 72 6c 64
}

func ExampleReader_ReadString() {
	r := bufioreader.NewReaderSize(strings.NewReader("hello\nworld"), 32)
	for i, s := range r.ReadString('\n') {
		fmt.Printf("[%d] %q\n", i, s)
	}
	if err := r.Err(); err != nil { // r.Err never returns io.EOF
		fmt.Printf("error: %v (remain:%q)\n", err, bufioreader.GetErrorBufferString(err))
	}
	// Output:
	// [0] "hello\n"
	// [1] "world"
}

func TestReader_Bytes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		delim    byte
		expected [][]byte
	}{
		{
			name:     "simple case",
			input:    "hello\nworld",
			delim:    '\n',
			expected: [][]byte{[]byte("hello\n"), []byte("world")},
		},
		{
			name:     "empty input",
			input:    "",
			delim:    '\n',
			expected: nil,
		},
		{
			name:     "no delimiter",
			input:    "helloworld",
			delim:    '\n',
			expected: [][]byte{[]byte("helloworld")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufioreader.NewReader(strings.NewReader(tt.input))
			var result [][]byte
			for _, s := range r.ReadBytes(tt.delim) {
				result = append(result, s)
			}
			if r.Err() != nil {
				t.Errorf("unexpected error: %v", r.Err())
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}

	t.Run("an error", func(t *testing.T) {
		r := bufioreader.NewReader(io.MultiReader(
			strings.NewReader("hello\nworld\nexample"),
			iotest.ErrReader(errors.New("an error")),
		))
		var result [][]byte
		for _, s := range r.ReadBytes('\n') {
			result = append(result, s)
		}
		if r.Err() == nil {
			t.Error("expected error")
		}
		if bufioreader.GetErrorBufferString(r.Err()) != "example" {
			t.Error("expected remaining data")
		}
		expected := [][]byte{[]byte("hello\n"), []byte("world\n")}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestReader_Slice(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		delim    byte
		expected [][]byte
	}{
		{
			name:     "simple case",
			input:    "hello\nworld",
			delim:    '\n',
			expected: [][]byte{[]byte("hello\n"), []byte("world")},
		},
		{
			name:     "empty input",
			input:    "",
			delim:    '\n',
			expected: nil,
		},
		{
			name:     "no delimiter",
			input:    "helloworld",
			delim:    '\n',
			expected: [][]byte{[]byte("helloworld")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufioreader.NewReader(strings.NewReader(tt.input))
			var result [][]byte
			for _, s := range r.ReadSlice(tt.delim) {
				result = append(result, bytes.Clone(s))
			}
			if r.Err() != nil {
				t.Errorf("unexpected error: %v", r.Err())
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestReader_String(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		delim    byte
		expected []string
	}{
		{
			name:     "simple case",
			input:    "hello\nworld",
			delim:    '\n',
			expected: []string{"hello\n", "world"},
		},
		{
			name:     "empty input",
			input:    "",
			delim:    '\n',
			expected: nil,
		},
		{
			name:     "no delimiter",
			input:    "helloworld",
			delim:    '\n',
			expected: []string{"helloworld"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := bufioreader.NewReader(strings.NewReader(tt.input))
			var result []string
			for _, s := range r.ReadString(tt.delim) {
				result = append(result, s)
			}
			if r.Err() != nil {
				t.Errorf("unexpected error: %v", r.Err())
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}

	t.Run("break", func(t *testing.T) {
		r := bufioreader.NewReader(strings.NewReader("hello\nworld\nhello\nworld\nhello\nworld"))
		var result []string
		for i, s := range r.ReadString('\n') {
			result = append(result, s)
			if i == 3 {
				break
			}
		}
		if r.Err() != nil {
			t.Errorf("unexpected error: %v", r.Err())
		}
		expected := []string{"hello\n", "world\n", "hello\n", "world\n"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}
