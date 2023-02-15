package block_parser

import (
	"das-account-indexer/tables"
	"das_database/dao"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/witness"
	"gorm.io/gorm"
)

func (b *BlockParser) ActionReverseRecordRoot(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameReverseRecordRootCellType); err != nil {
		resp.Err = fmt.Errorf("isCurrentVersion err: %s", err)
		return
	} else if !isCV {
		return
	}
	log.Info("ActionReverseRecordRoot:", req.BlockNumber, req.TxHash)

	smtBuilder := witness.NewReverseSmtBuilder()
	txReverseSmtRecord, err := smtBuilder.FromTx(req.Tx)
	if err != nil {
		resp.Err = err
		return
	}

	if err := b.DbDao.Transaction(func(tx *gorm.DB) error {
		for idx, v := range txReverseSmtRecord {
			outpoint := common.OutPoint2String(req.TxHash, uint(idx))
			accountId := common.Bytes2Hex(common.GetAccountIdByAccount(v.NextAccount))
			algorithmId := common.DasAlgorithmId(v.SignType)
			reverseInfo := &dao.TableReverseInfo{
				BlockNumber:    req.BlockNumber,
				BlockTimestamp: req.BlockTimestamp,
				Outpoint:       outpoint,
				AlgorithmId:    algorithmId,
				ChainType:      algorithmId.ToChainType(),
				Address:        v.Address,
				Account:        v.NextAccount,
				AccountId:      accountId,
				ReverseType:    dao.ReverseTypeSmt,
			}

			if v.PrevAccount != "" {
				if err := tx.Where("address=? and reverse_type=?", v.Address, dao.ReverseTypeSmt).Delete(&tables.TableReverseInfo{}).Error; err != nil {
					return err
				}
			}
			if v.Action == witness.ReverseSmtRecordActionUpdate {
				if err := tx.Create(reverseInfo).Error; err != nil {
					return err
				}
			}
		}
		return nil
	}); err != nil {
		resp.Err = err
		return
	}
	return
}
