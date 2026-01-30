package a

type A[T any] = any

func p[T any](_ ...A[T]) {}

type B interface {
	M(A[A[string]])
}

type b struct{}

func (b) M(_ A[A[string]]) {}

type invalidB struct{}

func (invalidB) M(_ A[A[bool]]) {}

type C struct {
	caller B
}

type D struct {
	field A[string]
}

// Option is a generic named type for testing phantom type checking through generic wrappers.
type Option[T any] struct {
	value *T
}

type Writer interface {
	Write(v Option[A[string]]) error
}

type goodWriter struct{}

func (goodWriter) Write(_ Option[A[string]]) error { return nil }

type badWriter struct{}

func (badWriter) Write(_ Option[A[bool]]) error { return nil }
