package rpc

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/nervosnetwork/ckb-sdk-go/types"
)

type epoch struct {
	CompactTarget hexutil.Uint64 `json:"compact_target"`
	Length        hexutil.Uint64 `json:"length"`
	Number        hexutil.Uint64 `json:"number"`
	StartNumber   hexutil.Uint64 `json:"start_number"`
}

type header struct {
	CompactTarget    hexutil.Uint   `json:"compact_target"`
	Dao              types.Hash     `json:"dao"`
	Epoch            hexutil.Uint64 `json:"epoch"`
	Hash             types.Hash     `json:"hash"`
	Nonce            hexutil.Big    `json:"nonce"`
	Number           hexutil.Uint64 `json:"number"`
	ParentHash       types.Hash     `json:"parent_hash"`
	ProposalsHash    types.Hash     `json:"proposals_hash"`
	Timestamp        hexutil.Uint64 `json:"timestamp"`
	TransactionsRoot types.Hash     `json:"transactions_root"`
	UnclesHash       types.Hash     `json:"uncles_hash"`
	Version          hexutil.Uint   `json:"version"`
}

type outPoint struct {
	TxHash types.Hash   `json:"tx_hash"`
	Index  hexutil.Uint `json:"index"`
}

type cellDep struct {
	OutPoint outPoint      `json:"out_point"`
	DepType  types.DepType `json:"dep_type"`
}

type cellInput struct {
	Since          hexutil.Uint64 `json:"since"`
	PreviousOutput outPoint       `json:"previous_output"`
}

type script struct {
	CodeHash types.Hash           `json:"code_hash"`
	HashType types.ScriptHashType `json:"hash_type"`
	Args     hexutil.Bytes        `json:"args"`
}

type cellOutput struct {
	Capacity hexutil.Uint64 `json:"capacity"`
	Lock     *script        `json:"lock"`
	Type     *script        `json:"type"`
}

type transaction struct {
	Version     hexutil.Uint    `json:"version"`
	Hash        types.Hash      `json:"hash"`
	CellDeps    []cellDep       `json:"cell_deps"`
	HeaderDeps  []types.Hash    `json:"header_deps"`
	Inputs      []cellInput     `json:"inputs"`
	Outputs     []cellOutput    `json:"outputs"`
	OutputsData []hexutil.Bytes `json:"outputs_data"`
	Witnesses   []hexutil.Bytes `json:"witnesses"`
}

type inTransaction struct {
	Version     hexutil.Uint    `json:"version"`
	CellDeps    []cellDep       `json:"cell_deps"`
	HeaderDeps  []types.Hash    `json:"header_deps"`
	Inputs      []cellInput     `json:"inputs"`
	Outputs     []cellOutput    `json:"outputs"`
	OutputsData []hexutil.Bytes `json:"outputs_data"`
	Witnesses   []hexutil.Bytes `json:"witnesses"`
}

type uncleBlock struct {
	Header    header   `json:"header"`
	Proposals []string `json:"proposals"`
}

type block struct {
	Header       header        `json:"header"`
	Proposals    []string      `json:"proposals"`
	Transactions []transaction `json:"transactions"`
	Uncles       []uncleBlock  `json:"uncles"`
}

type cell struct {
	BlockHash     types.Hash     `json:"block_hash"`
	Capacity      hexutil.Uint64 `json:"capacity"`
	Lock          *script        `json:"lock"`
	OutPoint      *outPoint      `json:"out_point"`
	Type          *script        `json:"type"`
	Cellbase      bool           `json:"cellbase,omitempty"`
	OutputDataLen hexutil.Uint64 `json:"output_data_len,omitempty"`
}

type cellData struct {
	Content hexutil.Bytes `json:"content"`
	Hash    types.Hash    `json:"hash"`
}

type cellInfo struct {
	Data   *cellData  `json:"data"`
	Output cellOutput `json:"output"`
}

