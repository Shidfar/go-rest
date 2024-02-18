package fooservice

import (
	"context"
	"reflect"
	"testing"
)

func Test_service_GetObject(t *testing.T) {
	type args struct {
		ctx context.Context
		req GetObjRequest
	}
	tests := []struct {
		name    string
		args    args
		want    GetObjResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := service{}
			got, err := s.GetObject(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}
