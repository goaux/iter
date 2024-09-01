package bufioreader_test

import (
	"errors"
	"io"
	"reflect"
	"testing"

	"github.com/goaux/iter/bufioreader"
)

func TestGetErrorBuffer(t *testing.T) {
	got := bufioreader.GetErrorBuffer(io.EOF)
	if got != nil {
		t.Errorf("expected nil")
	}
}

func TestGetErrorBufferString(t *testing.T) {
	got := bufioreader.GetErrorBufferString(io.EOF)
	if got != "" {
		t.Errorf(`expected ""`)
	}
}

func TestError_Buffer(t *testing.T) {
	src := []byte{0, 1, 2, 3, 4, 5}
	err := bufioreader.NewError(nil, src)
	got := bufioreader.GetErrorBuffer(err)
	if !reflect.DeepEqual(got, src) {
		t.Errorf("expected %v, got %v", src, got)
	}
}

func TestError_Unwrap(t *testing.T) {
	err := bufioreader.NewError(io.EOF, nil)
	err = errors.Unwrap(err)
	if err != io.EOF {
		t.Error("err must be io.EOF")
	}
}
