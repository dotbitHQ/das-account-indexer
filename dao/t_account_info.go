package dao

import (
	"das-account-indexer/tables"
	"errors"
	"github.com/dotbitHQ/das-lib/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

func (d *DbDao) UpdateAccountInfo(account *tables.TableAccountInfo, records []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"block_number", "block_timestamp", "outpoint", "next_account_id",
				"owner_algorithm_id", "owner_chain_type", "owner",
				"manager_algorithm_id", "manager_chain_type", "manager",
				"status", "registered_at", "expired_at",
			}),
		}).Create(&account).Error; err != nil {
			return err
		}

		if err := tx.Where(" account_id=? ", account.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}

		if len(records) > 0 {
			if err := tx.Create(&records).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *DbDao) TransferAccountToDid(accountInfo tables.TableAccountInfo, didCellInfo tables.TableDidCellInfo, recordsInfos []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("block_number", "outpoint", "status").
			Where("account_id = ?", accountInfo.AccountId).
			Updates(accountInfo).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id = ?", didCellInfo.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}
		if len(recordsInfos) > 0 {
			if err := tx.Create(&recordsInfos).Error; err != nil {
				return err
			}
		}

		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"args", "account", "expired_at",
				"created_at", "updated_at",
			}),
		}).Create(&didCellInfo).Error; err != nil {
			return err
		}

		return nil
	})
}

func (d *DbDao) UpdateAccountInfoList(accounts []tables.TableAccountInfo, records []tables.TableRecordsInfo, accountIdList []string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"block_number", "block_timestamp", "outpoint", "next_account_id",
				"owner_algorithm_id", "owner_chain_type", "owner",
				"manager_algorithm_id", "manager_chain_type", "manager",
				"status", "registered_at", "expired_at",
			}),
		}).Create(&accounts).Error; err != nil {
			return err
		}

		if len(accountIdList) > 0 {
			if err := tx.Where(" account_id IN(?) ", accountIdList).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
				return err
			}
		}

		if len(records) > 0 {
			if err := tx.Create(&records).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *DbDao) FindAccountInfoByAccountId(accountId string) (accountInfo tables.TableAccountInfo, err error) {
	err = d.db.Where(" account_id=? ", accountId).Find(&accountInfo).Error
	return
}

func (d *DbDao) FindAccountInfoListByAccountIds(accountIds []string) (list []tables.TableAccountInfo, err error) {
	err = d.db.Where(" account_id IN(?) ", accountIds).Find(&list).Error
	return
}

func (d *DbDao) FindAccountListByAddress(chainType common.ChainType, address string) (list []tables.TableAccountInfo, err error) {
	err = d.db.Where(" owner_chain_type=? AND owner=? and expired_at >= ?", chainType, address, time.Now().Unix()-90*86400).Find(&list).Error
	return
}

func (d *DbDao) FindAccountNameListByAddress(chainType common.ChainType, address, role string) (list []tables.TableAccountInfo, err error) {
	if role == "" || role == "owner" {
		err = d.db.Select("account,registered_at,expired_at").Where(" owner_chain_type=? AND owner=? AND `status`!=? and expired_at >= ?", chainType, address, tables.AccountStatusOnLock, time.Now().Unix()-90*86400).Find(&list).Error
	} else if role == "manager" {
		err = d.db.Select("account,registered_at,expired_at").Where(" manager_chain_type=? AND manager=? AND `status`!=? and expired_at >= ?", chainType, address, tables.AccountStatusOnLock, time.Now().Unix()-90*86400).Find(&list).Error
	}
	return
}

func (d *DbDao) EnableSubAccount(accountInfo tables.TableAccountInfo) error {
	return d.db.Select("block_number", "block_timestamp", "outpoint", "enable_sub_account", "renew_sub_account_price").
		Where("account_id = ?", accountInfo.AccountId).Updates(accountInfo).Error
}

