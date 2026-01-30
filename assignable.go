package phantom

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "phantom checks for phantom types"

var AssignableAnalyzer = &analysis.Analyzer{
	Name: "assignable",
	Doc:  doc,
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.Ident)(nil),
		(*ast.AssignStmt)(nil),
		(*ast.DeclStmt)(nil),
		(*ast.CallExpr)(nil),
		(*ast.CompositeLit)(nil),
		(*ast.FuncDecl)(nil),
		(*ast.FuncLit)(nil),
		(*ast.ReturnStmt)(nil),
		(*ast.SendStmt)(nil),
		(*ast.BinaryExpr)(nil),
		(*ast.IndexExpr)(nil),
		(*ast.SwitchStmt)(nil),
	}

	// Check package-level variable declarations
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			gendecl, ok := decl.(*ast.GenDecl)
			if !ok || gendecl.Tok != token.VAR {
				continue
			}
			for _, spec := range gendecl.Specs {
				valuespec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				if len(valuespec.Names) == len(valuespec.Values) {
					for i := range len(valuespec.Names) {
						assignableTo(pass, valuespec.Pos(), valuespec.Values[i], valuespec.Names[i])
					}
				}
			}
		}
	}

	var funcStack []*types.Signature

	inspect.Nodes(nodeFilter, func(n ast.Node, push bool) bool {
		if !push {
			switch n.(type) {
			case *ast.FuncDecl, *ast.FuncLit:
				funcStack = funcStack[:len(funcStack)-1]
			}
			return true
		}

		switch n := n.(type) {
		case *ast.FuncDecl:
			obj := pass.TypesInfo.ObjectOf(n.Name)
			if obj == nil {
				return true
			}
			sig, ok := obj.Type().(*types.Signature)
			if !ok {
				return true
			}
			funcStack = append(funcStack, sig)
			return true
		case *ast.FuncLit:
			sig, ok := pass.TypesInfo.TypeOf(n).(*types.Signature)
			if !ok {
				return true
			}
			funcStack = append(funcStack, sig)
			return true
		case *ast.ReturnStmt:
			if len(funcStack) == 0 || len(n.Results) == 0 {
				return true
			}
			sig := funcStack[len(funcStack)-1]
			results := sig.Results()
			if results.Len() == len(n.Results) {
				for i, expr := range n.Results {
					assignableTo(pass, expr.Pos(), expr, results.At(i))
				}
			}
			return true
		case *ast.AssignStmt:
			switch {
			case len(n.Lhs) == len(n.Rhs):
				for i := range len(n.Lhs) {
					assignableTo(pass, n.Pos(), n.Rhs[i], n.Lhs[i])
				}
			case len(n.Rhs) == 1:
				switch expr := n.Rhs[0].(type) {
				case *ast.CallExpr:
					signature, ok := pass.TypesInfo.TypeOf(expr.Fun).(*types.Signature)
					if !ok {
						return true
					}

					if signature.Results().Len() == len(n.Lhs) {
						for i := range len(n.Lhs) {
							assignableTo(pass, n.Pos(), signature.Results().At(i), n.Lhs[i])
						}
					}
				case *ast.IndexExpr:
					tuple, _ := pass.TypesInfo.TypeOf(expr).(*types.Tuple)
					if len(n.Lhs) == 2 && tuple.Len() == 2 {
						assignableTo(pass, n.Pos(), tuple.At(0), n.Lhs[0])
					}
				case *ast.TypeAssertExpr:
					if len(n.Lhs) == 2 {
						assignableTo(pass, n.Pos(), expr.Type, n.Lhs[0])
					}
				}
			}
		case *ast.DeclStmt:
			gendecl, ok := n.Decl.(*ast.GenDecl)
			if !ok || gendecl.Tok != token.VAR {
				return true
			}

			for _, spec := range gendecl.Specs {
				valuespec, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}

				switch {
				case len(valuespec.Names) == len(valuespec.Values):
					for i := range len(valuespec.Names) {
						assignableTo(pass, valuespec.Pos(), valuespec.Values[i], valuespec.Names[i])
					}
				case len(valuespec.Values) == 1:
					switch expr := valuespec.Values[0].(type) {
					case *ast.CallExpr:
						signature, ok := pass.TypesInfo.TypeOf(expr.Fun).(*types.Signature)
						if !ok {
							return true
						}

						if signature.Results().Len() == len(valuespec.Names) {
							for i := range len(valuespec.Names) {
								assignableTo(pass, n.Pos(), signature.Results().At(i), valuespec.Names[i])
							}
						}
					}
				}

			}
		case *ast.CompositeLit:
			litType := pass.TypesInfo.TypeOf(n)
			if litType == nil {
				return true
			}
			switch u := litType.Underlying().(type) {
			case *types.Struct:
				for i, elt := range n.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						ident, ok := kv.Key.(*ast.Ident)
						if !ok {
							continue
						}
						for j := range u.NumFields() {
							field := u.Field(j)
							if field.Name() == ident.Name {
								assignableTo(pass, kv.Pos(), kv.Value, field)
								break
							}
						}
					} else if i < u.NumFields() {
						assignableTo(pass, elt.Pos(), elt, u.Field(i))
					}
				}
			case *types.Slice:
				elemType := u.Elem()
				for _, elt := range n.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						assignableTo(pass, kv.Pos(), kv.Value, elemType)
					} else {
						assignableTo(pass, elt.Pos(), elt, elemType)
					}
				}
			case *types.Array:
				elemType := u.Elem()
				for _, elt := range n.Elts {
					if kv, ok := elt.(*ast.KeyValueExpr); ok {
						assignableTo(pass, kv.Pos(), kv.Value, elemType)
					} else {
						assignableTo(pass, elt.Pos(), elt, elemType)
					}
				}
			case *types.Map:
				keyType := u.Key()
				elemType := u.Elem()
				for _, elt := range n.Elts {
					kv, ok := elt.(*ast.KeyValueExpr)
					if !ok {
						continue
					}
					assignableTo(pass, kv.Pos(), kv.Key, keyType)
					assignableTo(pass, kv.Pos(), kv.Value, elemType)
				}
			}
		case *ast.SendStmt:
			chanType, ok := pass.TypesInfo.TypeOf(n.Chan).Underlying().(*types.Chan)
			if !ok {
				return true
			}
			assignableTo(pass, n.Arrow, n.Value, chanType.Elem())
		case *ast.BinaryExpr:
			if n.Op == token.EQL || n.Op == token.NEQ {
				assignableTo(pass, n.OpPos, n.Y, n.X)
			}
		case *ast.IndexExpr:
			// Check map key type on map index access
			mapType, ok := pass.TypesInfo.TypeOf(n.X).Underlying().(*types.Map)
			if !ok {
				return true
			}
			assignableTo(pass, n.Lbrack, n.Index, mapType.Key())
		case *ast.SwitchStmt:
			if n.Tag != nil {
				for _, stmt := range n.Body.List {
					if cc, ok := stmt.(*ast.CaseClause); ok {
						for _, expr := range cc.List {
							assignableTo(pass, expr.Pos(), expr, n.Tag)
						}
					}
				}
			}
		case *ast.CallExpr:
			// Handle builtin functions whose parameters are instantiated per-call
			if ident, ok := n.Fun.(*ast.Ident); ok {
				obj := pass.TypesInfo.ObjectOf(ident)
				if obj != nil {
					if builtin, ok := obj.(*types.Builtin); ok {
						switch builtin.Name() {
						case "copy":
							if len(n.Args) == 2 {
								assignableTo(pass, n.Pos(), n.Args[1], n.Args[0])
							}
							return true
						}
					}
				}
			}

			signature, ok := pass.TypesInfo.TypeOf(n.Fun).(*types.Signature)
			if !ok {
				return true
			}

			paramsLen := signature.Params().Len()
			argsLen := len(n.Args)

			if signature.Variadic() {
				// Check fixed parameters (exclude the variadic parameter)
				for i := 0; i < paramsLen-1; i++ {
					assignableTo(pass, n.Pos(), n.Args[i], signature.Params().At(i))
				}
				// Check variadic arguments against the variadic parameter's element type
				if argsLen > paramsLen-1 {
					variadicParam := signature.Params().At(paramsLen - 1)
					if slice, ok := variadicParam.Type().(*types.Slice); ok {
						elementType := slice.Elem()
						
						// Check if the call uses slice expansion (...)
						if n.Ellipsis.IsValid() {
							// Handle slice expansion: the last argument should be assignable to []T
							lastArgIndex := argsLen - 1
							if lastArgIndex >= paramsLen-1 {
								assignableTo(pass, n.Pos(), n.Args[lastArgIndex], variadicParam.Type())
							}
						} else {
							// Handle regular variadic arguments
							for i := paramsLen - 1; i < argsLen; i++ {
								assignableTo(pass, n.Pos(), n.Args[i], elementType)
							}
						}
					}
				}
			} else {
				for i := range len(n.Args) {
					assignableTo(pass, n.Pos(), n.Args[i], signature.Params().At(i))
				}
			}
		}
		return true
	})

	return nil, nil
}

