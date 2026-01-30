package a

func (b) N() A[string] { return A[bool](nil) } // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`

func multiReturn() (A[string], A[bool]) {
	return A[bool](nil), A[bool](nil) // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
}

func multiReturnOK() (A[string], A[bool]) {
	return A[string](nil), A[bool](nil) // OK
}
