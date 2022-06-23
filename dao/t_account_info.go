package dao

import (
	"das-account-indexer/tables"
	"github.com/DeAccountSystems/das-lib/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (d *DbDao) FindAccountListByAddress(chainType common.ChainType, address string) (list []tables.TableAccountInfo, err error) {
	err = d.db.Where(" owner_chain_type=? AND owner=? ", chainType, address).Find(&list).Error
	return
}

func (d *DbDao) FindAccountNameListByAddress(chainType common.ChainType, address string) (list []tables.TableAccountInfo, err error) {
	err = d.db.Select("account").Where(" owner_chain_type=? AND owner=? ", chainType, address).Find(&list).Error
	return
}

func (d *DbDao) EnableSubAccount(accountInfo tables.TableAccountInfo) error {
	return d.db.Select("block_number", "block_timestamp", "outpoint", "enable_sub_account", "renew_sub_account_price").
		Where("account_id = ?", accountInfo.AccountId).Updates(accountInfo).Error
}

func (d *DbDao) CreateSubAccount(subAccountIds []string, accountInfos []tables.TableAccountInfo, accountInfo tables.TableAccountInfo) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if len(subAccountIds) > 0 {
			if err := tx.Where(" account_id IN(?) ", subAccountIds).
				Delete(&tables.TableRecordsInfo{}).Error; err != nil {
				return err
			}
		}
		if len(accountInfos) > 0 {
			if err := tx.Clauses(clause.OnConflict{
				DoUpdates: clause.AssignmentColumns([]string{
					"block_number", "block_timestamp", "outpoint",
					"owner_chain_type", "owner", "owner_algorithm_id",
					"manager_chain_type", "manager", "manager_algorithm_id",
					"registered_at", "expired_at", "status",
					"enable_sub_account", "renew_sub_account_price", "nonce",
				}),
			}).Create(&accountInfos).Error; err != nil {
				return err
			}
		}

		if err := tx.Select("block_number", "block_timestamp", "outpoint").
			Where("account_id = ?", accountInfo.AccountId).
			Updates(accountInfo).Error; err != nil {
			return err
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

func (d *DbDao) RecycleExpiredAccount(previousAccountId, previousNextAccountId, accountId string, enableSubAccount uint8) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&tables.TableAccountInfo{}).Where("account_id=?", previousAccountId).
			Update("next_account_id", previousNextAccountId).Error; err != nil {
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
