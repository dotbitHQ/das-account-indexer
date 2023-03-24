package block_parser

import (
	"context"
	"das-account-indexer/config"
	"das-account-indexer/dao"
	"das-account-indexer/notify"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/witness"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"sync"
	"sync/atomic"
	"time"
)

var (
	log                 = mylog.NewLogger("block_parser", mylog.LevelDebug)
	IsLatestBlockNumber = false
	CurrentBlockNumber  = uint64(0)
)

type BlockParser struct {
	DasCore              *core.DasCore
	MapTransactionHandle map[common.DasAction]FuncTransactionHandle
	CurrentBlockNumber   uint64
	DbDao                *dao.DbDao
	ConcurrencyNum       uint64
	ConfirmNum           uint64
	Ctx                  context.Context
	Cancel               context.CancelFunc
	Wg                   *sync.WaitGroup

	errCountHandle int
}

func (b *BlockParser) initCurrentBlockNumber() error {
	if block, err := b.DbDao.FindCurrentBlockInfo(); err != nil {
		return err
	} else if block.Id > 0 {
		b.CurrentBlockNumber = block.BlockNumber
		CurrentBlockNumber = block.BlockNumber
	}
	return nil
}

func (b *BlockParser) getTipBlockNumber() (uint64, error) {
	if blockNumber, err := b.DasCore.Client().GetTipBlockNumber(b.Ctx); err != nil {
		return 0, fmt.Errorf("GetTipBlockNumber err:%s", err.Error())
	} else {
		return blockNumber, nil
	}
}

func (b *BlockParser) RunParser() error {
	b.registerTransactionHandle()
	if err := b.initCurrentBlockNumber(); err != nil {
		return fmt.Errorf("initCurrentBlockNumber err: %s", err.Error())
	}

	atomic.AddUint64(&b.CurrentBlockNumber, 1)
	b.Wg.Add(1)
	go func() {
		for {
			select {
			default:
				// get the new height and compare with current height
				latestBlockNumber, err := b.getTipBlockNumber()
				if err != nil {
					log.Error("get latest block number err:", err.Error())
					time.Sleep(time.Second)
				} else {
					// async
					if b.ConcurrencyNum > 1 && b.CurrentBlockNumber < (latestBlockNumber-b.ConfirmNum-b.ConcurrencyNum) {
						nowTime := time.Now()
						if err = b.parserConcurrencyMode(); err != nil {
							log.Error("parserConcurrencyMode err:", err.Error(), b.CurrentBlockNumber)
						}
						log.Info("parserConcurrencyMode time:", time.Since(nowTime).Seconds())
						IsLatestBlockNumber = false
						time.Sleep(time.Millisecond * 300)
					} else if b.CurrentBlockNumber < (latestBlockNumber - b.ConfirmNum) { // check rollback
						nowTime := time.Now()
						if err = b.parserSubMode(); err != nil {
							log.Error("parserSubMode err:", err.Error(), b.CurrentBlockNumber)
						}
						log.Info("parserSubMode time:", time.Since(nowTime).Seconds())
						IsLatestBlockNumber = true
						time.Sleep(time.Second * 1)
					} else {
						log.Info("RunParser:", IsLatestBlockNumber, b.CurrentBlockNumber, latestBlockNumber)
						IsLatestBlockNumber = true
						time.Sleep(time.Second * 10)
					}
					CurrentBlockNumber = b.CurrentBlockNumber
				}
			case <-b.Ctx.Done():
				b.Wg.Done()
				return
			}
		}
	}()
	return nil
}

func (b *BlockParser) parsingBlockData(block *types.Block) error {
	if err := b.checkContractVersion(); err != nil {
		return err
	}
	for _, tx := range block.Transactions {
		txHash := tx.Hash.Hex()
		blockNumber := block.Header.Number
		blockTimestamp := block.Header.Timestamp
		log.Info("parsingBlockData txHash:", txHash)

		if builder, err := witness.ActionDataBuilderFromTx(tx); err != nil {
			//log.Warn("ActionDataBuilderFromTx err:", err.Error())
		} else {
			// transaction parse by action
			req := FuncTransactionHandleReq{
				Tx:             tx,
				TxHash:         txHash,
				BlockNumber:    blockNumber,
				BlockTimestamp: blockTimestamp,
				Action:         builder.Action,
			}
			handle, ok := b.MapTransactionHandle[builder.Action]
			if !ok {
				log.Info("other handle:", txHash, builder.Action)
				handle = b.ActionUpdateAccountInfo
			}
			resp := handle(&req)
			if resp.Err != nil {
				log.Error("action handle resp:", builder.Action, blockNumber, txHash, resp.Err.Error())

				b.errCountHandle++
				if b.errCountHandle < 100 {
					// notify
					msg := "> Transaction hash：%s\n> Action：%s\n> Timestamp：%s\n> Error message：%s"
					msg = fmt.Sprintf(msg, txHash, builder.Action, time.Now().Format("2006-01-02 15:04:05"), resp.Err.Error())
					err = notify.SendLarkTextNotify(config.Cfg.Notice.LarkErrUrl, "DasAccountIndexer BlockParser", msg)
					if err != nil {
						log.Error("SendLarkTextNotify err:", err.Error())
					}
				}

				return resp.Err
			}
		}
	}
	b.errCountHandle = 0
	return nil
}

