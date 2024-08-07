package block_parser

import (
	"das-account-indexer/tables"
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
			address := common.FormatAddressPayload(v.Address, algorithmId)
			p2shP2wpkh, err := v.GetP2SHP2WPKH(b.DasCore.NetType())
			if err != nil {
				log.Error("GetP2SHP2WPKH err: %s", err.Error())
			}
			p2tr, err := v.GetP2TR(b.DasCore.NetType())
			if err != nil {
				log.Error("GetP2TR err: %s", err.Error())
			}
			reverseInfo := &tables.TableReverseInfo{
				BlockNumber:    req.BlockNumber,
				BlockTimestamp: req.BlockTimestamp,
				Outpoint:       outpoint,
				AlgorithmId:    algorithmId,
				//SubAlgorithmId:
				ChainType:   algorithmId.ToChainType(),
				Address:     address,
				Account:     v.NextAccount,
				ReverseType: tables.ReverseTypeSmt,
				P2shP2wpkh:  p2shP2wpkh,
				P2tr:        p2tr,
			}
			if v.PrevAccount != "" {
				if err := tx.Where("address=? and reverse_type=?", address, tables.ReverseTypeSmt).Delete(&tables.TableReverseInfo{}).Error; err != nil {
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