func assignableTo(pass *analysis.Pass, pos token.Pos, val, typ any) {
	var typ1, typ2 types.Type
	switch val := val.(type) {
	case ast.Expr:
		typ1 = pass.TypesInfo.TypeOf(val)
	case types.Object:
		typ1 = val.Type()
	case types.Type:
		typ1 = val
	}

	switch typ := typ.(type) {
	case ast.Expr:
		typ2 = pass.TypesInfo.TypeOf(typ)
	case types.Object:
		typ2 = typ.Type()
	case types.Type:
		typ2 = typ
	}

	if typ1 == nil || typ2 == nil {
		return
	}

	if !types.AssignableTo(typ1, typ2) {
		pass.Reportf(pos, "types are not assignable: %v to %v", typ1, typ2)
	}

	if !phantomAssignable(typ1, typ2) {
		pass.Reportf(pos, "type annotations are not assignable: %v to %v", typ1, typ2)
	}

	// Check interface satisfaction for phantom type mismatches
	ifaceType, ok := typ2.Underlying().(*types.Interface)
	if !ok {
		return
	}
	mset := types.NewMethodSet(typ1)
	for i := range ifaceType.NumMethods() {
		ifaceMethod := ifaceType.Method(i)
		ifaceSig, ok := ifaceMethod.Type().(*types.Signature)
		if !ok {
			continue
		}
		sel := mset.Lookup(ifaceMethod.Pkg(), ifaceMethod.Name())
		if sel == nil {
			continue
		}
		concreteSig, ok := sel.Type().(*types.Signature)
		if !ok {
			continue
		}
		if !signaturesPhantomMatch(ifaceSig, concreteSig) {
			pass.Reportf(pos, "%v does not implement %v (wrong type for method %s)", typ1, typ2, ifaceMethod.Name())
			return
		}
	}
}

