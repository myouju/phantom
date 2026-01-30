package a

func genericWrapperTests() {
	// Assignment of generic named type wrapping phantom type
	var _ Option[A[string]] = Option[A[bool]]{}   // want `type annotations are not assignable: a\.Option\[a\.A\[bool\]\] to a\.Option\[a\.A\[string\]\]`
	var _ Option[A[string]] = Option[A[string]]{} // OK

	// Function call with generic wrapper
	func(_ Option[A[string]]) {}(Option[A[bool]]{}) // want `type annotations are not assignable: a\.Option\[a\.A\[bool\]\] to a\.Option\[a\.A\[string\]\]`
	func(_ Option[A[string]]) {}(Option[A[string]]{}) // OK

	// Interface satisfaction with generic wrapper
	var _ Writer = goodWriter{} // OK
	var _ Writer = badWriter{}  // want `badWriter does not implement a\.Writer \(wrong type for method Write\)`
}
