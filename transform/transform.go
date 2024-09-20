// Package transform provides functions transforming iterators.
package transform

import "iter"

// Concat returns a single iterator concatenating the passed in iterators.
func Concat[T any](iterators ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, iterator := range iterators {
			for v := range iterator {
				if !yield(v) {
					return
				}
			}
		}
	}
}

// Concat2 returns a single iterator concatenating the passed in iterators.
func Concat2[S, T any](iterators ...iter.Seq2[S, T]) iter.Seq2[S, T] {
	return func(yield func(S, T) bool) {
		for _, iterator := range iterators {
			for s, t := range iterator {
				if !yield(s, t) {
					return
				}
			}
		}
	}
}

// Skip returns an iterator that skips the first skip element in the sequence.
// If skip is a negative number, it returns an iterator that adds skip's
// absolute value number of zero values to the beginning of the sequence.
func Skip[T any](iterator iter.Seq[T], skip int) iter.Seq[T] {
	switch {
	case skip > 0:
		return func(yield func(T) bool) {
			next, stop := iter.Pull(iterator)
			defer stop()
			for range skip {
				if _, ok := next(); !ok {
					return
				}
			}
			for {
				if v, ok := next(); !ok || !yield(v) {
					break
				}
			}
		}

	case skip < 0:
		return func(yield func(T) bool) {
			var zero T
			for range -skip {
				if !yield(zero) {
					return
				}
			}
			next, stop := iter.Pull(iterator)
			defer stop()
			for {
				if v, ok := next(); !ok || !yield(v) {
					break
				}
			}
		}

	default:
		return iterator
	}
}

// Skip2 returns an iterator that skips the first skip element in the sequence.
// If skip is a negative number, it returns an iterator that adds skip's
// absolute value number of zero values to the beginning of the sequence.
func Skip2[S, T any](iterator iter.Seq2[S, T], skip int) iter.Seq2[S, T] {
	switch {
	case skip > 0:
		return func(yield func(S, T) bool) {
			next, stop := iter.Pull2(iterator)
			defer stop()
			for range skip {
				if _, _, ok := next(); !ok {
					return
				}
			}
			for {
				if s, t, ok := next(); !ok || !yield(s, t) {
					break
				}
			}
		}

	case skip < 0:
		return func(yield func(S, T) bool) {
			var zeroS S
			var zeroT T
			for range -skip {
				if !yield(zeroS, zeroT) {
					return
				}
			}
			next, stop := iter.Pull2(iterator)
			defer stop()
			for {
				if s, t, ok := next(); !ok || !yield(s, t) {
					break
				}
			}
		}

	default:
		return iterator
	}
}

// Resize returns an iterator that changes the number of elements in a
// sequence. If you specify a size larger than the number of elements in the
// sequence, the iterator returns zero values for the missing elements.
// A negative size is considered to be 0.
func Resize[T any](iterator iter.Seq[T], size int) iter.Seq[T] {
	return func(yield func(T) bool) {
		next, stop := iter.Pull(iterator)
		defer stop()
		for i := 0; i < size; i++ {
			v, ok := next()
			if !ok {
				var zero T
				for ; i < size && yield(zero); i++ {
					// empty
				}
				break
			}
			if !yield(v) {
				break
			}
		}
	}
}

// Resize2 returns an iterator that changes the number of elements in a
// sequence. If you specify a size larger than the number of elements in the
// sequence, the iterator returns zero values for the missing elements.
// A negative size is considered to be 0.
func Resize2[S, T any](iterator iter.Seq2[S, T], size int) iter.Seq2[S, T] {
	return func(yield func(S, T) bool) {
		next, stop := iter.Pull2(iterator)
		defer stop()
		for i := 0; i < size; i++ {
			s, t, ok := next()
			if !ok {
				var zeroS S
				var zeroT T
				for ; i < size && yield(zeroS, zeroT); i++ {
					// empty
				}
				break
			}
			if !yield(s, t) {
				break
			}
		}
	}
}

// Swap converts an [iter.Seq2][S, T] to an [iter.Seq2][T, S].
// The order of the sequence is not changed.
func Swap[S, T any](iterator iter.Seq2[S, T]) iter.Seq2[T, S] {
	return func(yield func(T, S) bool) {
		for s, t := range iterator {
			if !yield(t, s) {
				break
			}
		}
	}
}

// Keys converts [iter.Seq2][K, V] to [iter.Seq][K].
// The order of the sequence is not changed.
func Keys[K, V any](iterator iter.Seq2[K, V]) iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range iterator {
			if !yield(k) {
				break
			}
		}
	}
}

// Values converts [iter.Seq2][K, V] to [iter.Seq][V].
// The order of the sequence is not changed.
func Values[K, V any](iterator iter.Seq2[K, V]) iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range iterator {
			if !yield(v) {
				break
			}
		}
	}
}

// Select returns an iterator over sequences of individual values with a condition.
// When called as seq(yield), seq calls f(v), then yield(v) if f returns true
// for each value v in the sequence, stopping early if yield returns false.
func Select[T any](iterator iter.Seq[T], f func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for it := range iterator {
			if f(it) && !yield(it) {
				break
			}
		}
	}
}

// Select2 returns an iterator over sequences of pairs of values with a condition, most
// commonly key-value pairs. When called as seq(yield), seq calls f(k, v), then
// if f returns true calls yield(k, v) for each pair (k, v) in the sequence,
// stopping early if yield returns false.
func Select2[K, V any](iterator iter.Seq2[K, V], f func(K, V) bool) iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for k, v := range iterator {
			if f(k, v) && !yield(k, v) {
				break
			}
		}
	}
}

// Map transforms [iter.Seq][S] to [iter.Seq][T] using f, which transforms S to T.
func Map[S, T any](iterator iter.Seq[S], f func(S) T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range iterator {
			if !yield(f(v)) {
				break
			}
		}
	}
}

// Map2 transforms [iter.Seq2][S, T] to [iter.Seq2][U, V] using f, which transforms S and T to U and V.
func Map2[S, T, U, V any](iterator iter.Seq2[S, T], f func(S, T) (U, V)) iter.Seq2[U, V] {
	return func(yield func(U, V) bool) {
		for s, t := range iterator {
			if !yield(f(s, t)) {
				break
			}
		}
	}
}

// MapIn transforms [iter.Seq2][S, T] to [iter.Seq][U] using f, which transforms S and T to U.
func MapIn[S, T, U any](iterator iter.Seq2[S, T], f func(S, T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for s, t := range iterator {
			if !yield(f(s, t)) {
				break
			}
		}
	}
}

// MapOut transforms [iter.Seq][S] to [iter.Seq2][T, U] using f, which transforms S to T and U.
func MapOut[S, T, U any](iterator iter.Seq[S], f func(S) (T, U)) iter.Seq2[T, U] {
	return func(yield func(T, U) bool) {
		for s := range iterator {
			if !yield(f(s)) {
				break
			}
		}
	}
}
