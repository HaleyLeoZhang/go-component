package db

import (
	"context"
	"github.com/HaleyLeoZhang/go-component/driver/xlog"
	"github.com/jinzhu/gorm"
	"os"
	"testing"
)

var (
	db  *gorm.DB
	ctx context.Context
)

func TestMain(m *testing.M) {
	ctx = context.Background()
	dbConfig := &Config{
		Name:     "local-db",
		Type:     "mysql",
		Host:     "192.168.56.110",
		Port:     3306,
		Database: "curl_avatar",
		User:     "yth_blog",
		Password: "http://hlzblog.top",
	}
	var err error
	db, err = New(dbConfig)
	if err != nil {
		xlog.Errorf("Init DB failed!")
		return
	}
	os.Exit(m.Run())
}