func (d *DbDao) CreateSubAccount(subAccountIds []string, accountInfos []tables.TableAccountInfo, parentAccountInfo tables.TableAccountInfo, records []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if len(subAccountIds) > 0 {
			if err := tx.Where(" account_id IN(?) ", subAccountIds).
				Delete(&tables.TableRecordsInfo{}).Error; err != nil {
				return err
			}
		}
		if len(records) > 0 {
			if err := tx.Create(&records).Error; err != nil {
				return err
			}
		}
		if len(accountInfos) > 0 {
			if err := tx.Clauses(clause.Insert{
				Modifier: "IGNORE",
			}).Create(&accountInfos).Error; err != nil {
				return err
			}
		}
		if parentAccountInfo.AccountId != "" {
			if err := tx.Select("block_number", "block_timestamp", "outpoint").
				Where("account_id = ?", parentAccountInfo.AccountId).
				Updates(parentAccountInfo).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *DbDao) EditOwnerSubAccount(accountInfo tables.TableAccountInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("block_number", "block_timestamp", "outpoint",
			"manager_chain_type", "manager", "manager_algorithm_id",
			"owner_chain_type", "owner", "owner_algorithm_id", "nonce").
			Where("account_id = ?", accountInfo.AccountId).
			Updates(accountInfo).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id = ?", accountInfo.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (d *DbDao) EditManagerSubAccount(accountInfo tables.TableAccountInfo) error {
	return d.db.Select("block_number", "block_timestamp", "outpoint",
		"manager_chain_type", "manager", "manager_algorithm_id", "nonce").
		Where("account_id = ?", accountInfo.AccountId).
		Updates(accountInfo).Error
}

func (d *DbDao) EditRecordsSubAccount(accountInfo tables.TableAccountInfo, recordsInfos []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("block_number", "block_timestamp", "outpoint", "nonce").
			Where("account_id = ?", accountInfo.AccountId).
			Updates(accountInfo).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id = ?", accountInfo.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}

		if len(recordsInfos) > 0 {
			if err := tx.Create(&recordsInfos).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *DbDao) RenewSubAccount(accountInfos []tables.TableAccountInfo) error {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{
			"block_number", "block_timestamp", "outpoint",
			"expired_at", "nonce",
		}),
	}).Create(&accountInfos).Error
}

func (d *DbDao) RecycleSubAccount(accountId []string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id IN(?)", accountId).Delete(&tables.TableAccountInfo{}).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id IN(?)", accountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}

		return nil
	})
}

func (d *DbDao) GetAccountInfoByParentAccountId(parentAccountId string) (accountInfos []tables.TableAccountInfo, err error) {
	err = d.db.Where("parent_account_id=?", parentAccountId).Find(&accountInfos).Error
	return
}

func (d *DbDao) GetAccountInfoByAccountId(accountId string) (accountInfo tables.TableAccountInfo, err error) {
	err = d.db.Where("account_id=?", accountId).First(&accountInfo).Error
	return
}