type cellWithStatus struct {
	Cell   cellInfo `json:"cell"`
	Status string   `json:"status"`
}

type transactionWithStatus struct {
	Transaction transaction `json:"transaction"`
	TxStatus    struct {
		BlockHash *types.Hash             `json:"block_hash"`
		Status    types.TransactionStatus `json:"status"`
	} `json:"tx_status"`
}

type blockReward struct {
	Primary        hexutil.Big `json:"primary"`
	ProposalReward hexutil.Big `json:"proposal_reward"`
	Secondary      hexutil.Big `json:"secondary"`
	Total          hexutil.Big `json:"total"`
	TxFee          hexutil.Big `json:"tx_fee"`
}

type dryRunTransactionResult struct {
	Cycles hexutil.Uint64 `json:"cycles"`
}

type estimateFeeRateResult struct {
	FeeRate hexutil.Uint64 `json:"fee_rate"`
}

type lockHashIndexState struct {
	BlockHash   types.Hash     `json:"block_hash"`
	BlockNumber hexutil.Uint64 `json:"block_number"`
	LockHash    types.Hash     `json:"lock_hash"`
}

type transactionPoint struct {
	BlockNumber hexutil.Uint64 `json:"block_number"`
	Index       hexutil.Uint   `json:"index"`
	TxHash      types.Hash     `json:"tx_hash"`
}

type liveCell struct {
	CellOutput cellOutput       `json:"cell_output"`
	CreatedBy  transactionPoint `json:"created_by"`
}

type cellTransaction struct {
	ConsumedBy *transactionPoint `json:"consumed_by"`
	CreatedBy  *transactionPoint `json:"created_by"`
}

type nodeAddress struct {
	Address string         `json:"address"`
	Score   hexutil.Uint64 `json:"score"`
}

type node struct {
	Addresses  []*nodeAddress `json:"addresses"`
	IsOutbound bool           `json:"is_outbound"`
	NodeId     string         `json:"node_id"`
	Version    string         `json:"version"`
}

type bannedAddress struct {
	Address   string         `json:"address"`
	BanReason string         `json:"ban_reason"`
	BanUntil  hexutil.Uint64 `json:"ban_until"`
	CreatedAt hexutil.Uint64 `json:"created_at"`
}

type txPoolInfo struct {
	LastTxsUpdatedAt hexutil.Uint64 `json:"last_txs_updated_at"`
	Orphan           hexutil.Uint64 `json:"orphan"`
	Pending          hexutil.Uint64 `json:"pending"`
	Proposed         hexutil.Uint64 `json:"proposed"`
	TotalTxCycles    hexutil.Uint64 `json:"total_tx_cycles"`
	TotalTxSize      hexutil.Uint64 `json:"total_tx_size"`
}

type alertMessage struct {
	Id          string         `json:"id"`
	Message     string         `json:"message"`
	NoticeUntil hexutil.Uint64 `json:"notice_until"`
	Priority    string         `json:"priority"`
}

type blockchainInfo struct {
	Alerts                 []*alertMessage `json:"alerts"`
	Chain                  string          `json:"chain"`
	Difficulty             hexutil.Big     `json:"difficulty"`
	Epoch                  hexutil.Uint64  `json:"epoch"`
	IsInitialBlockDownload bool            `json:"is_initial_block_download"`
	MedianTime             hexutil.Uint64  `json:"median_time"`
}

type blockEconomicState struct {
	Issuance    blockIssuance `json:"issuance"`
	MinerReward minerReward   `json:"miner_reward"`
	TxsFee      hexutil.Big   `json:"txs_fee"`
	FinalizedAt types.Hash    `json:"finalized_at"`
}

type blockIssuance struct {
	Primary   hexutil.Big `json:"primary"`
	Secondary hexutil.Big `json:"secondary"`
}

