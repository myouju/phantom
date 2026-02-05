package a

func missingTests() {
	// Switch statement comparisons
	var tag A[string] = A[string](nil)
	var valString A[string] = A[string](nil)
	var valBool A[bool] = A[bool](nil)

	switch tag {
	case valString: // OK
	case valBool: // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	}

	// Channel receive assignment
	chBool := make(chan A[bool])
	var sString A[string]

	sString = <-chBool // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	_ = sString

	select {
	case sString = <-chBool: // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	default:
	}

	// Range over map
	m := map[A[string]]A[int]{}
	for k, _ := range m {
		var _ A[bool] = k // want `type annotations are not assignable: a\.A\[string\] to a\.A\[bool\]`
	}
}
