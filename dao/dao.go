package dao

import (
	"das-account-indexer/config"
	"das-account-indexer/tables"
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

	// AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
	// It will change existing column’s type if its size, precision, nullable changed.
	// It WON’T delete unused columns to protect your data.
	if err = db.AutoMigrate(
		&tables.TableAccountInfo{},
		&tables.TableBlockInfo{},
		&tables.TableRecordsInfo{},
		&tables.TableReverseInfo{},
	); err != nil {
		return nil, err
	}

	return &DbDao{db: db}, nil
}
