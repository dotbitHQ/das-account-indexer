package tables

import "time"

type TableRecordsInfo struct {
	Id        uint64    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	AccountId string    `json:"account_id" gorm:"column:account_id"`
	Account   string    `json:"account" gorm:"column:account"`
	Key       string    `json:"key" gorm:"column:key"`
	Type      string    `json:"type" gorm:"column:type"`
	Label     string    `json:"label" gorm:"column:label"`
	Value     string    `json:"value" gorm:"column:value"`
	Ttl       string    `json:"ttl" gorm:"column:ttl"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at"`
}

const (
	TableNameRecordsInfo = "t_records_info"
)

func (t *TableRecordsInfo) TableName() string {
	return TableNameRecordsInfo
}
