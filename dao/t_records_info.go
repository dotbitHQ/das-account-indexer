package dao

import "das-account-indexer/tables"

func (d *DbDao) FindAccountRecords(account string) (list []tables.TableRecordsInfo, err error) {
	err = d.db.Where(" account=? ", account).Find(&list).Error
	return
}

func (d *DbDao) FindRecordsByAccounts(accounts []string) (list []tables.TableRecordsInfo, err error) {
	err = d.db.Where(" account IN(?) ", accounts).Find(&list).Error
	return
}

func (d *DbDao) FindRecordByAccountAddressValue(account, value string) (r tables.TableRecordsInfo, err error) {
	err = d.db.Where(" `account`=? AND `type`='address' AND `value`=? ", account, value).Find(&r).Limit(1).Error
	return
}
