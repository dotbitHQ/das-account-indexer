package toolib

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func NewGormDataBase(addr, user, password, dbName string, maxOpenConn, maxIdleConn int) (*gorm.DB, error) {
	const conn = "%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=true&loc=Local"
	dataSource := fmt.Sprintf(conn, user, password, addr, dbName)
	db, err := gorm.Open("mysql", dataSource)
	if err != nil {
		return nil, fmt.Errorf("gorm open :%s", err)
	}
	db.DB().SetMaxOpenConns(maxOpenConn)
	db.DB().SetMaxIdleConns(maxIdleConn)
	return db, nil
}
