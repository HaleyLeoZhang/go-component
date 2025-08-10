package xlog

import (
	"context"
	"testing"
)

func TestAll(t *testing.T) {
	var ctx = context.Background()
	ctx = WithLogID(ctx, GenerateLogID())

	Infof(ctx, "测试")
	Infof(ctx, "测试 %d", 2)
	Warnf(ctx, "测试 Warnf")
	Warnf(ctx, "测试 Warnf %d", 2)
	Errorf(ctx, "测试 Err ")
	Errorf(ctx, "测试 Err %d", 2)
}
