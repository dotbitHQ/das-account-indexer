package block_parser

import (
	"das-account-indexer/tables"
	"encoding/hex"
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

	txReverseSmtRecord := make([]*witness.ReverseSmtRecord, 0)
	if err := witness.ParseFromTx(req.Tx, common.ActionDataTypeReverseSmt, &txReverseSmtRecord); err != nil {
		resp.Err = err
		return
	}

	if err := b.DbDao.Transaction(func(tx *gorm.DB) error {
		for idx, v := range txReverseSmtRecord {
			outpoint := common.OutPoint2String(req.TxHash, uint(idx))
			algorithmId := common.DasAlgorithmId(v.SignType)
			address := hex.EncodeToString(v.Address)
			if algorithmId == common.DasAlgorithmIdTron {
				address = common.TronPreFix + address
			}
			reverseInfo := &tables.TableReverseInfo{
				BlockNumber:    req.BlockNumber,
				BlockTimestamp: req.BlockTimestamp,
				Outpoint:       outpoint,
				AlgorithmId:    algorithmId,
				ChainType:      algorithmId.ToChainType(),
				Address:        address,
				Account:        v.NextAccount,
				ReverseType:    tables.ReverseTypeSmt,
			}
			addresses := []string{address, common.HexPreFix + address}

			if v.PrevAccount != "" {
				if err := tx.Where("address in (?) and reverse_type=?", addresses, tables.ReverseTypeSmt).Delete(&tables.TableReverseInfo{}).Error; err != nil {
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
