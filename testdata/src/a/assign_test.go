package a

import "fmt"

func assignTests() {
	var _ any = A[any](nil)                                             // OK
	var _ A[any] = 100                                                  // OK
	var _ A[any] = any(nil)                                             // OK
	var _ A[string] = A[bool](nil)                                      // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]) = A[bool](nil)                                      // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]), *new(A[bool]) = func() (_, _ A[bool]) { return }() // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]), *new(bool) = map[int]A[bool]{}[0]                  // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	*new(A[string]), *new(bool) = any(nil).(A[bool])                    // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _, _ A[string] = func() (_ A[string], _ A[bool]) { return }()   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`

	// Multiple assignment from function returning phantom types
	var _, _ = func() (A[string], A[bool]) { return nil, nil }() // OK (no explicit type on LHS)

	// Range loop variable assignment
	for _, v := range []A[string]{A[string](nil)} {
		var _ A[bool] = v // want `type annotations are not assignable: a\.A\[string\] to a\.A\[bool\]`
		_ = v             // OK
	}

	fmt.Printf("I'm %s", "Gopher") // OK
}
