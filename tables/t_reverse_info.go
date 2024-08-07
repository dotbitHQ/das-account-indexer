package tables

import (
	"github.com/dotbitHQ/das-lib/common"
	"time"
)

type TableReverseInfo struct {
	Id             uint64                   `json:"id" gorm:"column:id;primary_key;type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT ''"`
	BlockNumber    uint64                   `json:"block_number" gorm:"column:block_number;type:bigint(20) NOT NULL DEFAULT '0' COMMENT ''"`
	BlockTimestamp uint64                   `json:"block_timestamp" gorm:"column:block_timestamp;type:bigint(20) NOT NULL DEFAULT '0' COMMENT ''"`
	Outpoint       string                   `json:"outpoint" gorm:"column:outpoint;uniqueIndex:uk_outpoint;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	AlgorithmId    common.DasAlgorithmId    `json:"algorithm_id" gorm:"column:algorithm_id;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	SubAlgorithmId common.DasSubAlgorithmId `json:"sub_algorithm_id" gorm:"column:sub_algorithm_id;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	ChainType      common.ChainType         `json:"chain_type" gorm:"column:chain_type;index:k_address;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	Address        string                   `json:"address" gorm:"column:address;index:k_address;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	Account        string                   `json:"account" gorm:"column:account;index:k_account;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	Capacity       uint64                   `json:"capacity" gorm:"column:capacity;type:bigint(20) NOT NULL DEFAULT '0' COMMENT ''"`
	ReverseType    uint32                   `json:"reverse_type" gorm:"column:reverse_type;type:tinyint(1) NOT NULL DEFAULT '0' COMMENT '0: old reverse type，1：new outpoint struct'"`
	P2shP2wpkh     string                   `json:"p2sh_p2wpkh" gorm:"column:p2sh_p2wpkh; index:k_p2sh_p2wpkh; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	P2tr           string                   `json:"p2tr" gorm:"column:p2tr; index:k_p2tr; type:varchar(255) NOT NULL DEFAULT '' COMMENT '';"`
	CreatedAt      time.Time                `json:"created_at" gorm:"column:created_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''"`
	UpdatedAt      time.Time                `json:"updated_at" gorm:"column:updated_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''"`
}

const (
	TableNameReverseInfo = "t_reverse_info"

	ReverseTypeOld = 0
	ReverseTypeSmt = 1
)

func (t *TableReverseInfo) TableName() string {
	return TableNameReverseInfo
}
