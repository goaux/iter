package transform_test

import (
	"maps"
	"slices"
	"testing"

	"github.com/goaux/iter/transform"
)

func TestZipIndex(t *testing.T) {
	t.Run("", func(t *testing.T) {
		zip := maps.Collect(transform.ZipIndex(slices.Values([]string{"a", "b", "c"})))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[0] != "a" {
			t.Error("must be a")
		}
		if zip[1] != "b" {
			t.Error("must be b")
		}
		if zip[2] != "c" {
			t.Error("must be c")
		}
	})

	t.Run("", func(t *testing.T) {
		zip := maps.Collect(transform.ZipIndex(slices.Values([]string{})))
		if len(zip) != 0 {
			t.Error("len(zip) must be 0")
		}
	})
}

func TestZip(t *testing.T) {
	t.Run("same length", func(t *testing.T) {
		zip := maps.Collect(transform.Zip(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[30] != 33 {
			t.Error("must be 33")
		}
	})

	t.Run("len(lhs) < len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.Zip(
			slices.Values([]int{10, 20}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 2 {
			t.Error("len(zip) must be 2")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
	})

	t.Run("len(lhs) > len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.Zip(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22}),
		))
		if len(zip) != 2 {
			t.Error("len(zip) must be 2")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
	})

	t.Run("zero length", func(t *testing.T) {
		zip := maps.Collect(transform.Zip(
			slices.Values([]int{10}),
			slices.Values([]int{}),
		))
		if len(zip) != 0 {
			t.Error("len(zip) must be 0")
		}
	})
}

func TestZipLeft(t *testing.T) {
	t.Run("same length", func(t *testing.T) {
		zip := maps.Collect(transform.ZipLeft(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[30] != 33 {
			t.Error("must be 33")
		}
	})

	t.Run("len(lhs) < len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.ZipLeft(
			slices.Values([]int{10, 20}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 2 {
			t.Error("len(zip) must be 2")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
	})

	t.Run("len(lhs) > len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.ZipLeft(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[30] != 0 {
			t.Error("must be 0")
		}
	})

	t.Run("zero length", func(t *testing.T) {
		zip := maps.Collect(transform.ZipLeft(
			slices.Values([]int{}),
			slices.Values([]int{}),
		))
		if len(zip) != 0 {
			t.Error("len(zip) must be 0")
		}
	})
}

func TestZipRight(t *testing.T) {
	t.Run("same length", func(t *testing.T) {
		zip := maps.Collect(transform.ZipRight(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[30] != 33 {
			t.Error("must be 33")
		}
	})

	t.Run("len(lhs) < len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.ZipRight(
			slices.Values([]int{10, 20}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[0] != 33 {
			t.Error("must be 22")
		}
	})

	t.Run("len(lhs) > len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.ZipRight(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22}),
		))
		if len(zip) != 2 {
			t.Error("len(zip) must be 2")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
	})

	t.Run("zero length", func(t *testing.T) {
		zip := maps.Collect(transform.ZipRight(
			slices.Values([]int{}),
			slices.Values([]int{}),
		))
		if len(zip) != 0 {
			t.Error("len(zip) must be 0")
		}
	})
}

func TestZipAll(t *testing.T) {
	t.Run("same length", func(t *testing.T) {
		zip := maps.Collect(transform.ZipAll(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[30] != 33 {
			t.Error("must be 33")
		}
	})

	t.Run("len(lhs) < len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.ZipAll(
			slices.Values([]int{10, 20}),
			slices.Values([]int{11, 22, 33}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[0] != 33 {
			t.Error("must be 33")
		}
	})

	t.Run("len(lhs) > len(rhs)", func(t *testing.T) {
		zip := maps.Collect(transform.ZipAll(
			slices.Values([]int{10, 20, 30}),
			slices.Values([]int{11, 22}),
		))
		if len(zip) != 3 {
			t.Error("len(zip) must be 3")
		}
		if zip[10] != 11 {
			t.Error("must be 11")
		}
		if zip[20] != 22 {
			t.Error("must be 22")
		}
		if zip[30] != 0 {
			t.Error("must be 0")
		}
	})

	t.Run("zero length", func(t *testing.T) {
		zip := maps.Collect(transform.ZipAll(
			slices.Values([]int{}),
			slices.Values([]int{}),
		))
		if len(zip) != 0 {
			t.Error("len(zip) must be 0")
		}
	})
}