// subscribe mode
func (b *BlockParser) parserSubMode() error {
	log.Info("parserSubMode:", b.CurrentBlockNumber)
	block, err := b.DasCore.Client().GetBlockByNumber(b.Ctx, b.CurrentBlockNumber)
	if err != nil {
		return fmt.Errorf("GetBlockByNumber err: %s", err.Error())
	} else {
		blockHash := block.Header.Hash.Hex()
		parentHash := block.Header.ParentHash.Hex()
		log.Info("parserSubMode:", b.CurrentBlockNumber, blockHash, parentHash)
		// check rollback
		if fork, err := b.checkFork(parentHash); err != nil {
			return fmt.Errorf("checkFork err: %s", err.Error())
		} else if fork {
			log.Warn("checkFork is true:", b.CurrentBlockNumber, blockHash, parentHash)
			atomic.AddUint64(&b.CurrentBlockNumber, ^uint64(0))
		} else if err = b.parsingBlockData(block); err != nil {
			return fmt.Errorf("parsingBlockData err: %s", err.Error())
		} else {
			if err = b.DbDao.CreateBlockInfo(b.CurrentBlockNumber, blockHash, parentHash); err != nil {
				return fmt.Errorf("CreateBlockInfo err: %s", err.Error())
			} else {
				atomic.AddUint64(&b.CurrentBlockNumber, 1)
			}
			if err = b.DbDao.DeleteBlockInfo(b.CurrentBlockNumber - 20); err != nil {
				return fmt.Errorf("DeleteBlockInfo err: %s", err.Error())
			}
		}
	}
	return nil
}

// rollback checking
func (b *BlockParser) checkFork(parentHash string) (bool, error) {
	block, err := b.DbDao.FindBlockInfoByBlockNumber(b.CurrentBlockNumber - 1)
	if err != nil {
		return false, fmt.Errorf("FindBlockInfoByBlockNumber err: %s", err.Error())
	}
	if block.Id == 0 {
		return false, nil
	} else if block.BlockHash != parentHash {
		log.Warn("CheckFork:", b.CurrentBlockNumber, parentHash, block.BlockHash)
		return true, nil
	}
	return false, nil
}

func (b *BlockParser) parserConcurrencyMode() error {
	log.Info("parserConcurrencyMode:", b.CurrentBlockNumber, b.ConcurrencyNum)
	for i := uint64(0); i < b.ConcurrencyNum; i++ {
		block, err := b.DasCore.Client().GetBlockByNumber(b.Ctx, b.CurrentBlockNumber)
		if err != nil {
			return fmt.Errorf("GetBlockByNumber err: %s [%d]", err.Error(), b.CurrentBlockNumber)
		}
		blockHash := block.Header.Hash.Hex()
		parentHash := block.Header.ParentHash.Hex()
		log.Info("parserConcurrencyMode:", b.CurrentBlockNumber, blockHash, parentHash)

		if err = b.parsingBlockData(block); err != nil {
			return fmt.Errorf("parsingBlockData err: %s", err.Error())
		} else {
			if err = b.DbDao.CreateBlockInfo(b.CurrentBlockNumber, blockHash, parentHash); err != nil {
				return fmt.Errorf("CreateBlockInfo err: %s", err.Error())
			} else {
				atomic.AddUint64(&b.CurrentBlockNumber, 1)
			}
		}
	}
	if err := b.DbDao.DeleteBlockInfo(b.CurrentBlockNumber - 20); err != nil {
		return fmt.Errorf("DeleteBlockInfo err: %s", err.Error())
	}
	return nil
}

var contractNames = []common.DasContractName{
	//common.DasContractNameApplyRegisterCellType,
	//common.DasContractNamePreAccountCellType,
	//common.DasContractNameProposalCellType,
	common.DasContractNameConfigCellType,
	common.DasContractNameAccountCellType,
	common.DasContractNameAccountSaleCellType,
	common.DASContractNameSubAccountCellType,
	common.DASContractNameOfferCellType,
	//common.DasContractNameBalanceCellType,
	//common.DasContractNameIncomeCellType,
	common.DasContractNameReverseRecordCellType,
	//common.DASContractNameEip712LibCellType,
	common.DasContractNameReverseRecordRootCellType,
}

func (b *BlockParser) checkContractVersion() error {
	for _, v := range contractNames {
		defaultVersion, chainVersion, err := b.DasCore.CheckContractVersion(v)
		if err != nil {
			if err == core.ErrContractMajorVersionDiff {
				log.Errorf("contract[%s] version diff, chain[%s], service[%s].", v, chainVersion, defaultVersion)
				log.Error("Please update the service. [https://github.com/dotbitHQ/das-account-indexer]")
				if b.Cancel != nil && !config.Cfg.Server.NotExit {
					b.Cancel()
				}
				return err
			}
			return fmt.Errorf("CheckContractVersion err: %s", err.Error())
		}
	}
	return nil
}
