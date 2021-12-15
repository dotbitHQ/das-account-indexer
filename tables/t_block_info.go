package tables

import "time"

type TableBlockInfo struct {
	Id          uint64    `json:"id" gorm:"column:id;primary_key;AUTO_INCREMENT"`
	BlockNumber uint64    `json:"block_number" gorm:"column:block_number"`
	BlockHash   string    `json:"block_hash" gorm:"column:block_hash"`
	ParentHash  string    `json:"parent_hash" gorm:"column:parent_hash"`
	CreatedAt   time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"column:updated_at"`
}

const (
	TableNameBlockInfo = "t_block_info"
)

func (t *TableBlockInfo) TableName() string {
	return TableNameBlockInfo
}
