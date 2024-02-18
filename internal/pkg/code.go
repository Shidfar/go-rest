package pkg

import (
	"fmt"
	"go/ast"
	"strings"
)

func getImports(imps []*ast.ImportSpec) (res map[string]string, err error) {
	res = make(map[string]string)
	extractNameAndPath := func(spec *ast.ImportSpec) (name, path string) {
		if spec.Name != nil {
			name = spec.Name.Name
			path = spec.Path.Value
			return
		}
		parts := strings.Split(spec.Path.Value, "/")
		name = strings.Replace(parts[len(parts)-1], `"`, "", -1)
		path = strings.Replace(spec.Path.Value, `"`, "", -1)
		return
	}
	for _, imp := range imps {
		name, path := extractNameAndPath(imp)
		if currPath, ok := res[name]; ok && currPath != path {
			err = fmt.Errorf("import path missmatch for %s: %s and %s", name, currPath, path)
			return
		}
		res[name] = path
	}
	return
}

func getMethods(decls []ast.Decl, iface string) (res []*ast.Field, err error) {
	for _, decl := range decls {
		switch typedDecl := decl.(type) {
		case *ast.BadDecl:
			//fmt.Printf(" > BadDecl > %#v\n", typedDecl)
			continue
		case *ast.FuncDecl:
			//fmt.Printf(" > FuncDecl > %#v\n", typedDecl)
			continue
		case *ast.GenDecl:
			//fmt.Printf(" > GenDecl > %#v\n", typedDecl)
			methods := parseGenDecl(typedDecl, iface)
			if len(methods) > 0 && len(methods[0].Doc.List) > 0 {
				res = methods
				//fmt.Printf("[*] %#v\n", methods[0].Doc.List[0].Text)
				//fmt.Printf("[*] %#v\n", methods[0].Doc.List[1].Text)
				//fmt.Printf("[*] %#v\n", methods[0].Doc.List[2].Text)
				//fmt.Printf("[*] %#v\n", methods[0].Names[0].Name)
				//fmt.Printf("[*] %#v %#v.%#v\n", methods[0].Type.(*ast.FuncType).Params.List[0].Names[0].Name, methods[0].Type.(*ast.FuncType).Params.List[0].Type.(*ast.SelectorExpr).X.(*ast.Ident).Name, methods[0].Type.(*ast.FuncType).Params.List[0].Type.(*ast.SelectorExpr).Sel.Name)
				//fmt.Printf("[*] %#v %#v.%#v\n", methods[0].Type.(*ast.FuncType).Params.List[1].Names[0].Name, methods[0].Type.(*ast.FuncType).Params.List[1].Type.(*ast.Ident).Name, methods[0].Type.(*ast.FuncType).Params.List[1].Type.(*ast.Ident).Obj.Decl.(*ast.TypeSpec).Name.Name)
			}
			//else {
			//	fmt.Println("[-]")
			//}
			//fmt.Println()
			continue
		}
	}
	return
}

func parseGenDecl(decl *ast.GenDecl, iface string) (res []*ast.Field) {
	for _, spec := range decl.Specs {
		//fmt.Printf("    ---- %#v\n", spec)
		switch typedSpec := spec.(type) {
		case *ast.ImportSpec:
			continue
			//fmt.Printf(" -- ImportSpec --> %#v", typedSpec)
		case *ast.TypeSpec:
			//fmt.Printf(" -- TypeSpec --> %#v", typedSpec)
			if typedSpec.Name.Name != iface {
				continue
			}
			if t, ok := typedSpec.Type.(*ast.InterfaceType); ok {
				res = t.Methods.List
			}
		case *ast.ValueSpec:
			continue
			//fmt.Printf(" -- ValueSpec --> %#v", typedSpec)
		}
		fmt.Println()
	}
	return
}
