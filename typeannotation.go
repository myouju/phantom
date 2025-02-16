package typeannotation

import (
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const doc = "typeannotation checks types with type param for type alias"

var Analyzer = &analysis.Analyzer{
	Name: "typeannotation",
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
			if len(n.Lhs) == len(n.Rhs) {
				for i := range len(n.Lhs) {
					assignableTo(pass, n.Pos(), n.Rhs[i], n.Lhs[i])
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

				if len(valuespec.Names) == len(valuespec.Values) {
					for i := range len(valuespec.Names) {
						assignableTo(pass, valuespec.Pos(), valuespec.Values[i], valuespec.Names[i])
					}
				}
			}
		case *ast.CallExpr:
			signature, ok := pass.TypesInfo.TypeOf(n.Fun).(*types.Signature)
			if !ok {
				return
			}

			if signature.Params().Len() == len(n.Args) {
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
	case types.Type:
		typ1 = val
	}

	switch typ := typ.(type) {
	case ast.Expr:
		typ2 = pass.TypesInfo.TypeOf(typ)
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
