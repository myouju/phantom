package a

import "fmt"

type A[T any] = any

func p[T any](_ ...A[T]) {}

type B interface {
	M(A[A[string]])
}

type b struct{}

func (b) M(_ A[A[string]]) {}
func (b) N() A[string]     { return A[bool](nil) } // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`

type C struct {
	caller B
}

type D struct {
	field A[string]
}

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
	fmt.Printf("I'm %s", "Gopher")                                      // OK
	p[string]()                                                         // OK
	p[any](A[bool](nil))                                                // OK
	p[bool]([]A[bool]{A[bool](nil), A[bool](nil)}...)                   // OK
	p[string](A[bool](nil))                                             // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	p[string](A[string](nil), A[bool](nil))                             // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`

	// Test append function with slice expansion
	slice1 := []A[string]{A[string](nil)}
	slice2 := []A[string]{A[string](nil)}
	slice3 := []A[bool]{A[bool](nil)}
	_ = append(slice1, slice2...) // OK - same types
	_ = append(slice1, slice3...) // want `type annotations are not assignable: \[\]a\.A\[bool\] to \[\]a\.A\[string\]`

	// Test with bob.Mod types like in zoo project
	type Mod[T any] interface{ Apply(T) }
	type SelectQuery struct{}

	buildMods := func() []Mod[*SelectQuery] { return nil }
	mods := []Mod[*SelectQuery]{}
	filterMods := buildMods()
	_ = append(mods, filterMods...) // OK

	c := C{caller: b{}}
	c.caller.M(A[A[bool]](nil))   // want `type annotations are not assignable: a\.A\[a\.A\[bool\]\] to a\.A\[a\.A\[string\]\]`
	c.caller.M(A[A[string]](nil)) // OK

	var _ D = D{field: A[bool](nil)}   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	var _ D = D{field: A[string](nil)} // OK
}
