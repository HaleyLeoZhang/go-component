package db

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Config struct {
	Name         string        `yaml:"name" json:"name"` // 用于 Trace 识别
	Type         string        `yaml:"type" json:"type"`
	Host         string        `yaml:"host" json:"host"`
	Port         int           `yaml:"port" json:"port"`
	Database     string        `yaml:"database" json:"database"`
	User         string        `yaml:"user" json:"user"`
	Password     string        `yaml:"password" json:"password"`
	MaxIdleConns int           `yaml:"maxIdleConns" json:"maxIdleConns"`
	MaxOpenConns int           `yaml:"maxOpenConns" json:"maxOpenConns"`
	MaxLeftTime  time.Duration `yaml:"maxLeftTime" json:"maxLeftTime"`
}

func New(conf *Config) (db *gorm.DB, err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host, conf.Port, conf.Database)
	db, err = gorm.Open(conf.Type, dsn)

	if err != nil {
		fmt.Println("db.Init.Error", err)
		return nil, err
	}

	db.DB().SetMaxIdleConns(conf.MaxIdleConns)
	db.DB().SetMaxOpenConns(conf.MaxOpenConns)
	db.DB().SetConnMaxLifetime(conf.MaxLeftTime)
	return db, nil
}

func Close(db *gorm.DB) {
	defer db.Close()
}
