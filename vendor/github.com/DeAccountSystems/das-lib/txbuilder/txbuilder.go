package txbuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/DeAccountSystems/das-lib/sign"
	"github.com/nervosnetwork/ckb-sdk-go/rpc"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/mylog"
)

var log = mylog.NewLogger("txbuilder", mylog.LevelDebug)

type DasTxBuilder struct {
	*DasTxBuilderBase                                  // for base
	*DasTxBuilderTransaction                           // for tx
	DasMMJson                                          // for mmjson
	mapCellDep               map[string]*types.CellDep // for memory
}

func NewDasTxBuilderBase(ctx context.Context, dasCore *core.DasCore, handle sign.HandleSignCkbMessage, serverArgs string) *DasTxBuilderBase {
	var base DasTxBuilderBase
	base.ctx = ctx
	base.dasCore = dasCore
	base.handleServerSign = handle
	base.serverArgs = serverArgs
	return &base
}

func NewDasTxBuilderFromBase(base *DasTxBuilderBase, tx *DasTxBuilderTransaction) *DasTxBuilder {
	var b DasTxBuilder
	b.DasTxBuilderBase = base
	b.DasTxBuilderTransaction = tx
	if tx == nil {
		b.DasTxBuilderTransaction = &DasTxBuilderTransaction{}
		b.MapInputsCell = make(map[string]*types.CellWithStatus)
	}
	b.mapCellDep = make(map[string]*types.CellDep)
	return &b
}

type DasTxBuilderBase struct {
	ctx              context.Context
	dasCore          *core.DasCore
	handleServerSign sign.HandleSignCkbMessage
	serverArgs       string
}

type DasTxBuilderTransaction struct {
	Transaction     *types.Transaction               `json:"transaction"`
	MapInputsCell   map[string]*types.CellWithStatus `json:"map_inputs_cell"`
	ServerSignGroup []int                            `json:"server_sign_group"`
}

type DasMMJson struct {
	account            string
	accountDasLockArgs []byte
	salePrice          uint64
	offers             int // cancel offer count
}

type BuildTransactionParams struct {
	CellDeps    []*types.CellDep    `json:"cell_deps"`
	Inputs      []*types.CellInput  `json:"inputs"`
	Outputs     []*types.CellOutput `json:"outputs"`
	OutputsData [][]byte            `json:"outputs_data"`
	Witnesses   [][]byte            `json:"witnesses"`
}

func (d *DasTxBuilder) BuildTransaction(p *BuildTransactionParams) error {
	err := d.newTx()
	if err != nil {
		return fmt.Errorf("newBaseTx err: %s", err.Error())
	}

	err = d.addInputsForTx(p.Inputs)
	if err != nil {
		return fmt.Errorf("addInputsForBaseTx err: %s", err.Error())
	}

	err = d.addOutputsForTx(p.Outputs, p.OutputsData)
	if err != nil {
		return fmt.Errorf("addOutputsForBaseTx err: %s", err.Error())
	}

	d.Transaction.Witnesses = append(d.Transaction.Witnesses, p.Witnesses...)

	if err := d.addMapCellDepWitnessForBaseTx(p.CellDeps); err != nil {
		return fmt.Errorf("addMapCellDepWitnessForBaseTx err: %s", err.Error())
	}

	return nil
}

func (d *DasTxBuilder) TxString() string {
	txStr, _ := rpc.TransactionString(d.Transaction)
	return txStr
}

func (d *DasTxBuilder) GetDasTxBuilderTransactionString() string {
	bys, err := json.Marshal(d.DasTxBuilderTransaction)
	if err != nil {
		return ""
	}
	return string(bys)
}