func (d *DbDao) RecycleExpiredAccount(accountInfo tables.TableAccountInfo, accountId string, enableSubAccount uint8) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("block_number", "block_timestamp", "outpoint", "next_account_id").
			Where("account_id=?", accountInfo.AccountId).
			Updates(accountInfo).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id=?", accountId).Delete(&tables.TableAccountInfo{}).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id=?", accountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}

		if enableSubAccount == 1 {
			if err := tx.Where("parent_account_id=?", accountId).Delete(&tables.TableAccountInfo{}).Error; err != nil {
				return err
			}

			if err := tx.Where("parent_account_id=?", accountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (d *DbDao) UpdateAccountOutpoint(accountId, outpoint string) error {
	return d.db.Model(tables.TableAccountInfo{}).
		Where("account_id=?", accountId).
		Updates(map[string]interface{}{
			"outpoint": outpoint,
		}).Error
}

func (d *DbDao) GetSubAccountListByParentAccountId(parentAccountId string, limit, offset int) (list []tables.TableAccountInfo, err error) {
	err = d.db.Where("parent_account_id=?", parentAccountId).
		Order("account").Limit(limit).Offset(offset).
		Find(&list).Error
	return
}

func (d *DbDao) GetSubAccountListCountByParentAccountId(parentAccountId string) (count int64, err error) {
	err = d.db.Model(tables.TableAccountInfo{}).Where("parent_account_id=?", parentAccountId).Count(&count).Error
	return
}

func (d *DbDao) GetSubAccByParentAccountIdOfAddress(parentAccountId, subAccountId, address string, verifyType uint) (count int64, err error) {
	var queryField string
	var queryValue string
	if subAccountId != "" {
		queryField = "account_id"
		queryValue = subAccountId
	} else {
		queryField = "parent_account_id"
		queryValue = parentAccountId
	}
	if verifyType == 0 {
		err = d.db.Model(tables.TableAccountInfo{}).Where(queryField+" =? and owner=? ", queryValue, address).Count(&count).Error
		return
	} else {
		err = d.db.Model(tables.TableAccountInfo{}).Where(queryField+" =? and manager=? ", queryValue, address).Count(&count).Error
		return
	}
}

func (d *DbDao) DelSubAccounts(subAccIds []string) error {
	if len(subAccIds) == 0 {
		return nil
	}
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("account_id IN(?)", subAccIds).
			Delete(&tables.TableAccountInfo{}).Error; err != nil {
			return err
		}
		if err := tx.Where("account_id IN(?)", subAccIds).
			Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *DbDao) UpdateAccounts(accounts []map[string]interface{}) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		for _, account := range accounts {
			accId, ok := account["account_id"]
			if !ok {
				return errors.New("account_id is not exist")
			}
			action := account["action"]
			delete(account, "action")

			if err := tx.Model(&tables.TableAccountInfo{}).Where("account_id=?", accId).
				Updates(account).Error; err != nil {
				return err
			}

			if action == common.SubActionFullfillApproval {
				if err := tx.Where("account_id = ?", accId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (d *DbDao) GetAccountByAccIds(accIds []string) (list []*tables.TableAccountInfo, err error) {
	if len(accIds) == 0 {
		return
	}
	err = d.db.Where("account_id in (?)", accIds).Find(&list).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = nil
	}
	return
}

func (d *DbDao) BidExpiredAccountAuction(accountInfo tables.TableAccountInfo, recordsInfos []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		//update account_info
		if err := tx.Select("status", "expired_at", "registered_at", "block_number", "outpoint", "owner_chain_type", "owner", "owner_algorithm_id", "owner_sub_aid", "manager_chain_type", "manager", "manager_algorithm_id", "manager_sub_aid").
			Where("account_id = ?", accountInfo.AccountId).
			Updates(accountInfo).Error; err != nil {
			return err
		}
		//delete record
		if err := tx.Where("account_id=?", accountInfo.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}
		//default record
		if len(recordsInfos) > 0 {
			if err := tx.Create(&recordsInfos).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (d *DbDao) AccountUpgrade(accountInfo tables.TableAccountInfo, didCellInfo tables.TableDidCellInfo, recordsInfos []tables.TableRecordsInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Select("status", "block_number", "outpoint").
			Where("account_id = ?", accountInfo.AccountId).
			Updates(accountInfo).Error; err != nil {
			return err
		}

		if err := tx.Where("account_id = ?", didCellInfo.AccountId).Delete(&tables.TableRecordsInfo{}).Error; err != nil {
			return err
		}

		if len(recordsInfos) > 0 {
			if err := tx.Create(&recordsInfos).Error; err != nil {
				return err
			}
		}
		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"args", "account", "expired_at",
				"created_at", "updated_at",
			}),
		}).Create(&didCellInfo).Error; err != nil {
			return err
		}

		return nil
	})
}
