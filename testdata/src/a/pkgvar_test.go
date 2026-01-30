package a

var pkgVar A[string] = A[bool](nil) // want `type annotations are not assignable: a\.A\[bool\] to a\.A\[string\]`
var pkgVarOK A[string] = A[string](nil)
