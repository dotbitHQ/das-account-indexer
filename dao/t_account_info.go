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

func (d *DbDao) FindAccountInfoByAccountName(accountName string) (accountInfo tables.TableAccountInfo, err error) {
	err = d.db.Where(" account=? ", accountName).Find(&accountInfo).Error
	return
}

func (d *DbDao) FindAccountListByAddress(chainType common.ChainType, address string) (list []tables.TableAccountInfo, err error) {
	err = d.db.Where(" owner_chain_type=? AND owner=? ", chainType, address).Find(&list).Error
	return
}
