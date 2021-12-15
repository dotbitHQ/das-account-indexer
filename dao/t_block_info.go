package dao

import (
	"das-account-indexer/tables"
	"gorm.io/gorm/clause"
)

func (d *DbDao) CreateBlockInfo(blockNumber uint64, blockHash, parentHash string) error {
	return d.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"block_hash", "parent_hash"}),
	}).Create(&tables.TableBlockInfo{
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
		ParentHash:  parentHash,
	}).Error
}

func (d *DbDao) DeleteBlockInfo(blockNumber uint64) error {
	return d.db.Where("block_number < ?", blockNumber).Delete(&tables.TableBlockInfo{}).Error
}

func (d *DbDao) FindCurrentBlockInfo() (blockInfo tables.TableBlockInfo, err error) {
	err = d.db.Order("block_number DESC").Limit(1).Find(&blockInfo).Error
	return
}

func (d *DbDao) FindBlockInfoByBlockNumber(blockNumber uint64) (blockInfo tables.TableBlockInfo, err error) {
	err = d.db.Where("block_number = ?", blockNumber).Limit(1).Find(&blockInfo).Error
	return
}
