package transform_test

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"testing"

	"github.com/goaux/iter/transform"
	"github.com/goaux/results"
)

func ExampleConcat() {
	got := slices.Collect(
		transform.Concat(
			slices.Values([]int{1, 2}),
			slices.Values([]int{3, 4, 5}),
			slices.Values([]int{6}),
		),
	)
	fmt.Println(got)
	// Output:
	// [1 2 3 4 5 6]
}

func ExampleConcat2() {
	got := maps.Collect(
		transform.Concat2(
			transform.Zip(
				slices.Values([]int{1, 2, 3}),
				slices.Values([]string{"a", "b", "c"}),
			),
			transform.Zip(
				slices.Values([]int{4, 5}),
				slices.Values([]string{"d", "e"}),
			),
			transform.Zip(
				slices.Values([]int{6}),
				slices.Values([]string{"f"}),
			),
		),
	)
	fmt.Println(got)
	// Output:
	// map[1:a 2:b 3:c 4:d 5:e 6:f]
}

func TestConcat(t *testing.T) {
	s := slices.Collect(transform.Concat(
		slices.Values([]int{1, 2, 3}),
		slices.Values([]int{4, 5}),
		slices.Values([]int{6}),
	))
	if !slices.Equal(s, []int{1, 2, 3, 4, 5, 6}) {
		t.Error("must be equal")
	}
}

func TestConcat2(t *testing.T) {
	m := maps.Collect(transform.Concat2(
		maps.All(map[int]string{1: "a", 2: "b", 3: "c"}),
		maps.All(map[int]string{4: "c", 5: "d"}),
		maps.All(map[int]string{6: "e"}),
	))
	want := map[int]string{
		1: "a", 2: "b", 3: "c",
		4: "c", 5: "d",
		6: "e",
	}
	if !maps.Equal(m, want) {
		t.Error("must be equal")
	}
}

func TestSkip(t *testing.T) {
	t.Run("0 < n < len", func(t *testing.T) {
		s := slices.Collect(transform.Skip(slices.Values([]int{11, 22, 33}), 2))
		if !slices.Equal(s, []int{33}) {
			t.Error("must be equal")
		}
	})

	t.Run("n = len", func(t *testing.T) {
		s := slices.Collect(transform.Skip(slices.Values([]int{11, 22, 33}), 3))
		if !slices.Equal(s, []int{}) {
			t.Error("must be equal")
		}
	})

	t.Run("n > len", func(t *testing.T) {
		s := slices.Collect(transform.Skip(slices.Values([]int{11, 22, 33}), 5))
		if !slices.Equal(s, []int{}) {
			t.Error("must be equal")
		}
	})

	t.Run("n = 0", func(t *testing.T) {
		s := slices.Collect(transform.Skip(slices.Values([]int{11, 22, 33}), 0))
		if !slices.Equal(s, []int{11, 22, 33}) {
			t.Error("must be equal")
		}
	})

	t.Run("n < 0", func(t *testing.T) {
		s := slices.Collect(transform.Skip(slices.Values([]int{11, 22, 33}), -2))
		if !slices.Equal(s, []int{0, 0, 11, 22, 33}) {
			t.Error("must be equal")
		}
	})
}

func TestSkip2(t *testing.T) {
	t.Run("0 < n < len", func(t *testing.T) {
		m := maps.Collect(transform.Skip2(slices.All([]int{11, 22, 33}), 2))
		want := map[int]int{2: 33}
		if !maps.Equal(m, want) {
			t.Error("must be equal")
		}
	})

	t.Run("n = len", func(t *testing.T) {
		s := maps.Collect(transform.Skip2(slices.All([]int{11, 22, 33}), 3))
		want := map[int]int{}
		if !maps.Equal(s, want) {
			t.Error("must be equal")
		}
	})

	t.Run("n > len", func(t *testing.T) {
		s := maps.Collect(transform.Skip2(slices.All([]int{11, 22, 33}), 5))
		want := map[int]int{}
		if !maps.Equal(s, want) {
			t.Error("must be equal")
		}
	})

	t.Run("n = 0", func(t *testing.T) {
		s := maps.Collect(transform.Skip2(slices.All([]int{11, 22, 33}), 0))
		want := map[int]int{0: 11, 1: 22, 2: 33}
		if !maps.Equal(s, want) {
			t.Error("must be equal")
		}
	})

	t.Run("n < 0", func(t *testing.T) {
		s := slices.Collect(
			transform.Keys(
				transform.Map2(
					transform.Skip2(
						transform.Zip(slices.Values([]int{3, 4, 5}), slices.Values([]int{33, 44, 55})),
						-2,
					),
					func(k, v int) (int, struct{}) {
						return k + v, struct{}{}
					},
				),
			),
		)
		want := []int{0, 0, 36, 48, 60}
		if !slices.Equal(s, want) {
			t.Error("must be equal")
		}
	})
}