type minerReward struct {
	Primary   hexutil.Big `json:"primary"`
	Secondary hexutil.Big `json:"secondary"`
	Committed hexutil.Big `json:"committed"`
	Proposal  hexutil.Big `json:"proposal"`
}

type rationalU256 struct {
	Denom hexutil.Big `json:"denom"`
	Numer hexutil.Big `json:"numer"`
}

type proposalWindow struct {
	Closest  hexutil.Uint64 `json:"closest"`
	Farthest hexutil.Uint64 `json:"farthest"`
}

type consensus struct {
	Id                                   string         `json:"ID"`
	GenesisHash                          types.Hash     `json:"genesis_hash"`
	DaoTypeHash                          types.Hash     `json:"dao_type_hash"`
	Secp256k1Blake160SighashAllTypeHash  types.Hash     `json:"secp256k1_blake160_sighash_all_type_hash"`
	Secp256k1Blake160MultisigAllTypeHash types.Hash     `json:"secp256k1_blake160_multisig_all_type_hash"`
	InitialPrimaryEpochReward            hexutil.Uint64 `json:"initial_primary_epoch_reward"`
	SecondaryEpochReward                 hexutil.Uint64 `json:"secondary_epoch_reward"`
	MaxUnclesNum                         hexutil.Uint64 `json:"max_uncles_num"`
	OrphanRateTarget                     rationalU256   `json:"orphan_rate_target"`
	EpochDurationTarget                  hexutil.Uint64 `json:"epoch_duration_target"`
	TxProposalWindow                     proposalWindow `json:"tx_proposal_window"`
	ProposerRewardRatio                  rationalU256   `json:"proposer_reward_ratio"`
	CellbaseMaturity                     hexutil.Uint64 `json:"cellbase_maturity"`
	MedianTimeBlockCount                 hexutil.Uint64 `json:"median_time_block_count"`
	MaxBlockCycles                       hexutil.Uint64 `json:"max_block_cycles"`
	BlockVersion                         hexutil.Uint   `json:"block_version"`
	TxVersion                            hexutil.Uint   `json:"tx_version"`
	TypeIdCodeHash                       types.Hash     `json:"type_id_code_hash"`
	MaxBlockProposalsLimit               hexutil.Uint64 `json:"max_block_proposals_limit"`
	PrimaryEpochRewardHalvingInterval    hexutil.Uint64 `json:"primary_epoch_reward_halving_interval"`
	PermanentDifficultyInDummy           bool           `json:"permanent_difficulty_in_dummy"`
}

type syncState struct {
	Ibd                     bool           `json:"ibd"`
	BestKnownBlockNumber    hexutil.Uint64 `json:"best_known_block_number"`
	BestKnownBlockTimestamp hexutil.Uint64 `json:"best_known_block_timestamp"`
	OrphanBlocksCount       hexutil.Uint64 `json:"orphan_blocks_count"`
	InflightBlocksCount     hexutil.Uint64 `json:"inflight_blocks_count"`
	FastTime                hexutil.Uint64 `json:"fast_time"`
	LowTime                 hexutil.Uint64 `json:"low_time"`
	NormalTime              hexutil.Uint64 `json:"normal_time"`
}

func toSyncState(syncState syncState) *types.SyncState {
	return &types.SyncState{
		syncState.Ibd,
		uint64(syncState.BestKnownBlockNumber),
		uint64(syncState.BestKnownBlockTimestamp),
		uint64(syncState.OrphanBlocksCount),
		uint64(syncState.InflightBlocksCount),
		uint64(syncState.FastTime),
		uint64(syncState.LowTime),
		uint64(syncState.NormalTime),
	}
}

type transactionProof struct {
	Proof         proof      `json:"proof"`
	BlockHash     types.Hash `json:"block_hash"`
	WitnessesRoot types.Hash `json:"witnesses_root"`
}

type proof struct {
	Indices []hexutil.Uint `json:"indices"`
	Lemmas  []types.Hash   `json:"lemmas"`
}

