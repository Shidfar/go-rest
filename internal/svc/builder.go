package svc

import "github.com/Shidfar/go-rest/internal/types"

type Builder struct {
	Pkg         string
	PkgPath     string
	ServiceName string
	Funcs       []types.FuncMeta
}
