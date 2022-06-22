package tables

import "time"

type TableRecordsInfo struct {
	Id              uint64    `json:"id" gorm:"column:id;primary_key;type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT ''"`
	AccountId       string    `json:"account_id" gorm:"column:account_id;index:k_account_id;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT 'hash of account'"`
	ParentAccountId string    `json:"parent_account_id" gorm:"column:parent_account_id;index:k_parent_account_id;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	Account         string    `json:"account" gorm:"column:account;index:k_account;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT '' COMMENT ''"`
	Key             string    `json:"key" gorm:"column:key;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT ''"`
	Type            string    `json:"type" gorm:"column:type;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT ''"`
	Label           string    `json:"label" gorm:"column:label;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT ''"`
	Value           string    `json:"value" gorm:"column:value;index:k_value,length:768;type:varchar(1024) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT ''"`
	Ttl             string    `json:"ttl" gorm:"column:ttl;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci  NOT NULL DEFAULT ''"`
	CreatedAt       time.Time `json:"created_at" gorm:"column:created_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''"`
}

const (
	TableNameRecordsInfo = "t_records_info"
)

func (t *TableRecordsInfo) TableName() string {
	return TableNameRecordsInfo
}