func toReqTransactionProof(reqTransactionProof types.TransactionProof) transactionProof {

	return transactionProof{
		toReqProof(reqTransactionProof.Proof),
		reqTransactionProof.BlockHash,
		reqTransactionProof.WitnessesRoot,
	}
}

func toReqProof(reqProof *types.Proof) proof {
	indices := make([]hexutil.Uint, len(reqProof.Indices))
	for i, value := range reqProof.Indices {
		indices[i] = hexutil.Uint(value)
	}
	return proof{
		indices,
		reqProof.Lemmas,
	}
}

func toTransactionProof(transactionProof transactionProof) *types.TransactionProof {
	return &types.TransactionProof{
		toProof(transactionProof.Proof),
		transactionProof.BlockHash,
		transactionProof.WitnessesRoot,
	}
}

func toProof(proof proof) *types.Proof {
	return &types.Proof{
		toUints(proof.Indices),
		proof.Lemmas,
	}
}

func toHeader(head header) *types.Header {
	return &types.Header{
		CompactTarget:    uint(head.CompactTarget),
		Dao:              head.Dao,
		Epoch:            uint64(head.Epoch),
		Hash:             head.Hash,
		Nonce:            (*big.Int)(&head.Nonce),
		Number:           uint64(head.Number),
		ParentHash:       head.ParentHash,
		ProposalsHash:    head.ProposalsHash,
		Timestamp:        uint64(head.Timestamp),
		TransactionsRoot: head.TransactionsRoot,
		UnclesHash:       head.UnclesHash,
		Version:          uint(head.Version),
	}
}

func toTransaction(tx transaction) *types.Transaction {
	return &types.Transaction{
		Version:     uint(tx.Version),
		Hash:        tx.Hash,
		CellDeps:    toCellDeps(tx.CellDeps),
		HeaderDeps:  tx.HeaderDeps,
		Inputs:      toInputs(tx.Inputs),
		Outputs:     toOutputs(tx.Outputs),
		OutputsData: toBytesArray(tx.OutputsData),
		Witnesses:   toBytesArray(tx.Witnesses),
	}
}

func toTransactions(transactions []transaction) []*types.Transaction {
	result := make([]*types.Transaction, len(transactions))
	for i := 0; i < len(transactions); i++ {
		result[i] = toTransaction(transactions[i])
	}
	return result
}

func toBytesArray(bytes []hexutil.Bytes) [][]byte {
	result := make([][]byte, len(bytes))
	for i, data := range bytes {
		result[i] = data
	}
	return result
}

func toOutputs(outputs []cellOutput) []*types.CellOutput {
	result := make([]*types.CellOutput, len(outputs))
	for i := 0; i < len(outputs); i++ {
		output := outputs[i]
		result[i] = &types.CellOutput{
			Capacity: uint64(output.Capacity),
			Lock: &types.Script{
				CodeHash: output.Lock.CodeHash,
				HashType: output.Lock.HashType,
				Args:     output.Lock.Args,
			},
		}
		if output.Type != nil {
			result[i].Type = &types.Script{
				CodeHash: output.Type.CodeHash,
				HashType: output.Type.HashType,
				Args:     output.Type.Args,
			}
		}
	}
	return result
}

func toInputs(inputs []cellInput) []*types.CellInput {
	result := make([]*types.CellInput, len(inputs))
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		result[i] = &types.CellInput{
			Since: uint64(input.Since),
			PreviousOutput: &types.OutPoint{
				TxHash: input.PreviousOutput.TxHash,
				Index:  uint(input.PreviousOutput.Index),
			},
		}
	}
	return result
}

func toCellDeps(deps []cellDep) []*types.CellDep {
	result := make([]*types.CellDep, len(deps))
	for i := 0; i < len(deps); i++ {
		dep := deps[i]
		result[i] = &types.CellDep{
			OutPoint: &types.OutPoint{
				TxHash: dep.OutPoint.TxHash,
				Index:  uint(dep.OutPoint.Index),
			},
			DepType: dep.DepType,
		}
	}
	return result
}

