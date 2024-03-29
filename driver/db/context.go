package db

import (
	"context"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

// 链路追踪、Ctx超时检测
func Context(ctx context.Context, db *gorm.DB) (err error) {
	err = nil
	if ctx == nil {
		return
	}
	err = ctx.Err()
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	// TODO opentracing  https://www.jaegertracing.io/

	return
}