func TestResize(t *testing.T) {
	t.Run("less", func(t *testing.T) {
		s := slices.Collect(transform.Resize(slices.Values([]int{11, 22, 33}), 2))
		if !slices.Equal(s, []int{11, 22}) {
			t.Error("must be equal")
		}
	})

	t.Run("more", func(t *testing.T) {
		s := slices.Collect(transform.Resize(slices.Values([]int{11, 22, 33}), 5))
		if !slices.Equal(s, []int{11, 22, 33, 0, 0}) {
			t.Error("must be equal")
		}
	})

	t.Run("equal", func(t *testing.T) {
		s := slices.Collect(transform.Resize(slices.Values([]int{11, 22, 33}), 3))
		if !slices.Equal(s, []int{11, 22, 33}) {
			t.Error("must be equal")
		}
	})
}

func TestResize2(t *testing.T) {
	t.Run("less", func(t *testing.T) {
		m := maps.Collect(
			transform.Resize2(
				transform.Zip(slices.Values([]int{1, 2, 3}), slices.Values([]int{2, 4, 6})),
				2,
			),
		)
		want := map[int]int{1: 2, 2: 4}
		if !maps.Equal(m, want) {
			t.Error("must be equal")
		}
	})

	t.Run("more", func(t *testing.T) {
		s := slices.Collect(
			transform.Keys(
				transform.Map2(
					transform.Resize2(
						transform.Zip(slices.Values([]int{3, 4, 5}), slices.Values([]int{33, 44, 55})),
						5,
					),
					func(k, v int) (int, struct{}) {
						return k + v, struct{}{}
					},
				),
			),
		)
		want := []int{36, 48, 60, 0, 0}
		if !slices.Equal(s, want) {
			t.Error("must be equal")
		}
	})

	t.Run("equal", func(t *testing.T) {
		m := maps.Collect(
			transform.Resize2(
				transform.Zip(slices.Values([]int{1, 2, 3}), slices.Values([]int{2, 4, 6})),
				3,
			),
		)
		want := map[int]int{1: 2, 2: 4, 3: 6}
		if !maps.Equal(m, want) {
			t.Error("must be equal")
		}
	})
}

func TestSwap(t *testing.T) {
	m := maps.Collect(transform.Swap(slices.All([]string{"a", "b", "c"})))
	if len(m) != 3 {
		t.Error("len(m) must be 3")
	}
	if m["a"] != 0 {
		t.Error("must be 0")
	}
	if m["b"] != 1 {
		t.Error("must be 0")
	}
	if m["c"] != 2 {
		t.Error("must be 0")
	}
}

func TestKeys(t *testing.T) {
	keys := slices.Sorted(transform.Keys(maps.All(map[string]int{"a": 1, "b": 2, "c": 3})))
	if len(keys) != 3 {
		t.Error("len(keys) must be 3")
	}
	if !slices.Equal(keys, []string{"a", "b", "c"}) {
		t.Error("must be equal")
	}
}

func TestValues(t *testing.T) {
	keys := slices.Sorted(transform.Values(maps.All(map[string]int{"a": 1, "b": 2, "c": 3})))
	if len(keys) != 3 {
		t.Error("len(keys) must be 3")
	}
	if !slices.Equal(keys, []int{1, 2, 3}) {
		t.Error("must be equal")
	}
}

func TestSelect(t *testing.T) {
	ss := slices.Collect(transform.Select(
		slices.Values([]int{0, 1, 2, 3, 4, 5}),
		func(v int) bool { return v%2 == 1 },
	))
	if len(ss) != 3 {
		t.Error("len(ss) must be 3")
	}
	if !slices.Equal(ss, []int{1, 3, 5}) {
		t.Error("must be equal")
	}
}

