package types

import (
	"fmt"
	"go/ast"
	"strings"
)

const (
	Signature    = UsageCtx("Signature")
	Field        = UsageCtx("Field")
	DynamicField = UsageCtx("DynamicField")
	Use          = UsageCtx("Use")
	UseLower     = UsageCtx("UseLower")
)

func GetArgs(basePkg string, fields []*ast.Field) (Args, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("no fields")
	}

	if !isCtx(fields[0]) {
		return nil, fmt.Errorf("first param must be of type context.Context")
	}
	return parseArgs(basePkg, fields[1:])
}

func GetResults(basePkg string, results []*ast.Field) (Args, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no results")
	}
	if !isErr(results[len(results)-1]) {
		return nil, fmt.Errorf("last result must be an error")
	}
	return parseArgs(basePkg, results[:len(results)-1])
}

func parseArgs(basePkg string, fields []*ast.Field) (Args, error) {
	var args []Arg
	for i, f := range fields {
		t, err := getType(basePkg, basePkg, f.Type)
		if err != nil {
			return nil, err
		}

		if len(f.Names) == 0 {
			args = append(args, Arg{
				Name: nameOrDefault(f, fmt.Sprintf("arg%d", i)),
				Type: t,
			})
			continue
		}
		for _, name := range f.Names {
			args = append(args, Arg{
				Name: name.Name,
				Type: t,
			})
		}
	}
	return args, nil
}

func nameOrDefault(field *ast.Field, def string) string {
	if len(field.Names) < 1 {
		return def
	}
	return field.Names[0].Name
}

func getType(basePkg, pkg string, expr ast.Expr) (TypeDef, error) {
	switch t := expr.(type) {
	case *ast.Ident:
		var typ *Named
		switch {
		case basePkg != pkg:
			typ = &Named{Name: t.Name, Pkg: pkg}
		case t.Name == strings.ToLower(t.Name):
			typ = &Named{Name: t.Name}
		default:
			typ = &Named{Name: t.Name, Pkg: pkg}
		}
		return typ, nil

	case *ast.StarExpr:
		// This case handles pointers.
		td, err := getType(basePkg, pkg, t.X)
		if err != nil {
			return nil, err
		}
		return &Ptr{To: td}, nil

	case *ast.SelectorExpr:
		pkg := t.X.(*ast.Ident).Name
		return getType(basePkg, pkg, t.Sel)

	case *ast.ArrayType:
		td, err := getType(basePkg, pkg, t.Elt)
		if err != nil {
			return nil, err
		}
		return &Slice{Vals: td}, nil

	case *ast.Ellipsis:
		td, err := getType(basePkg, pkg, t.Elt)
		if err != nil {
			return nil, err
		}
		return &Vararg{Vals: td}, nil

	case *ast.MapType:
		keys, err := getType(basePkg, pkg, t.Key)
		if err != nil {
			return nil, err
		}
		vals, err := getType(basePkg, pkg, t.Value)
		if err != nil {
			return nil, err
		}
		return &Map{Keys: keys, Vals: vals}, nil

	case *ast.InterfaceType:
		return &Any{}, nil

	default:
		return nil, fmt.Errorf("unexpected type: %T", expr)
	}
}

func isCtx(field *ast.Field) bool {
	switch f := field.Type.(type) {
	case *ast.SelectorExpr:
		ident, ok := f.X.(*ast.Ident)
		if !ok {
			return false
		}
		return f.Sel.Name == "Context" && ident.Name == "context"
	default:
		return false
	}
}

func isErr(f *ast.Field) bool {
	switch f := f.Type.(type) {
	case *ast.Ident:
		return f.Name == "error"
	default:
		return false
	}
}
