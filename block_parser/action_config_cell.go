package block_parser

import (
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
)

func (b *BlockParser) ActionConfigCell(req *FuncTransactionHandleReq) (resp FuncTransactionHandleResp) {
	if isCV, err := isCurrentVersionTx(req.Tx, common.DasContractNameConfigCellType); err != nil {
		resp.Err = fmt.Errorf("ActionConfigCell err: %s", err.Error())
		return
	} else if !isCV {
		return
	}
	log.Info("das tx:", req.Action, req.TxHash)

	if err := b.DasCore.AsyncDasConfigCell(); err != nil {
		resp.Err = fmt.Errorf("AsyncDasConfigCell err: %s", err.Error())
		return
	}
	return
}
