package fooservice

import (
	"context"
	"fmt"
)

//go:generate go run github.com/Shidfar/go-rest FooService
type FooService interface {
	// GetObject does so-and-so
	// hello there what is the comment here?
	// #GET
	GetObject(ctx context.Context, req GetObjRequest) (GetObjResponse, error)
}

type GetObjRequest struct {
	Id        string `json:"id,omitempty"`
	PageSize  int    `json:"pageSize,omitempty"`
	Something any    `json:"something,omitempty"`
}

type GetObjResponse struct {
	Id     string   `json:"id,omitempty"`
	Values []string `json:"values,omitempty"`
	Token  string   `json:"token,omitempty"`
	Obj    any      `json:"obj,omitempty"`
}

type service struct {
}

func New() FooService {
	return &service{}
}

func (s service) GetObject(ctx context.Context, req GetObjRequest) (GetObjResponse, error) {
	fmt.Println("Hello world")
	return GetObjResponse{
		Id:     "foo",
		Values: []string{"foo", "bar", "biz", "buzz"},
		Token:  "foo-token",
		Obj: map[string]any{
			"rest": "in peace",
		},
	}, nil
	//return GetObjResponse{}, errors.New("oopsie daisy")
}
