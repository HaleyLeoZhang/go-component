package xlog

import (
	"context"
	"testing"
)

func TestInfof(t *testing.T) {
	var ctx = context.Background()
	type args struct {
		ctx    context.Context
		format string
		args   []interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			args: args{
				ctx:    ctx,
				format: "测试1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Infof(tt.args.ctx, tt.args.format)
		})
	}
}
