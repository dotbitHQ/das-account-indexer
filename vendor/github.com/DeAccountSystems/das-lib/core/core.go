package core

import (
	"context"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
	"golang.org/x/sync/syncmap"
	"sync"
)

var (
	log                      = mylog.NewLogger("das-core", mylog.LevelDebug)
	DasContractMap           syncmap.Map                               // map[contact name]{contract info}
	DasContractByTypeIdMap   = make(map[string]common.DasContractName) // map[contract type id]{contract name}
	DasConfigCellMap         syncmap.Map                               // map[ConfigCellTypeArgs]config cell info
	DasConfigCellByTxHashMap syncmap.Map                               // map[tx hash]{true}
	DasSoScriptMap           syncmap.Map                               // map[so script type]
)

type DasCore struct {
	client              rpc.Client
	ctx                 context.Context
	wg                  *sync.WaitGroup
	dasContractCodeHash string // contract code hash
	dasContractArgs     string // contract owner args
	thqCodeHash         string // time,height,quote cell code hash
	net                 common.DasNetType
	daf                 *DasAddressFormat
}

func NewDasCore(ctx context.Context, wg *sync.WaitGroup, opts ...DasCoreOption) *DasCore {
	var dc DasCore
	dc.ctx = ctx
	dc.wg = wg
	for _, opt := range opts {
		opt(&dc)
	}
	return &dc
}

func (d *DasCore) Client() rpc.Client {
	return d.client
}

func (d *DasCore) NetType() common.DasNetType {
	return d.net
}

func (d *DasCore) Daf() *DasAddressFormat {
	return d.daf
}

func (d *DasCore) GetDasLock() *types.Script {
	if d.net == common.DasNetTypeMainNet {
		return common.GetNormalLockScriptByMultiSig("0xc126635ece567c71c50f7482c5db80603852c306")
	} else {
		return common.GetNormalLockScript("0xefbf497f752ff7a655a8ec6f3c8f3feaaed6e410")
	}
}

func SetLogLevel(level int) {
	log = mylog.NewLogger("das-core", level)
}
