package tables

import (
	"time"

	"github.com/dotbitHQ/das-lib/common"
)

type TableAccountInfo struct {
	Id                   uint64                `json:"id" gorm:"column:id;primary_key;type:bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT ''"`
	BlockNumber          uint64                `json:"block_number" gorm:"column:block_number;type:bigint(20) NOT NULL DEFAULT '0' COMMENT ''"`
	BlockTimestamp       uint64                `json:"block_timestamp" gorm:"column:block_timestamp;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'Hash-Index'"`
	Outpoint             string                `json:"outpoint" gorm:"column:outpoint;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'Hash-Index'"`
	AccountId            string                `json:"account_id" gorm:"column:account_id;uniqueIndex:uk_account_id;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'hash of account'"`
	ParentAccountId      string                `json:"parent_account_id" gorm:"column:parent_account_id;index:k_parent_account_id;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	NextAccountId        string                `json:"next_account_id" gorm:"column:next_account_id;index:k_next_account_id;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'hash of next account'"`
	Account              string                `json:"account" gorm:"column:account;index:k_account;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT ''"`
	OwnerChainType       common.ChainType      `json:"owner_chain_type" gorm:"column:owner_chain_type;index:k_oct_o;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	Owner                string                `json:"owner" gorm:"column:owner;index:k_oct_o;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'owner address'"`
	OwnerAlgorithmId     common.DasAlgorithmId `json:"owner_algorithm_id" gorm:"column:owner_algorithm_id;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	ManagerChainType     common.ChainType      `json:"manager_chain_type" gorm:"column:manager_chain_type;index:k_mct_m;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	Manager              string                `json:"manager" gorm:"column:manager;index:k_mct_m;type:varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT '' COMMENT 'manager address'"`
	ManagerAlgorithmId   common.DasAlgorithmId `json:"manager_algorithm_id" gorm:"column:manager_algorithm_id;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	Status               AccountStatus         `json:"status" gorm:"column:status;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	EnableSubAccount     EnableSubAccount      `json:"enable_sub_account" gorm:"column:enable_sub_account;type:smallint(6) NOT NULL DEFAULT '0' COMMENT ''"`
	RenewSubAccountPrice uint64                `json:"renew_sub_account_price" gorm:"column:renew_sub_account_price;type:bigint(20) NOT NULL DEFAULT '0' COMMENT ''"`
	Nonce                uint64                `json:"nonce" gorm:"column:nonce;type:bigint(20) NOT NULL DEFAULT '0' COMMENT ''"`
	RegisteredAt         uint64                `json:"registered_at" gorm:"column:registered_at;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	ExpiredAt            uint64                `json:"expired_at" gorm:"column:expired_at;type:bigint(20) unsigned NOT NULL DEFAULT '0' COMMENT ''"`
	CreatedAt            time.Time             `json:"created_at" gorm:"column:created_at;index:k_created_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT ''"`
	UpdatedAt            time.Time             `json:"updated_at" gorm:"column:updated_at;index:k_updated_at;type:timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT ''"`
}

type AccountStatus int
type EnableSubAccount int

const (
	AccountStatusNormal    AccountStatus = 0
	AccountStatusOnSale    AccountStatus = 1
	AccountStatusOnAuction AccountStatus = 2
	AccountStatusOnLock    AccountStatus = 3

	AccountEnableStatusOff EnableSubAccount = 0
	AccountEnableStatusOn  EnableSubAccount = 1

	TableNameAccountInfo = "t_account_info"
)

func (t *TableAccountInfo) TableName() string {
	return TableNameAccountInfo
}