func toUints(uints []hexutil.Uint) []uint {
	result := make([]uint, len(uints))
	for i, value := range uints {
		result[i] = uint(value)
	}
	return result
}

func toUncles(uncles []uncleBlock) []*types.UncleBlock {
	result := make([]*types.UncleBlock, len(uncles))
	for i := 0; i < len(uncles); i++ {
		block := uncles[i]
		result[i] = &types.UncleBlock{
			Header:    toHeader(block.Header),
			Proposals: block.Proposals,
		}
	}
	return result
}

func toCells(cells []cell) []*types.Cell {
	result := make([]*types.Cell, len(cells))
	for i := 0; i < len(cells); i++ {
		cell := cells[i]
		result[i] = &types.Cell{
			BlockHash: cell.BlockHash,
			Capacity:  uint64(cell.Capacity),
			Lock: &types.Script{
				CodeHash: cell.Lock.CodeHash,
				HashType: cell.Lock.HashType,
				Args:     cell.Lock.Args,
			},
			OutPoint: &types.OutPoint{
				TxHash: cell.OutPoint.TxHash,
				Index:  uint(cell.OutPoint.Index),
			},
			Cellbase:      cell.Cellbase,
			OutputDataLen: uint64(cell.OutputDataLen),
		}
		if cell.Type != nil {
			result[i].Type = &types.Script{
				CodeHash: cell.Type.CodeHash,
				HashType: cell.Type.HashType,
				Args:     cell.Type.Args,
			}
		}
	}
	return result
}

func toCellWithStatus(status cellWithStatus) *types.CellWithStatus {
	result := &types.CellWithStatus{
		Cell: &types.CellInfo{
			Output: &types.CellOutput{
				Capacity: uint64(status.Cell.Output.Capacity),
			},
		},
		Status: status.Status,
	}

	if status.Cell.Output.Lock != nil {
		result.Cell.Output.Lock = &types.Script{
			CodeHash: status.Cell.Output.Lock.CodeHash,
			HashType: status.Cell.Output.Lock.HashType,
			Args:     status.Cell.Output.Lock.Args,
		}
	}

	if status.Cell.Data != nil {
		result.Cell.Data = &types.CellData{
			Content: status.Cell.Data.Content,
			Hash:    status.Cell.Data.Hash,
		}
	}

	if status.Cell.Output.Type != nil {
		result.Cell.Output.Type = &types.Script{
			CodeHash: status.Cell.Output.Type.CodeHash,
			HashType: status.Cell.Output.Type.HashType,
			Args:     status.Cell.Output.Type.Args,
		}
	}

	return result
}

func fromCellDeps(deps []*types.CellDep) []cellDep {
	result := make([]cellDep, len(deps))
	for i := 0; i < len(deps); i++ {
		dep := deps[i]
		result[i] = cellDep{
			OutPoint: outPoint{
				TxHash: dep.OutPoint.TxHash,
				Index:  hexutil.Uint(dep.OutPoint.Index),
			},
			DepType: dep.DepType,
		}
	}
	return result
}

func fromInputs(inputs []*types.CellInput) []cellInput {
	result := make([]cellInput, len(inputs))
	for i := 0; i < len(inputs); i++ {
		input := inputs[i]
		result[i] = cellInput{
			Since: hexutil.Uint64(input.Since),
			PreviousOutput: outPoint{
				TxHash: input.PreviousOutput.TxHash,
				Index:  hexutil.Uint(input.PreviousOutput.Index),
			},
		}
	}
	return result
}

