package bufioscanner_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing/iotest"

	"github.com/goaux/iter/bufioscanner"
)

func ExampleScanner_Bytes() {
	s := bufioscanner.NewScanner(bytes.NewBufferString(" Where are you going \n for your next vacation? \n"))
	s.Split(bufio.ScanWords)
	for i, word := range s.Bytes() {
		fmt.Printf("[%d] % 02x\n", i, word)
	}
	if err := s.Err(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// Output:
	// [0] 57 68 65 72 65
	// [1] 61 72 65
	// [2] 79 6f 75
	// [3] 67 6f 69 6e 67
	// [4] 66 6f 72
	// [5] 79 6f 75 72
	// [6] 6e 65 78 74
	// [7] 76 61 63 61 74 69 6f 6e 3f
}

func ExampleScanner_Bytes_break() {
	s := bufioscanner.NewScanner(bytes.NewBufferString(" Where are you going \n for your next vacation? \n"))
	s.Split(bufio.ScanWords)
	for i, word := range s.Bytes() {
		fmt.Printf("[%d] % 02x\n", i, word)
		if i == 3 {
			break
		}
	}
	if err := s.Err(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// Output:
	// [0] 57 68 65 72 65
	// [1] 61 72 65
	// [2] 79 6f 75
	// [3] 67 6f 69 6e 67
}

func ExampleScanner_Text() {
	s := bufioscanner.NewScanner(
		io.MultiReader(
			bytes.NewBufferString(" Where are you going \n for your next vacation? \n"),
			iotest.ErrReader(errors.New("io error?")),
		),
	)
	s.Split(bufio.ScanWords)
	for i, word := range s.Text() {
		fmt.Printf("[%d] %q\n", i, word)
	}
	if err := s.Err(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// Output:
	// [0] "Where"
	// [1] "are"
	// [2] "you"
	// [3] "going"
	// [4] "for"
	// [5] "your"
	// [6] "next"
	// [7] "vacation?"
	// error: io error?
}

func ExampleScanner_Text_break() {
	s := bufioscanner.NewScanner(bytes.NewBufferString(" Where are you going \n for your next vacation? \n"))
	s.Split(bufio.ScanWords)
	for i, word := range s.Text() {
		fmt.Printf("[%d] %q\n", i, word)
		if i == 3 {
			break
		}
	}
	if err := s.Err(); err != nil {
		fmt.Printf("error: %v\n", err)
	}
	// Output:
	// [0] "Where"
	// [1] "are"
	// [2] "you"
	// [3] "going"
}
