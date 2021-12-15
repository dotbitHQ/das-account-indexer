package dao

import (
	"das-account-indexer/tables"
	"github.com/DeAccountSystems/das-lib/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (d *DbDao) CreateReverseInfo(reverse *tables.TableReverseInfo) error {
	if err := d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{
			"block_number", "block_timestamp", "outpoint",
			"algorithm_id", "chain_type", "address",
			"account", "capacity",
		}),
	}).Create(&reverse).Error; err != nil {
		return err
	}
	return nil
}

func (d *DbDao) UpdateReverseInfo(reverse *tables.TableReverseInfo, lastOutpoint string) error {
	return d.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.AssignmentColumns([]string{
				"block_number", "block_timestamp", "outpoint",
				"algorithm_id", "chain_type", "address",
				"account", "capacity",
			}),
		}).Create(&reverse).Error; err != nil {
			return err
		}
		if err := tx.Where(" outpoint=? ", lastOutpoint).Delete(&tables.TableReverseInfo{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func (d *DbDao) DeleteReverseInfo(outpoints []string) error {
	return d.db.Where(" outpoint IN (?) ", outpoints).Delete(&tables.TableReverseInfo{}).Error
}

func (d *DbDao) FindLatestReverseRecord(chainType common.ChainType, address string) (r tables.TableReverseInfo, err error) {
	err = d.db.Where(" chain_type=? AND address=? ", chainType, address).Order(" block_number DESC,outpoint DESC ").Limit(1).Find(&r).Error
	return
}
