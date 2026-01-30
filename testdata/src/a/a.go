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
