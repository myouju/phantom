package a

func compositeTests() {
	// Struct literal - key-value
	var _ D = D{field: A[bool](nil)}   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _ D = D{field: A[string](nil)} // OK

	// Struct literal - positional
	var _ D = D{A[bool](nil)}   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _ D = D{A[string](nil)} // OK

	// Slice literal
	var _ []A[string] = []A[string]{A[bool](nil)}   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _ []A[string] = []A[string]{A[string](nil)}  // OK

	// Array literal
	var _ [1]A[string] = [1]A[string]{A[bool](nil)}   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _ [1]A[string] = [1]A[string]{A[string](nil)} // OK

	// Map literal - value
	var _ map[string]A[string] = map[string]A[string]{"k": A[bool](nil)}   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _ map[string]A[string] = map[string]A[string]{"k": A[string](nil)} // OK

	// Map literal - phantom key
	type K[T any] = string
	var _ map[K[string]]int = map[K[string]]int{K[bool]("x"): 1}   // want `type annotations are not assignable: a\.K\[bool\] to a\.K\[string\]`
	var _ map[K[string]]int = map[K[string]]int{K[string]("x"): 1} // OK

	// Nested composite literal - struct with slice field
	type E struct {
		items []A[string]
	}
	var _ E = E{items: []A[string]{A[bool](nil)}}   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _ E = E{items: []A[string]{A[string](nil)}} // OK
}
