package dao

import (
	"das-account-indexer/config"
	"fmt"
	"github.com/scorpiotzh/toolib"
	"gorm.io/gorm"
)

type DbDao struct {
	db *gorm.DB
}

func NewGormDB(dbMysql config.DbMysql) (*DbDao, error) {
	db, err := toolib.NewGormDB(dbMysql.Addr, dbMysql.User, dbMysql.Password, dbMysql.DbName, dbMysql.MaxOpenConn, dbMysql.MaxIdleConn)
	if err != nil {
		return nil, fmt.Errorf("toolib.NewGormDB err: %s", err.Error())
	}
	return &DbDao{db: db}, nil
}
