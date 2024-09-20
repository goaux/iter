package transform

import "iter"

// ZipIndex converts an [iter.Seq][T] to an [iter.Seq2][int, T].
// The int is the index of the element in the sequence.
// The order of the sequence is not changed.
func ZipIndex[T any](iterator iter.Seq[T]) iter.Seq2[int, T] {
	return func(yield func(int, T) bool) {
		i := 0
		for v := range iterator {
			if !yield(i, v) {
				break
			}
			i++
		}
	}
}

// Zip returns a composite iterator iter.Seq2[S, T] from two iterators iter.Seq[S] and iter.Seq[T].
// Zip iterates over the smaller of the two sequences.
func Zip[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T] {
	return func(yield func(S, T) bool) {
		nextL, stopL := iter.Pull(lhs)
		defer stopL()
		nextR, stopR := iter.Pull(rhs)
		defer stopR()
		for {
			l, okL := nextL()
			r, okR := nextR()
			if !okL || !okR {
				return
			}
			if !yield(l, r) {
				return
			}
		}
	}
}

// ZipLeft returns a composite iterator iter.Seq2[S, T] from two iterators iter.Seq[S] and iter.Seq[T].
// ZipLeft iterates as many times as there are elements on the left.
// If there are fewer elements on the right than the left, zero values are used for the right.
func ZipLeft[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T] {
	return func(yield func(S, T) bool) {
		nextL, stopL := iter.Pull(lhs)
		defer stopL()
		nextR, stopR := iter.Pull(rhs)
		defer stopR()
		for {
			l, okL := nextL()
			r, okR := nextR()
			if !okL {
				return
			}
			if !yield(l, r) {
				return
			}
			if !okR {
				for okL {
					if !yield(l, r) {
						return
					}
					l, okL = nextL()
				}
				return
			}
		}
	}
}

// ZipRight returns a composite iterator iter.Seq2[S, T] from two iterators iter.Seq[S] and iter.Seq[T].
// ZipRight iterates as many times as there are elements on the right.
// If there are fewer elements on the left than the right, zero values are used for the left.
func ZipRight[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T] {
	return func(yield func(S, T) bool) {
		nextL, stopL := iter.Pull(lhs)
		defer stopL()
		nextR, stopR := iter.Pull(rhs)
		defer stopR()
		for {
			l, okL := nextL()
			r, okR := nextR()
			if !okR {
				return
			}
			if !yield(l, r) {
				return
			}
			if !okL {
				for okR {
					if !yield(l, r) {
						return
					}
					r, okR = nextR()
				}
				return
			}
		}
	}
}

// ZipAll returns a composite iterator iter.Seq2[S, T] from two iterators iter.Seq[S] and iter.Seq[T].
// ZipAll iterates over the greater of the two sequences.
// Zero values are used for any missing elements in the sequence with fewer elements.
func ZipAll[S, T any](lhs iter.Seq[S], rhs iter.Seq[T]) iter.Seq2[S, T] {
	return func(yield func(S, T) bool) {
		nextL, stopL := iter.Pull(lhs)
		defer stopL()
		nextR, stopR := iter.Pull(rhs)
		defer stopR()
		for {
			l, okL := nextL()
			r, okR := nextR()
			if !okL {
				for okR {
					if !yield(l, r) {
						return
					}
					r, okR = nextR()
				}
				return
			}
			if !okR {
				for okL {
					if !yield(l, r) {
						return
					}
					l, okL = nextL()
				}
				return
			}
			if !yield(l, r) {
				return
			}
		}
	}
}
