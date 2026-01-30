package a

func moreTests() {
	// Binary expressions (Comparison)
	var x A[string] = A[string](nil)
	var y A[bool] = A[bool](nil)
	_ = x == x // OK
	_ = x == y // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`

	// Index expressions (Map key)
	type Key[T any] = string
	var m map[Key[string]]int
	var kString Key[string] = "key"
	var kBool Key[bool] = "key"

	// Read
	_ = m[kString] // OK
	_ = m[kBool]   // want `type annotations are not assignable: a\.Key\[bool\] to a\.Key\[string\]`

	// Write
	m[kString] = 1 // OK
	m[kBool] = 1   // want `type annotations are not assignable: a\.Key\[bool\] to a\.Key\[string\]`

	// Built-in copy
	sString := []A[string]{A[string](nil)}
	sBool := []A[bool]{A[bool](nil)}
	copy(sString, sString) // OK
	copy(sString, sBool)   // want `type annotations are not assignable: \[\]a\.A\[bool\] to \[\]a\.A\[string\]`
}