func fromOutputs(outputs []*types.CellOutput) []cellOutput {
	result := make([]cellOutput, len(outputs))
	for i := 0; i < len(outputs); i++ {
		output := outputs[i]
		result[i] = cellOutput{
			Capacity: hexutil.Uint64(output.Capacity),
			Lock: &script{
				CodeHash: output.Lock.CodeHash,
				HashType: output.Lock.HashType,
				Args:     output.Lock.Args,
			},
		}
		if output.Type != nil {
			result[i].Type = &script{
				CodeHash: output.Type.CodeHash,
				HashType: output.Type.HashType,
				Args:     output.Type.Args,
			}
		}
	}
	return result
}

func fromBytesArray(bytes [][]byte) []hexutil.Bytes {
	result := make([]hexutil.Bytes, len(bytes))
	for i, data := range bytes {
		result[i] = data
	}
	return result
}

func fromTransaction(tx *types.Transaction) inTransaction {
	result := inTransaction{
		Version:     hexutil.Uint(tx.Version),
		HeaderDeps:  tx.HeaderDeps,
		CellDeps:    fromCellDeps(tx.CellDeps),
		Inputs:      fromInputs(tx.Inputs),
		Outputs:     fromOutputs(tx.Outputs),
		OutputsData: fromBytesArray(tx.OutputsData),
		Witnesses:   fromBytesArray(tx.Witnesses),
	}
	return result
}

func toNode(result node) *types.Node {
	ret := &types.Node{
		IsOutbound: result.IsOutbound,
		NodeId:     result.NodeId,
		Version:    result.Version,
	}
	ret.Addresses = make([]*types.NodeAddress, len(result.Addresses))
	for i := 0; i < len(result.Addresses); i++ {
		address := result.Addresses[i]
		ret.Addresses[i] = &types.NodeAddress{
			Address: address.Address,
			Score:   uint64(address.Score),
		}
	}
	return ret
}

func toConsensus(c consensus) *types.Consensus {
	result := &types.Consensus{
		Id:                                   c.Id,
		GenesisHash:                          c.GenesisHash,
		DaoTypeHash:                          c.DaoTypeHash,
		Secp256k1Blake160SighashAllTypeHash:  c.Secp256k1Blake160SighashAllTypeHash,
		Secp256k1Blake160MultisigAllTypeHash: c.Secp256k1Blake160MultisigAllTypeHash,
		InitialPrimaryEpochReward:            uint64(c.InitialPrimaryEpochReward),
		SecondaryEpochReward:                 uint64(c.SecondaryEpochReward),
		MaxUnclesNum:                         uint64(c.MaxUnclesNum),
		OrphanRateTarget: types.RationalU256{
			Denom: (*big.Int)(&c.OrphanRateTarget.Denom),
			Numer: (*big.Int)(&c.OrphanRateTarget.Numer),
		},
		EpochDurationTarget: uint64(c.EpochDurationTarget),
		TxProposalWindow: types.ProposalWindow{
			Closest:  uint64(c.TxProposalWindow.Closest),
			Farthest: uint64(c.TxProposalWindow.Farthest),
		},
		ProposerRewardRatio: types.RationalU256{
			Denom: (*big.Int)(&c.ProposerRewardRatio.Denom),
			Numer: (*big.Int)(&c.ProposerRewardRatio.Numer),
		},
		CellbaseMaturity:                  uint64(c.CellbaseMaturity),
		MedianTimeBlockCount:              uint64(c.MedianTimeBlockCount),
		MaxBlockCycles:                    uint64(c.MaxBlockCycles),
		BlockVersion:                      uint(c.BlockVersion),
		TxVersion:                         uint(c.TxVersion),
		TypeIdCodeHash:                    c.TypeIdCodeHash,
		MaxBlockProposalsLimit:            uint64(c.MaxBlockProposalsLimit),
		PrimaryEpochRewardHalvingInterval: uint64(c.PrimaryEpochRewardHalvingInterval),
		PermanentDifficultyInDummy:        c.PermanentDifficultyInDummy,
	}
	return result
}
