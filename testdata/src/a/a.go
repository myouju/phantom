package a

type A[T any] = any

func f() {
	var _ any = A[any](nil)                                             // OK
	var _ A[any] = 100                                                  // OK
	var _ A[any] = any(nil)                                             // OK
	var _ A[string] = A[bool](nil)                                      // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	func(A[string]) {}(A[bool](nil))                                    // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]) = A[bool](nil)                                      // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]), *new(A[bool]) = func() (_, _ A[bool]) { return }() // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]), *new(bool) = map[int]A[bool]{}[0]                  // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]), *new(bool) = any(nil).(A[bool])                    // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _, _ A[string] = func() (_ A[string], _ A[bool]) { return }()   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
}
