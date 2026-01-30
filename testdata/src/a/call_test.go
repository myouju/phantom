package a

func callTests() {
	// Direct function call
	func(A[string]) {}(A[bool](nil)) // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`

	// Variadic function calls
	p[string]()                                       // OK
	p[any](A[bool](nil))                              // OK
	p[bool]([]A[bool]{A[bool](nil), A[bool](nil)}...) // OK
	p[string](A[bool](nil))                           // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	p[string](A[string](nil), A[bool](nil))           // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`

	// Append with slice expansion
	slice1 := []A[string]{A[string](nil)}
	slice2 := []A[string]{A[string](nil)}
	slice3 := []A[bool]{A[bool](nil)}
	_ = append(slice1, slice2...) // OK
	_ = append(slice1, slice3...) // want `type annotations are not assignable: \[\]a\.A\[bool\] to \[\]a\.A\[string\]`

	// Non-phantom variadic (no false positive)
	type Mod[T any] interface{ Apply(T) }
	type SelectQuery struct{}
	buildMods := func() []Mod[*SelectQuery] { return nil }
	mods := []Mod[*SelectQuery]{}
	filterMods := buildMods()
	_ = append(mods, filterMods...) // OK

	// Method call with nested phantom types
	c := C{caller: b{}}
	c.caller.M(A[A[bool]](nil))   // want `type annotations are not assignable: a\.A\[a\.A\[bool\]\] to a\.A\[a\.A\[string\]\]`
	c.caller.M(A[A[string]](nil)) // OK

	c2 := C{caller: invalidB{}}    // want `invalidB does not implement a\.B \(wrong type for method M\)`
	c2.caller.M(A[A[bool]](nil))   // want `type annotations are not assignable: a\.A\[a\.A\[bool\]\] to a\.A\[a\.A\[string\]\]`
	c2.caller.M(A[A[string]](nil)) // OK

	// defer/go statement function calls
	defer func(_ A[string]) {}(A[bool](nil)) // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	go func(_ A[string]) {}(A[bool](nil))    // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
}
