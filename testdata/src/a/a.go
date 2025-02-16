package a

type A[T any] = any

func f() {
	var _ any = A[any](nil)        // OK
	var _ A[any] = 100             // OK
	var _ A[any] = any(nil)        // OK
	var _ A[string] = A[bool](nil) // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
}
