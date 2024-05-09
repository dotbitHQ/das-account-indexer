package tables

import (
	"time"
)

const (
	TableNameDidCellInfo = "t_did_cell_info"
)

type TableDidCellInfo struct {
	Id          uint64    `json:"id" gorm:"column:id;primaryKey;type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT ''"`
	BlockNumber uint64    `json:"block_number" gorm:"column:block_number;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	Outpoint    string    `json:"outpoint" gorm:"column:outpoint;uniqueIndex:uk_op;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'Hash-Index'"`
	AccountId   string    `json:"account_id" gorm:"column:account_id;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'hash of account'"`
	Account     string    `json:"account" gorm:"column:account;index:account;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	Args        string    `json:"args" gorm:"column:args;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'manager address'"`
	ExpiredAt   uint64    `json:"expired_at" gorm:"column:expired_at;index:k_expired_at;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''"`
}

func (t *TableDidCellInfo) TableName() string {
	return TableNameDidCellInfo
}

type DidCellStatus int

const (
	DidCellStatusNormal  DidCellStatus = 1
	DidCellStatusExpired DidCellStatus = 2
)

func (t *TableDidCellInfo) IsExpired() bool {
	if int64(t.ExpiredAt) <= time.Now().Unix() {
		return true
	}
	return false
}
