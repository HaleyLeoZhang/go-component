package xgin

import (
	"context"
	"testing"
	"time"
	"github.com/pkg/errors"
)

func TestNewGin(t *testing.T) {
	timeout := 1 * time.Second
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	ctx = context.WithValue(ctx,  "2323333", 1)
	defer cancelFunc()
	Cccccc(t, ctx, 1)
	//go Cccccc(t, ctx, 2)
	//go Cccccc(t, ctx, 3)
	//go Cccccc(t, ctx, 4)
	//for {
	//	select {
	//	case <-ctx.Done():
	//		t.Log("done")
	//		return
	//	}
	//	time.Sleep(2 * time.Second)
	//	go Cccccc(t, ctx, 1)
	//	//cancelFunc()
	//
	//}
}
func Cccccc(t *testing.T, cx context.Context, text int) (err error){

	time.Sleep(2 * time.Second)
	err = cx.Err()
	res := cx.Value("2323333")
	if err !=nil {
		err = errors.WithStack(err)
		t.Logf("Cccccc.23333(%v).res(%v).Err(%+v)", text, res, err)
		return
	}
		select {
		case <-cx.Done():
			t.Logf("Cccccc.done(%v)", text)
			t.Logf("Cccccc.Err(%v)", cx.Err())
			return
		default:
			t.Logf("Cccccc.Others(%v)", text)
		}
		t.Logf("Cccccc.except(%v)", text)
	//for {
	//	select {
	//	case <-cx.Done():
	//		t.Logf("Cccccc.done(%v)", text)
	//		return
	//	//default:
	//	//	t.Logf("Cccccc.Start(%v)", text)
	//	//	time.Sleep(1000 * time.Second)
	//	//	t.Logf("Cccccc.Others(%v)", text)
	//	}
	//}
	return
}
