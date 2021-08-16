package xgin

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"testing"
	"time"
)

func TestNewGin(t *testing.T) {
	cfg := &Config{
		Timeout: 4 * time.Second,
	}
	_ = New(cfg)
	gCtx := &gin.Context{}
	g := NewGin(gCtx)
	_ = LabCancel(t, g.C, 1)
	go LabCancel(t, g.C, 2)
	go LabCancel(t, g.C, 3)
	<-time.After(5 * time.Second)
}

func LabCancel(t *testing.T, cx context.Context, text int) (err error) {

	time.Sleep(2 * time.Second)
	err = cx.Err()
	res := cx.Value("2323333")
	if err != nil {
		err = errors.WithStack(err)
		t.Logf("Cccccc.23333(%v).res(%v).Err(%+v)", text, res, err)
		return
	}
	//t.Logf("Cccccc Err(%v)", cx.Err())
	for {
		select {
		case <-cx.Done():
			t.Logf("Cccccc done(%v)", text)
			t.Logf("Cccccc Err(%v)", cx.Err())
			return
		default:
			t.Logf("Cccccc.Others(%v)", text)
		}
	}
	return
}
