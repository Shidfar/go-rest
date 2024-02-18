package pkg

import (
	"go/ast"
	"go/build"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type PkgInfo struct {
	PkgName string
	Imports map[string]string
	Methods []*ast.Field
}

func GetPackageInfo(serviceName, pkgName string, files map[string]*ast.File) PkgInfo {
	pi := PkgInfo{
		PkgName: pkgName,
	}
	for _, file := range files {
		pi.Imports, _ = getImports(file.Imports)
		pi.Methods, _ = getMethods(file.Decls, serviceName)
	}
	return pi
}

func GetGoPackagePath(path string) string {
	var gopaths []string
	gopaths = filepath.SplitList(build.Default.GOPATH)
	for i, gopath := range gopaths {
		gopaths[i] = filepath.ToSlash(filepath.Join(gopath, "src"))
	}

	var dir string
	var err error
	if dir, err = filepath.Abs(path); err != nil {
		panic(err)
	}
	dir = filepath.ToSlash(dir)
	var rootDir string
	var ok bool
	if rootDir, ok = getModuleRoot(dir); ok {
		return rootDir
	}

	for _, gopath := range gopaths {
		if len(gopath) < len(dir) && strings.EqualFold(gopath, dir[0:len(gopath)]) {
			return dir[len(gopath)+1:]
		}
	}

	return ""
}

func getModuleRoot(path string) (string, bool) {
	modregex := regexp.MustCompile(`module ([^\s]*)`)

	var dir string
	var err error
	if dir, err = filepath.Abs(path); err != nil {
		panic(err)
	}
	dir = filepath.ToSlash(dir)

	modDir := dir
	assumedPart := ""
	for {
		var bs []byte
		if bs, err = os.ReadFile(filepath.Join(modDir, "go.mod")); err == nil {
			return string(modregex.FindSubmatch(bs)[1]) + assumedPart, true
		}
		assumedPart = "/" + filepath.Base(modDir) + assumedPart

		var parentDir string
		if parentDir, err = filepath.Abs(filepath.Join(modDir, "..")); err != nil {
			panic(err)
		}

		if parentDir == modDir {
			break
		}
		modDir = parentDir
	}
	return "", false
}
