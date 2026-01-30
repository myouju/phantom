package a

func sendTests() {
	ch := make(chan A[string])
	ch <- A[bool](nil)   // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
	ch <- A[string](nil) // OK
}
