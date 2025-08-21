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
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		switch n := n.(type) {
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
						return
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
				return
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
							return
						}

						if signature.Results().Len() == len(valuespec.Names) {
							for i := range len(valuespec.Names) {
								assignableTo(pass, n.Pos(), signature.Results().At(i), valuespec.Names[i])
							}
						}
					}
				}

			}
		case *ast.CallExpr:
			signature, ok := pass.TypesInfo.TypeOf(n.Fun).(*types.Signature)
			if !ok {
				return
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
						for i := paramsLen - 1; i < argsLen; i++ {
							assignableTo(pass, n.Pos(), n.Args[i], elementType)
						}
					}
				}
			} else {
				for i := range len(n.Args) {
					assignableTo(pass, n.Pos(), n.Args[i], signature.Params().At(i))
				}
			}
		}
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

	if !types.AssignableTo(typ1, typ2) {
		pass.Reportf(pos, "types are not assignable: %v to %v", typ1, typ2)
	}

	alias1, _ := typ1.(*types.Alias)
	alias2, _ := typ2.(*types.Alias)
	if alias1 == nil || alias2 == nil {
		return
	}

	if !types.Identical(alias1.Origin(), alias2.Origin()) {
		return
	}

	args1 := alias1.TypeArgs()
	args2 := alias2.TypeArgs()
	if args1.Len() == 0 || args2.Len() == 0 || args1.Len() != args2.Len() {
		return
	}

	for i := range args1.Len() {
		if !types.AssignableTo(args1.At(i), args2.At(i)) {
			pass.Reportf(pos, "type annotations are not assignable: %v to %v", typ1, typ2)
		}
	}
}