// signaturesPhantomMatch checks if two method signatures match with respect to phantom types.
func signaturesPhantomMatch(iface, concrete *types.Signature) bool {
	for i := range iface.Params().Len() {
		if !phantomAssignable(concrete.Params().At(i).Type(), iface.Params().At(i).Type()) {
			return false
		}
	}
	for i := range iface.Results().Len() {
		if !phantomAssignable(concrete.Results().At(i).Type(), iface.Results().At(i).Type()) {
			return false
		}
	}
	return true
}

// typeArgs extracts type arguments from alias or named generic types.
// Returns nil if the type has no type arguments.
func typeArgs(t types.Type) (origin types.Type, args *types.TypeList) {
	switch t := t.(type) {
	case *types.Alias:
		if ta := t.TypeArgs(); ta != nil && ta.Len() > 0 {
			return t.Origin(), ta
		}
	case *types.Named:
		if ta := t.TypeArgs(); ta != nil && ta.Len() > 0 {
			return t.Origin(), ta
		}
	}
	return nil, nil
}

// phantomAssignable checks whether t1 is assignable to t2, including
// recursive comparison of phantom type arguments on alias and named generic types.
func phantomAssignable(t1, t2 types.Type) bool {
	if !types.AssignableTo(t1, t2) {
		return false
	}

	// Check alias and named generic types
	origin1, args1 := typeArgs(t1)
	origin2, args2 := typeArgs(t2)
	if origin1 != nil && origin2 != nil {
		if types.Identical(origin1, origin2) {
			if args1.Len() == args2.Len() {
				for i := range args1.Len() {
					if !phantomAssignable(args1.At(i), args2.At(i)) {
						return false
					}
				}
			}
		}
		return true
	}

	// Check slice types with generic element types
	slice1, ok1 := t1.(*types.Slice)
	slice2, ok2 := t2.(*types.Slice)
	if ok1 && ok2 {
		return phantomAssignable(slice1.Elem(), slice2.Elem())
	}

	return true
}