func TestSelect2(t *testing.T) {
	ss := maps.Collect(transform.Select2(
		maps.All(map[int]string{1: "a", 2: "b", 3: "c"}),
		func(i int, s string) bool { return i == 1 || s == "c" },
	))
	if len(ss) != 2 {
		t.Error("len(ss) must be 2")
	}
	if ss[1] != "a" {
		t.Error("must be a")
	}
	if ss[3] != "c" {
		t.Error("must be c")
	}
}

func TestMap(t *testing.T) {
	ss := slices.Collect(transform.Map(slices.Values([]int{1, 3, 5}), strconv.Itoa))
	if len(ss) != 3 {
		t.Error("len(ss) must be 2")
	}
	if !slices.Equal(ss, []string{"1", "3", "5"}) {
		t.Error("must be equal")
	}
}

func TestMap2(t *testing.T) {
	m := maps.Collect(transform.Map2(
		maps.All(map[int]string{1: "11", 2: "22", 3: "33"}),
		func(i int, s string) (string, int) {
			return strconv.Itoa(i), results.Must1(strconv.Atoi(s))
		},
	))
	if len(m) != 3 {
		t.Error("len(ss) must be 3")
	}
	if m["1"] != 11 {
		t.Error("must be 11")
	}
	if m["2"] != 22 {
		t.Error("must be 22")
	}
	if m["3"] != 33 {
		t.Error("must be 33")
	}
}

func TestMapIn(t *testing.T) {
	type Pair struct {
		I int
		V string
	}
	s := slices.Collect(
		transform.MapIn(
			slices.All([]string{"a", "b", "c"}),
			func(i int, v string) Pair { return Pair{I: i, V: v} },
		),
	)
	want := []Pair{{0, "a"}, {1, "b"}, {2, "c"}}
	if !slices.Equal(s, want) {
		t.Error("must be equal")
	}
}

func TestMapOut(t *testing.T) {
	type Pair struct {
		I int
		V string
	}
	m := maps.Collect(
		transform.MapOut(
			slices.Values([]Pair{{0, "a"}, {1, "b"}, {2, "c"}}),
			func(i Pair) (int, string) { return i.I, i.V },
		),
	)
	want := map[int]string{0: "a", 1: "b", 2: "c"}
	if !maps.Equal(m, want) {
		t.Error("must be equal")
	}
}

func TestSelectMap(t *testing.T) {
	ss := slices.Collect(transform.SelectMap(
		slices.Values([]int{1, 2, 3, 4, 5, 6}),
		func(i int) (string, bool) {
			if i%2 == 0 {
				return strconv.Itoa(i), true
			}
			return "", false
		},
	))
	if len(ss) != 3 {
		t.Error("len(ss) must be 3")
	}
	if !slices.Equal(ss, []string{"2", "4", "6"}) {
		t.Error("must be equal")
	}
}

func TestSelectMap2(t *testing.T) {
	m := maps.Collect(transform.SelectMap2(
		maps.All(map[int]string{1: "a", 2: "b", 3: "c"}),
		func(i int, s string) (string, rune, bool) {
			if i%2 == 1 {
				return strconv.Itoa(i), []rune(s)[0], true
			}
			return "", 0, false
		},
	))
	if len(m) != 2 {
		t.Error("len(ss) must be 2")
	}
	if m["1"] != 'a' {
		t.Error("must be a")
	}
	if m["3"] != 'c' {
		t.Error("must be c")
	}
}

func TestSelectMapIn(t *testing.T) {
	type Pair struct {
		I int
		V string
	}
	s := slices.Collect(
		transform.SelectMapIn(
			slices.All([]string{"a", "b", "c"}),
			func(i int, v string) (Pair, bool) { return Pair{I: i, V: v}, i%2 == 0 },
		),
	)
	want := []Pair{{0, "a"}, {2, "c"}}
	if !slices.Equal(s, want) {
		t.Error("must be equal")
	}

}

func TestSelectMapOut(t *testing.T) {
	type Pair struct {
		I int
		V string
	}
	m := maps.Collect(
		transform.SelectMapOut(
			slices.Values([]Pair{{0, "a"}, {1, "b"}, {2, "c"}}),
			func(i Pair) (int, string, bool) { return i.I, i.V, i.I%2 == 0 },
		),
	)
	want := map[int]string{0: "a", 2: "c"}
	if !maps.Equal(m, want) {
		t.Error("must be equal")
	}
}
