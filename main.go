package main

import (
	"errors"
	"fmt"
	"github.com/Shidfar/go-rest/internal/pkg"
	"github.com/Shidfar/go-rest/internal/svc"
	"github.com/Shidfar/go-rest/internal/svc/bp"
	"github.com/Shidfar/go-rest/internal/types"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)
	if len(os.Args) < 2 {
		log.Fatal("missing service name or param name")
	}

	serviceName := os.Args[1]

	pkgPath := pkg.GetGoPackagePath(".")

	start := time.Now()
	defer func() {
		log.Printf("generated service %s (%s) in [%v]\n", serviceName, pkgPath, time.Since(start))
	}()

	var pkgName string
	var files map[string]*ast.File
	var err error
	if pkgName, files, err = parseGoFiles("."); err != nil {
		fmt.Println("failed to parse:", err)
		return
	}

	//fmt.Println(pi)
	pi := pkg.GetPackageInfo(serviceName, pkgName, files)
	var funcs []types.FuncMeta
	for _, fun := range pi.Methods {
		fmeta, _ := makeFuncMeta(pkgName, fun, 500, 3)
		funcs = append(funcs, fmeta)
		// if any other meta presented we could change it here
		//if err =
	}

	if err := os.Mkdir(pkgName+"http", 0777); err != nil {
		if !errors.Is(err, os.ErrExist) {
			log.Fatal(err)
		}
	}

	ef, err := os.Create(pkgName + "http/endpoint.go")
	if err != nil {
		log.Fatal(err)
	}
	tf, _ := os.Create(pkgName + "http/transport.go")
	hf, _ := os.Create(pkgName + "http/handler.go")
	defer ef.Close()
	defer tf.Close()
	defer hf.Close()

	builder := svc.Builder{
		Pkg:         pkgName,
		PkgPath:     pkgPath,
		ServiceName: serviceName,
		Funcs:       funcs,
	}
	endpointTemplate := bp.NewEndpointTemplate(serviceName)
	transportTemplate := bp.NewTransportTemplate(serviceName)
	handlerTemplate := bp.NewHandlerTemplate(serviceName)

	if err := endpointTemplate.Execute(ef, builder); err != nil {
		log.Fatal(err)
	}
	if err := transportTemplate.Execute(tf, builder); err != nil {
		log.Fatal(err)
	}
	if err := handlerTemplate.Execute(hf, builder); err != nil {
		log.Fatal(err)
	}

}

func makeFuncMeta(basePkg string, fun *ast.Field, timeout, retries int) (types.FuncMeta, error) {
	if len(fun.Names) > 1 {
		return types.FuncMeta{}, fmt.Errorf("expected only one func name")
	}
	fname := fun.Names[0].Name

	var ftype *ast.FuncType
	var ok bool
	if ftype, ok = fun.Type.(*ast.FuncType); !ok {
		return types.FuncMeta{}, fmt.Errorf("not a function: %s. Service interface cannot have any embeded interfaces", fname)
	}

	args, _ := types.GetArgs(basePkg, ftype.Params.List)
	ress, _ := types.GetResults(basePkg, ftype.Results.List)
	//fmt.Println(args[0].Name, args[0].Type)
	//fmt.Println(ress[0].Name, ress[0].Type)

	return types.FuncMeta{
		Timeout:       time.Duration(timeout),
		RetryAttempts: retries,
		Arguments:     args,
		Returns:       ress,
		Name:          fname,
	}, nil
}

func filterTests(fileInfo fs.FileInfo) bool {
	if fileInfo.IsDir() {
		return false
	}

	if strings.HasSuffix(fileInfo.Name(), "_test.go") {
		return false
	}

	return true
}

func parseGoFiles(path string) (pkgName string, res map[string]*ast.File, err error) {
	res = make(map[string]*ast.File)
	fileSet := token.NewFileSet()
	pkgs, _ := parser.ParseDir(fileSet, path, filterTests, parser.ParseComments)
	if len(pkgs) > 1 {
		err = fmt.Errorf("only 1 package expected other than *_test: %v", pkgs)
		return
	}
	for key, value := range pkgs {
		//fmt.Printf(" > %s: %#v", key, value)
		res = value.Files
		pkgName = key
		//for k, v := range value.Files {
		//	//type File struct {
		//	//	Doc     *CommentGroup // associated documentation; or nil
		//	//	Package token.Pos     // position of "package" keyword
		//	//	Name    *Ident        // package name
		//	//	Decls   []Decl        // top-level declarations; or nil
		//	//
		//	//	FileStart, FileEnd token.Pos       // start and end of entire file
		//	//	Scope              *Scope          // package scope (this file only)
		//	//	Imports            []*ImportSpec   // imports in this file
		//	//	Unresolved         []*Ident        // unresolved identifiers in this file
		//	//	Comments           []*CommentGroup // list of all comments in the source file
		//	//}
		//	fmt.Println(" > key > ", k, ":")
		//	fmt.Println(" > doc > ", v.Doc)
		//	fmt.Println(" > package > ", v.Package)
		//	fmt.Println(" > name > ", v.Name)
		//	fmt.Println(" > decls > :")
		//	for _, d := range v.Decls {
		//		fmt.Println(" > > > ", d)
		//	}
		//	fmt.Println(" > fileStart > ", v.FileStart)
		//	fmt.Println(" > fileEnd > ", v.FileEnd)
		//	fmt.Println(" > scope > ", v.Scope)
		//	fmt.Println(" > imports > ", v.Imports)
		//	fmt.Println(" > unresolved > ", v.Unresolved)
		//	fmt.Println(" > comments > ", v.Comments)
		//}
	}
	return
}
