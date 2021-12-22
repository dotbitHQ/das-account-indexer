package dao

import "das-account-indexer/tables"

func (d *DbDao) FindAccountRecordsByAccountId(accountId string) (list []tables.TableRecordsInfo, err error) {
	err = d.db.Where(" account_id=? ", accountId).Find(&list).Error
	return
}

func (d *DbDao) FindRecordsByAccountIds(accountIds []string) (list []tables.TableRecordsInfo, err error) {
	err = d.db.Where(" account_id IN(?) ", accountIds).Find(&list).Error
	return
}

func (d *DbDao) FindRecordByAccountIdAddressValue(accountId, value string) (r tables.TableRecordsInfo, err error) {
	err = d.db.Where(" `account_id`=? AND `type`='address' AND `value`=? ", accountId, value).Find(&r).Limit(1).Error
	return
}
