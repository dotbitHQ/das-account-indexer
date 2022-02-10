package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/shopspring/decimal"
	"strings"
	"time"
)

func ParserWitnessData(witnessByte []byte) interface{} {
	if len(witnessByte) <= common.WitnessDasTableTypeEndIndex+1 {
		return parserDefaultWitness(witnessByte)
	}
	if string(witnessByte[0:common.WitnessDasCharLen]) != common.WitnessDas {
		return parserDefaultWitness(witnessByte)
	}
	actionDataType := common.Bytes2Hex(witnessByte[common.WitnessDasCharLen:common.WitnessDasTableTypeEndIndex])

	switch actionDataType {
	case common.ActionDataTypeActionData:
		return ParserActionData(witnessByte)
	case common.ActionDataTypeAccountCell:
		return ParserAccountCell(witnessByte)
	case common.ActionDataTypeAccountSaleCell:
		return ParserAccountSaleCell(witnessByte)
	case common.ActionDataTypeAccountAuctionCell:
		return ParserAccountAuctionCell(witnessByte)
	case common.ActionDataTypeProposalCell:
		return ParserProposalCell(witnessByte)
	case common.ActionDataTypePreAccountCell:
		return ParserPreAccountCell(witnessByte)
	case common.ActionDataTypeIncomeCell:
		return ParserIncomeCell(witnessByte)
	case common.ActionDataTypeOfferCell:
		return ParserOfferCell(witnessByte)

	case common.ConfigCellTypeArgsAccount:
		return ParserConfigCellAccount(witnessByte)
	case common.ConfigCellTypeArgsApply:
		return ParserConfigCellApply(witnessByte)
	case common.ConfigCellTypeArgsIncome:
		return ParserConfigCellIncome(witnessByte)
	case common.ConfigCellTypeArgsMain:
		return ParserConfigCellMain(witnessByte)
	case common.ConfigCellTypeArgsPrice:
		return ParserConfigCellPrice(witnessByte)
	case common.ConfigCellTypeArgsProposal:
		return ParserConfigCellProposal(witnessByte)
	case common.ConfigCellTypeArgsProfitRate:
		return ParserConfigCellProfitRate(witnessByte)
	case common.ConfigCellTypeArgsRecordNamespace:
		return ParserConfigCellRecordNamespace(witnessByte)
	case common.ConfigCellTypeArgsRelease:
		return ParserConfigCellRelease(witnessByte)
	case common.ConfigCellTypeArgsUnavailable:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellUnavailable")
	case common.ConfigCellTypeArgsSecondaryMarket:
		return ParserConfigCellSecondaryMarket(witnessByte)
	case common.ConfigCellTypeArgsReverseRecord:
		return ParserConfigCellReverseRecord(witnessByte)

	case common.ConfigCellTypeArgsPreservedAccount00:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount00")
	case common.ConfigCellTypeArgsPreservedAccount01:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount01")
	case common.ConfigCellTypeArgsPreservedAccount02:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount02")
	case common.ConfigCellTypeArgsPreservedAccount03:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount03")
	case common.ConfigCellTypeArgsPreservedAccount04:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount04")
	case common.ConfigCellTypeArgsPreservedAccount05:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount05")
	case common.ConfigCellTypeArgsPreservedAccount06:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount06")
	case common.ConfigCellTypeArgsPreservedAccount07:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount07")
	case common.ConfigCellTypeArgsPreservedAccount08:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount08")
	case common.ConfigCellTypeArgsPreservedAccount09:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount09")
	case common.ConfigCellTypeArgsPreservedAccount10:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount10")
	case common.ConfigCellTypeArgsPreservedAccount11:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount11")
	case common.ConfigCellTypeArgsPreservedAccount12:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount12")
	case common.ConfigCellTypeArgsPreservedAccount13:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount13")
	case common.ConfigCellTypeArgsPreservedAccount14:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount14")
	case common.ConfigCellTypeArgsPreservedAccount15:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount15")
	case common.ConfigCellTypeArgsPreservedAccount16:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount16")
	case common.ConfigCellTypeArgsPreservedAccount17:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount17")
	case common.ConfigCellTypeArgsPreservedAccount18:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount18")
	case common.ConfigCellTypeArgsPreservedAccount19:
		return ParserConfigCellUnavailable(witnessByte, "ConfigCellPreservedAccount19")

	case common.ConfigCellTypeArgsCharSetEmoji:
		return ParserConfigCellTypeArgsCharSetEmoji(witnessByte)
	case common.ConfigCellTypeArgsCharSetDigit:
		return ParserConfigCellTypeArgsCharSetDigit(witnessByte)
	case common.ConfigCellTypeArgsCharSetEn:
		return ParserConfigCellTypeArgsCharSetEn(witnessByte)
	case common.ConfigCellTypeArgsCharSetHanS:
		return ParserConfigCellTypeArgsCharSetHanS(witnessByte)
	case common.ConfigCellTypeArgsCharSetHanT:
		return ParserConfigCellTypeArgsCharSetHanT(witnessByte)

	default:
		return parserDefaultWitness(witnessByte)
	}
}

func parserDefaultWitness(witnessByte []byte) interface{} {
	return map[string]interface{}{
		"name":    "unknown",
		"witness": common.Bytes2Hex(witnessByte),
	}
}

func parserData(data *molecule.Data) (dataEntityOpts []map[string]interface{}) {
	if data.New() != nil && !data.New().IsNone() {
		dataEntityOpts = append(dataEntityOpts, map[string]interface{}{
			"type":   "new",
			"entity": data.New(),
		})
	}
	if data.Old() != nil && !data.Old().IsNone() {
		dataEntityOpts = append(dataEntityOpts, map[string]interface{}{
			"type":   "old",
			"entity": data.Old(),
		})
	}
	if data.Dep() != nil && !data.Dep().IsNone() {
		dataEntityOpts = append(dataEntityOpts, map[string]interface{}{
			"type":   "dep",
			"entity": data.Dep(),
		})
	}

	return dataEntityOpts
}

func parserScript(script *molecule.Script) map[string]interface{} {
	if script == nil {
		return nil
	}

	return map[string]interface{}{
		"code_hash": common.Bytes2Hex(script.CodeHash().RawData()),
		"hash_type": common.Bytes2Hex(script.HashType().AsSlice()),
		"args":      common.Bytes2Hex(script.Args().RawData()),
	}
}

func parserConfig(priceConfig *molecule.PriceConfig) map[string]interface{} {
	if priceConfig == nil {
		return nil
	}

	length, _ := molecule.Bytes2GoU8(priceConfig.Length().RawData())
	newP, _ := molecule.Bytes2GoU64(priceConfig.New().RawData())
	renew, _ := molecule.Bytes2GoU64(priceConfig.Renew().RawData())

	return map[string]interface{}{
		"length": length,
		"new":    ConvertDollar(newP),
		"renew":  ConvertDollar(renew),
	}
}

func ParserActionData(witnessByte []byte) interface{} {
	builder, err := ActionDataBuilderFromWitness(witnessByte)
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(builder.ActionData.AsSlice())),
		"name":         "ActionData",
		"data": map[string]interface{}{
			"action":      builder.Action,
			"action_hash": common.Bytes2Hex(builder.ActionData.Action().RawData()),
			"params":      builder.ParamsStr,
		},
	}
}

func ParserAccountCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	accountCells := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		var accountCell map[string]interface{}
		switch version {
		case common.GoDataEntityVersion1:
			accountCell = parserAccountCellV1(dataEntity)
		case common.GoDataEntityVersion2:
			accountCell = parserAccountCell(dataEntity)
		}
		if accountCell == nil {
			return parserDefaultWitness(witnessByte)
		}
		accountCells[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": accountCell["witness_hash"],
			"entity":       accountCell["entity"],
		}
	}

	return map[string]interface{}{
		"witness": common.Bytes2Hex(witnessByte),
		"name":    "AccountCell",
		"data":    accountCells,
	}
}

func parserAccountCellV1(dataEntity *molecule.DataEntity) map[string]interface{} {
	accountCellV1, _ := molecule.AccountCellDataV1FromSlice(dataEntity.Entity().RawData(), false)
	if accountCellV1 == nil {
		return nil
	}

	registeredAt, _ := molecule.Bytes2GoU64(accountCellV1.RegisteredAt().RawData())
	updatedAt, _ := molecule.Bytes2GoU64(accountCellV1.UpdatedAt().RawData())
	status, _ := molecule.Bytes2GoU64(accountCellV1.Status().RawData())
	var recordsMaps []map[string]interface{}
	for i := uint(0); i < accountCellV1.Records().Len(); i++ {
		record := accountCellV1.Records().Get(i)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		recordsMaps = append(recordsMaps, map[string]interface{}{
			"key":   string(record.RecordKey().RawData()),
			"type":  string(record.RecordType().RawData()),
			"label": string(record.RecordLabel().RawData()),
			"value": string(record.RecordValue().RawData()),
			"ttl":   ConvertMinute(ttl),
		})
	}

	return map[string]interface{}{
		"witness_hash": common.Bytes2Hex(common.Blake2b(accountCellV1.AsSlice())),
		"entity": map[string]interface{}{
			"id":            common.Bytes2Hex(accountCellV1.Id().RawData()),
			"account":       common.AccountCharsToAccount(accountCellV1.Account()),
			"registered_at": ConvertTimestamp(int64(registeredAt)),
			"updated_at":    ConvertTimestamp(int64(updatedAt)),
			"status":        status,
			"records":       recordsMaps,
		},
	}
}

func parserAccountCell(dataEntity *molecule.DataEntity) map[string]interface{} {
	accountCell, _ := molecule.AccountCellDataFromSlice(dataEntity.Entity().RawData(), false)
	if accountCell == nil {
		return nil
	}

	registeredAt, _ := molecule.Bytes2GoU64(accountCell.RegisteredAt().RawData())
	lastTransferAccountAt, _ := molecule.Bytes2GoU64(accountCell.LastTransferAccountAt().RawData())
	lastEditManagerAt, _ := molecule.Bytes2GoU64(accountCell.LastEditManagerAt().RawData())
	lastEditRecordsAt, _ := molecule.Bytes2GoU64(accountCell.LastEditRecordsAt().RawData())
	status, _ := molecule.Bytes2GoU64(accountCell.Status().RawData())
	var recordsMaps []map[string]interface{}
	for i := uint(0); i < accountCell.Records().Len(); i++ {
		record := accountCell.Records().Get(i)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		recordsMaps = append(recordsMaps, map[string]interface{}{
			"key":   string(record.RecordKey().RawData()),
			"type":  string(record.RecordType().RawData()),
			"label": string(record.RecordLabel().RawData()),
			"value": string(record.RecordValue().RawData()),
			"ttl":   ConvertMinute(ttl),
		})
	}

	return map[string]interface{}{
		"witness_hash": common.Bytes2Hex(common.Blake2b(accountCell.AsSlice())),
		"entity": map[string]interface{}{
			"id":                       common.Bytes2Hex(accountCell.Id().RawData()),
			"account":                  common.AccountCharsToAccount(accountCell.Account()),
			"registered_at":            ConvertTimestamp(int64(registeredAt)),
			"last_transfer_account_at": ConvertTimestamp(int64(lastTransferAccountAt)),
			"last_edit_manager_at":     ConvertTimestamp(int64(lastEditManagerAt)),
			"last_edit_records_at":     ConvertTimestamp(int64(lastEditRecordsAt)),
			"status":                   status,
			"records":                  recordsMaps,
		},
	}
}

func ParserAccountSaleCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	accountSaleCells := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		var accountSaleCell map[string]interface{}
		switch version {
		case common.GoDataEntityVersion1:
			accountSaleCell = parserAccountSaleCellV1(dataEntity)
		case common.GoDataEntityVersion2:
			accountSaleCell = parserAccountSaleCell(dataEntity)
		}
		if accountSaleCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		accountSaleCells[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": accountSaleCell["witness_hash"],
			"entity":       accountSaleCell["entity"],
		}
	}

	return map[string]interface{}{
		"witness": common.Bytes2Hex(witnessByte),
		"name":    "AccountSaleCell",
		"data":    accountSaleCells,
	}
}

func parserAccountSaleCellV1(dataEntity *molecule.DataEntity) map[string]interface{} {
	accountSaleCellV1, _ := molecule.AccountSaleCellDataV1FromSlice(dataEntity.Entity().RawData(), false)
	if accountSaleCellV1 == nil {
		return nil
	}
	price, _ := molecule.Bytes2GoU64(accountSaleCellV1.Price().RawData())
	startedAt, _ := molecule.Bytes2GoU64(accountSaleCellV1.StartedAt().RawData())

	return map[string]interface{}{
		"witness_hash": common.Bytes2Hex(common.Blake2b(accountSaleCellV1.AsSlice())),
		"entity": map[string]interface{}{
			"id":          common.Bytes2Hex(accountSaleCellV1.AccountId().RawData()),
			"account":     string(accountSaleCellV1.Account().RawData()),
			"price":       ConvertCapacity(price),
			"description": string(accountSaleCellV1.Description().RawData()),
			"started_at":  ConvertTimestamp(int64(startedAt)),
		},
	}
}

func parserAccountSaleCell(dataEntity *molecule.DataEntity) map[string]interface{} {
	accountSaleCell, _ := molecule.AccountSaleCellDataFromSlice(dataEntity.Entity().RawData(), false)
	if accountSaleCell == nil {
		return nil
	}
	price, _ := molecule.Bytes2GoU64(accountSaleCell.Price().RawData())
	startedAt, _ := molecule.Bytes2GoU64(accountSaleCell.StartedAt().RawData())

	return map[string]interface{}{
		"witness_hash": common.Bytes2Hex(common.Blake2b(accountSaleCell.AsSlice())),
		"entity": map[string]interface{}{
			"id":          common.Bytes2Hex(accountSaleCell.AccountId().RawData()),
			"account":     string(accountSaleCell.Account().RawData()),
			"price":       ConvertCapacity(price),
			"description": string(accountSaleCell.Description().RawData()),
			"started_at":  ConvertTimestamp(int64(startedAt)),
		},
	}
}

func ParserAccountAuctionCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	accountAuctionCells := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		accountAuctionCell, _ := molecule.AccountAuctionCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if accountAuctionCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		openingPrice, _ := molecule.Bytes2GoU64(accountAuctionCell.OpeningPrice().RawData())
		incrementRateEachBid, _ := molecule.Bytes2GoU32(accountAuctionCell.IncrementRateEachBid().RawData())
		startedAt, _ := molecule.Bytes2GoU64(accountAuctionCell.StartedAt().RawData())
		endedAt, _ := molecule.Bytes2GoU64(accountAuctionCell.EndedAt().RawData())
		currentBidPrice, _ := molecule.Bytes2GoU64(accountAuctionCell.CurrentBidPrice().RawData())
		prevBidderProfitRate, _ := molecule.Bytes2GoU32(accountAuctionCell.PrevBidderProfitRate().RawData())

		accountAuctionCells[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": common.Bytes2Hex(common.Blake2b(accountAuctionCell.AsSlice())),
			"entity": map[string]interface{}{
				"id":                      common.Bytes2Hex(accountAuctionCell.AccountId().RawData()),
				"account":                 string(accountAuctionCell.Account().RawData()),
				"description":             string(accountAuctionCell.Description().RawData()),
				"opening_price":           ConvertCapacity(openingPrice),
				"incrementRateEachBid":    ConvertRate(incrementRateEachBid),
				"started_at":              ConvertTimestamp(int64(startedAt)),
				"ended_at":                ConvertTimestamp(int64(endedAt)),
				"current_bidder_lock":     parserScript(accountAuctionCell.CurrentBidderLock()),
				"current_bid_price":       ConvertCapacity(currentBidPrice),
				"prev_bidder_profit_rate": ConvertRate(prevBidderProfitRate),
			},
		}
	}

	return map[string]interface{}{
		"witness": common.Bytes2Hex(witnessByte),
		"name":    "AccountAuctionCell",
		"data":    accountAuctionCells,
	}
}

func ParserProposalCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	proposalCells := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		proposalCell, _ := molecule.ProposalCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if proposalCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		createdAtHeight, _ := molecule.Bytes2GoU64(proposalCell.CreatedAtHeight().RawData())
		var slices []interface{}
		for i := uint(0); i < proposalCell.Slices().Len(); i++ {
			slice := proposalCell.Slices().Get(i)
			var proposalItems []interface{}
			for k := uint(0); k < slice.Len(); k++ {
				proposalItem := slice.Get(k)
				itemType, _ := molecule.Bytes2GoU8(proposalItem.ItemType().RawData())
				proposalItems = append(proposalItems, map[string]interface{}{
					"id":        common.Bytes2Hex(proposalItem.AccountId().RawData()),
					"item_type": itemType,
					"next":      common.Bytes2Hex(proposalItem.Next().RawData()),
				})
			}
			slices = append(slices, proposalItems)
		}

		proposalCells[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": common.Bytes2Hex(common.Blake2b(proposalCell.AsSlice())),
			"entity": map[string]interface{}{
				"proposal_lock":     parserScript(proposalCell.ProposerLock()),
				"created_at_height": createdAtHeight,
				"slices":            slices,
			},
		}
	}

	return map[string]interface{}{
		"witness": common.Bytes2Hex(witnessByte),
		"name":    "ProposalCell",
		"data":    proposalCells,
	}
}

func ParserPreAccountCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	preAccountCells := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		preAccountCell, _ := molecule.PreAccountCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if preAccountCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		inviterLock, _ := preAccountCell.InviterLock().IntoScript()
		channelLock, _ := preAccountCell.ChannelLock().IntoScript()
		quote, _ := molecule.Bytes2GoU64(preAccountCell.Quote().RawData())
		invitedDiscount, _ := molecule.Bytes2GoU32(preAccountCell.InvitedDiscount().RawData())
		createdAt, _ := molecule.Bytes2GoU64(preAccountCell.CreatedAt().RawData())

		preAccountCells[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": common.Bytes2Hex(common.Blake2b(preAccountCell.AsSlice())),
			"entity": map[string]interface{}{
				"account":          common.AccountCharsToAccount(preAccountCell.Account()),
				"owner_lock_args":  common.Bytes2Hex(preAccountCell.OwnerLockArgs().RawData()),
				"inviter_id":       common.Bytes2Hex(preAccountCell.InviterId().RawData()),
				"refund_lock":      parserScript(preAccountCell.RefundLock()),
				"inviter_lock":     parserScript(inviterLock),
				"channel_lock":     parserScript(channelLock),
				"price":            parserConfig(preAccountCell.Price()),
				"quote":            quote,
				"invited_discount": ConvertRate(invitedDiscount),
				"created_at":       ConvertTimestamp(int64(createdAt)),
			},
		}
	}

	return map[string]interface{}{
		"witness": common.Bytes2Hex(witnessByte),
		"name":    "PreAccountCell",
		"data":    preAccountCells,
	}
}

func ParserIncomeCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	incomeCells := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		incomeCell, _ := molecule.IncomeCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if incomeCell == nil {
			return parserDefaultWitness(witnessByte)
		}

		var recordsMaps []map[string]interface{}
		for i := uint(0); i < incomeCell.Records().Len(); i++ {
			record := incomeCell.Records().Get(i)
			capacity, _ := molecule.Bytes2GoU64(record.Capacity().RawData())
			recordsMaps = append(recordsMaps, map[string]interface{}{
				"belong_to": map[string]interface{}{
					"code_hash": common.Bytes2Hex(record.BelongTo().CodeHash().RawData()),
					"hash_type": common.Bytes2Hex(record.BelongTo().HashType().AsSlice()),
					"args":      common.Bytes2Hex(record.BelongTo().Args().RawData()),
				},
				"capacity": ConvertCapacity(capacity),
			})
		}

		incomeCells[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": common.Bytes2Hex(common.Blake2b(incomeCell.AsSlice())),
			"entity": map[string]interface{}{
				"creator": map[string]interface{}{
					"code_hash": common.Bytes2Hex(incomeCell.Creator().CodeHash().RawData()),
					"hash_type": common.Bytes2Hex(incomeCell.Creator().HashType().AsSlice()),
				},
				"records": recordsMaps,
			},
		}
	}

	return map[string]interface{}{
		"witness": common.Bytes2Hex(witnessByte),
		"name":    "IncomeCell",
		"data":    incomeCells,
	}
}

func ParserOfferCell(witnessByte []byte) interface{} {
	data, _ := molecule.DataFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if data == nil {
		return parserDefaultWitness(witnessByte)
	}

	offerCells := map[string]interface{}{}
	for _, v := range parserData(data) {
		dataEntity, _ := molecule.DataEntityFromSlice(v["entity"].(*molecule.DataEntityOpt).AsSlice(), false)
		if dataEntity == nil {
			return parserDefaultWitness(witnessByte)
		}

		version, _ := molecule.Bytes2GoU32(dataEntity.Version().RawData())
		index, _ := molecule.Bytes2GoU32(dataEntity.Index().RawData())
		offerCell, _ := molecule.OfferCellDataFromSlice(dataEntity.Entity().RawData(), false)
		if offerCell == nil {
			return parserDefaultWitness(witnessByte)
		}
		price, _ := molecule.Bytes2GoU64(offerCell.Price().RawData())

		offerCells[v["type"].(string)] = map[string]interface{}{
			"version":      version,
			"index":        index,
			"witness_hash": common.Bytes2Hex(common.Blake2b(offerCell.AsSlice())),
			"entity": map[string]interface{}{
				"account":      string(offerCell.Account().RawData()),
				"price":        price,
				"message":      string(offerCell.Message().RawData()),
				"inviter_lock": parserScript(offerCell.InviterLock()),
				"channel_lock": parserScript(offerCell.ChannelLock()),
			},
		}
	}

	return map[string]interface{}{
		"witness": common.Bytes2Hex(witnessByte),
		"name":    "OfferCell",
		"data":    offerCells,
	}
}

func ParserConfigCellAccount(witnessByte []byte) interface{} {
	configCellAccount, _ := molecule.ConfigCellAccountFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellAccount == nil {
		return parserDefaultWitness(witnessByte)
	}

	maxLength, _ := molecule.Bytes2GoU32(configCellAccount.MaxLength().RawData())
	basicCapacity, _ := molecule.Bytes2GoU64(configCellAccount.BasicCapacity().RawData())
	preparedFeeCapacity, _ := molecule.Bytes2GoU64(configCellAccount.PreparedFeeCapacity().RawData())
	expirationGracePeriod, _ := molecule.Bytes2GoU32(configCellAccount.ExpirationGracePeriod().RawData())
	recordMinTtl, _ := molecule.Bytes2GoU32(configCellAccount.RecordMinTtl().RawData())
	recordSizeLimit, _ := molecule.Bytes2GoU32(configCellAccount.RecordSizeLimit().RawData())
	transferAccountFee, _ := molecule.Bytes2GoU64(configCellAccount.TransferAccountFee().RawData())
	editManagerFee, _ := molecule.Bytes2GoU64(configCellAccount.EditManagerFee().RawData())
	editRecordsFee, _ := molecule.Bytes2GoU64(configCellAccount.EditRecordsFee().RawData())
	commonFee, _ := molecule.Bytes2GoU64(configCellAccount.CommonFee().RawData())
	transferAccountThrottle, _ := molecule.Bytes2GoU32(configCellAccount.TransferAccountThrottle().RawData())
	editManagerThrottle, _ := molecule.Bytes2GoU32(configCellAccount.EditManagerThrottle().RawData())
	editRecordsThrottle, _ := molecule.Bytes2GoU32(configCellAccount.EditRecordsThrottle().RawData())
	commonThrottle, _ := molecule.Bytes2GoU32(configCellAccount.CommonThrottle().RawData())

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellAccount.AsSlice())),
		"name":         "ConfigCellAccount",
		"data": map[string]interface{}{
			"max_length":                maxLength,
			"basic_capacity":            ConvertCapacity(basicCapacity),
			"prepared_fee_capacity":     ConvertCapacity(preparedFeeCapacity),
			"expiration_grace_period":   ConvertDay(expirationGracePeriod),
			"record_min_ttl":            ConvertMinute(recordMinTtl),
			"record_size_limit":         recordSizeLimit,
			"transfer_account_fee":      ConvertCapacity(transferAccountFee),
			"edit_manager_fee":          ConvertCapacity(editManagerFee),
			"edit_records_fee":          ConvertCapacity(editRecordsFee),
			"common_fee":                ConvertCapacity(commonFee),
			"transfer_account_throttle": ConvertMinute(transferAccountThrottle),
			"edit_manager_throttle":     ConvertMinute(editManagerThrottle),
			"edit_records_throttle":     ConvertMinute(editRecordsThrottle),
			"common_throttle":           ConvertMinute(commonThrottle),
		},
	}
}

func ParserConfigCellApply(witnessByte []byte) interface{} {
	configCellApply, _ := molecule.ConfigCellApplyFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellApply == nil {
		return parserDefaultWitness(witnessByte)
	}

	applyMinWaitingBlockNumber, _ := molecule.Bytes2GoU32(configCellApply.ApplyMinWaitingBlockNumber().RawData())
	applyMaxWaitingBlockNumber, _ := molecule.Bytes2GoU32(configCellApply.ApplyMaxWaitingBlockNumber().RawData())
	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellApply.AsSlice())),
		"name":         "ConfigCellApply",
		"data": map[string]interface{}{
			"apply_min_waiting_block_number": applyMinWaitingBlockNumber,
			"apply_max_waiting_block_number": applyMaxWaitingBlockNumber,
		},
	}
}

func ParserConfigCellIncome(witnessByte []byte) interface{} {
	configCellIncome, _ := molecule.ConfigCellIncomeFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellIncome == nil {
		return parserDefaultWitness(witnessByte)
	}

	basicCapacity, _ := molecule.Bytes2GoU64(configCellIncome.BasicCapacity().RawData())
	maxRecords, _ := molecule.Bytes2GoU32(configCellIncome.MaxRecords().RawData())
	minTransferCapacity, _ := molecule.Bytes2GoU64(configCellIncome.MinTransferCapacity().RawData())
	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellIncome.AsSlice())),
		"name":         "ConfigCellIncome",
		"data": map[string]interface{}{
			"basic_capacity":        ConvertCapacity(basicCapacity),
			"max_records":           maxRecords,
			"min_transfer_capacity": ConvertCapacity(minTransferCapacity),
		},
	}
}

func ParserConfigCellMain(witnessByte []byte) interface{} {
	configCellMain, _ := molecule.ConfigCellMainFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellMain == nil {
		return parserDefaultWitness(witnessByte)
	}

	status, _ := molecule.Bytes2GoU8(configCellMain.Status().RawData())
	ckbSignAllIndex, _ := molecule.Bytes2GoU32(configCellMain.DasLockOutPointTable().CkbSignall().Index().RawData())
	ckbMultiSignIndex, _ := molecule.Bytes2GoU32(configCellMain.DasLockOutPointTable().CkbMultisign().Index().RawData())
	ckbAnyoneCanPayIndex, _ := molecule.Bytes2GoU32(configCellMain.DasLockOutPointTable().CkbAnyoneCanPay().Index().RawData())
	ethIndex, _ := molecule.Bytes2GoU32(configCellMain.DasLockOutPointTable().Eth().Index().RawData())
	tronIndex, _ := molecule.Bytes2GoU32(configCellMain.DasLockOutPointTable().Tron().Index().RawData())
	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellMain.AsSlice())),
		"name":         "ConfigCellMain",
		"data": map[string]interface{}{
			"status": status,
			"type_id_table": map[string]interface{}{
				"account_cell":         common.Bytes2Hex(configCellMain.TypeIdTable().AccountCell().RawData()),
				"apply_register_cell":  common.Bytes2Hex(configCellMain.TypeIdTable().ApplyRegisterCell().RawData()),
				"balance_cell":         common.Bytes2Hex(configCellMain.TypeIdTable().BalanceCell().RawData()),
				"income_cell":          common.Bytes2Hex(configCellMain.TypeIdTable().IncomeCell().RawData()),
				"pre_account_cell":     common.Bytes2Hex(configCellMain.TypeIdTable().PreAccountCell().RawData()),
				"proposal_cell":        common.Bytes2Hex(configCellMain.TypeIdTable().ProposalCell().RawData()),
				"account_sale_cell":    common.Bytes2Hex(configCellMain.TypeIdTable().AccountSaleCell().RawData()),
				"account_auction_cell": common.Bytes2Hex(configCellMain.TypeIdTable().AccountAuctionCell().RawData()),
				"offer_cell":           common.Bytes2Hex(configCellMain.TypeIdTable().OfferCell().RawData()),
				"reverse_record_cell":  common.Bytes2Hex(configCellMain.TypeIdTable().ReverseRecordCell().RawData()),
			},
			"das_lock_out_point_table": map[string]interface{}{
				"ckb_sign_all": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbSignall().TxHash().RawData()),
					"index":   ckbSignAllIndex,
				},
				"ckb_multi_sign": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbMultisign().TxHash().RawData()),
					"index":   ckbMultiSignIndex,
				},
				"ckb_anyone_can_pay": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().CkbAnyoneCanPay().TxHash().RawData()),
					"index":   ckbAnyoneCanPayIndex,
				},
				"eth": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().Eth().TxHash().RawData()),
					"index":   ethIndex,
				},
				"tron": map[string]interface{}{
					"tx_hash": common.Bytes2Hex(configCellMain.DasLockOutPointTable().Tron().TxHash().RawData()),
					"index":   tronIndex,
				},
			},
		},
	}
}

func ParserConfigCellPrice(witnessByte []byte) interface{} {
	configCellPrice, _ := molecule.ConfigCellPriceFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellPrice == nil {
		return parserDefaultWitness(witnessByte)
	}

	var prices []interface{}
	for i := uint(0); i < configCellPrice.Prices().Len(); i++ {
		prices = append(prices, parserConfig(configCellPrice.Prices().Get(i)))
	}

	invitedDiscount, _ := molecule.Bytes2GoU32(configCellPrice.Discount().InvitedDiscount().RawData())
	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellPrice.AsSlice())),
		"name":         "ConfigCellPrice",
		"data": map[string]interface{}{
			"discount": map[string]interface{}{
				"invited_discount": ConvertRate(invitedDiscount),
			},
			"prices": prices,
		},
	}
}

func ParserConfigCellProposal(witnessByte []byte) interface{} {
	configCellProposal, _ := molecule.ConfigCellProposalFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellProposal == nil {
		return parserDefaultWitness(witnessByte)
	}

	proposalMinConfirmInterval, _ := molecule.Bytes2GoU8(configCellProposal.ProposalMinConfirmInterval().RawData())
	proposalMinRecycleInterval, _ := molecule.Bytes2GoU8(configCellProposal.ProposalMinRecycleInterval().RawData())
	proposalMinExtendInterval, _ := molecule.Bytes2GoU8(configCellProposal.ProposalMinExtendInterval().RawData())
	proposalMaxAccountAffect, _ := molecule.Bytes2GoU32(configCellProposal.ProposalMaxAccountAffect().RawData())
	proposalMaxPreAccountContain, _ := molecule.Bytes2GoU32(configCellProposal.ProposalMaxPreAccountContain().RawData())
	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellProposal.AsSlice())),
		"name":         "ConfigCellProposal",
		"data": map[string]interface{}{
			"proposal_min_confirm_interval":    proposalMinConfirmInterval,
			"proposal_min_recycle_interval":    proposalMinRecycleInterval,
			"proposal_min_extend_interval":     proposalMinExtendInterval,
			"proposal_max_account_affect":      proposalMaxAccountAffect,
			"proposal_max_pre_account_contain": proposalMaxPreAccountContain,
		},
	}
}

func ParserConfigCellProfitRate(witnessByte []byte) interface{} {
	configCellProfitRate, _ := molecule.ConfigCellProfitRateFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellProfitRate == nil {
		return parserDefaultWitness(witnessByte)
	}

	inviter, _ := molecule.Bytes2GoU32(configCellProfitRate.Inviter().RawData())
	channel, _ := molecule.Bytes2GoU32(configCellProfitRate.Channel().RawData())
	proposalCreate, _ := molecule.Bytes2GoU32(configCellProfitRate.ProposalCreate().RawData())
	proposalConfirm, _ := molecule.Bytes2GoU32(configCellProfitRate.ProposalConfirm().RawData())
	incomeConsolidate, _ := molecule.Bytes2GoU32(configCellProfitRate.IncomeConsolidate().RawData())
	saleBuyerInviter, _ := molecule.Bytes2GoU32(configCellProfitRate.SaleBuyerInviter().RawData())
	saleBuyerChannel, _ := molecule.Bytes2GoU32(configCellProfitRate.SaleBuyerChannel().RawData())
	saleDas, _ := molecule.Bytes2GoU32(configCellProfitRate.SaleDas().RawData())
	auctionBidderInviter, _ := molecule.Bytes2GoU32(configCellProfitRate.AuctionBidderInviter().RawData())
	auctionBidderChannel, _ := molecule.Bytes2GoU32(configCellProfitRate.AuctionBidderChannel().RawData())
	auctionDas, _ := molecule.Bytes2GoU32(configCellProfitRate.AuctionDas().RawData())
	auctionPrevBidder, _ := molecule.Bytes2GoU32(configCellProfitRate.AuctionPrevBidder().RawData())

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellProfitRate.AsSlice())),
		"name":         "ConfigCellProfitRate",
		"data": map[string]interface{}{
			"inviter":                ConvertRate(inviter),
			"channel":                ConvertRate(channel),
			"proposal_create":        ConvertRate(proposalCreate),
			"proposal_confirm":       ConvertRate(proposalConfirm),
			"income_consolidate":     ConvertRate(incomeConsolidate),
			"sale_buyer_inviter":     ConvertRate(saleBuyerInviter),
			"sale_buyer_channel":     ConvertRate(saleBuyerChannel),
			"sale_das":               ConvertRate(saleDas),
			"auction_bidder_inviter": ConvertRate(auctionBidderInviter),
			"auction_bidder_channel": ConvertRate(auctionBidderChannel),
			"auction_das":            ConvertRate(auctionDas),
			"auction_prev_bidder":    ConvertRate(auctionPrevBidder),
		},
	}
}

func ParserConfigCellRecordNamespace(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"name":         "ConfigCellRecordNamespace",
		"data": map[string]interface{}{
			"length":                       dataLength,
			"config_cell_record_namespace": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellRelease(witnessByte []byte) interface{} {
	configCellRelease, _ := molecule.ConfigCellReleaseFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellRelease == nil {
		return parserDefaultWitness(witnessByte)
	}

	var releaseRules []interface{}
	for i := uint(0); i < configCellRelease.ReleaseRules().Len(); i++ {
		releaseRule := configCellRelease.ReleaseRules().Get(i)
		length, _ := molecule.Bytes2GoU32(releaseRule.Length().RawData())
		releaseStart, _ := molecule.Bytes2GoU64(releaseRule.ReleaseStart().RawData())
		releaseEnd, _ := molecule.Bytes2GoU64(releaseRule.ReleaseEnd().RawData())
		releaseRules = append(releaseRules, map[string]interface{}{
			"length":        length,
			"release_start": ConvertTimestamp(int64(releaseStart)),
			"release_end":   ConvertTimestamp(int64(releaseEnd)),
		})
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellRelease.AsSlice())),
		"name":         "ConfigCellRelease",
		"data": map[string]interface{}{
			"release_rules": releaseRules,
		},
	}
}

func ParserConfigCellUnavailable(witnessByte []byte, action string) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"name":         action,
		"data": map[string]interface{}{
			"length": dataLength,
		},
	}
}

func ParserConfigCellSecondaryMarket(witnessByte []byte) interface{} {
	configCellSecondaryMarket, _ := molecule.ConfigCellSecondaryMarketFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellSecondaryMarket == nil {
		return parserDefaultWitness(witnessByte)
	}

	commonFee, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.CommonFee().RawData())
	saleMinPrice, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.SaleMinPrice().RawData())
	saleExpirationLimit, _ := molecule.Bytes2GoU32(configCellSecondaryMarket.SaleExpirationLimit().RawData())
	saleDescriptionBytesLimit, _ := molecule.Bytes2GoU32(configCellSecondaryMarket.SaleDescriptionBytesLimit().RawData())
	saleCellBasicCapacity, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.SaleCellBasicCapacity().RawData())
	saleCellPreparedFeeCapacity, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.SaleCellPreparedFeeCapacity().RawData())
	auctionMaxExtendableDuration, _ := molecule.Bytes2GoU32(configCellSecondaryMarket.AuctionMaxExtendableDuration().RawData())
	auctionDurationIncrementEachBid, _ := molecule.Bytes2GoU32(configCellSecondaryMarket.AuctionDurationIncrementEachBid().RawData())
	auctionMinOpeningPrice, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.AuctionMinOpeningPrice().RawData())
	auctionMinIncrementRateEachBid, _ := molecule.Bytes2GoU32(configCellSecondaryMarket.AuctionMinIncrementRateEachBid().RawData())
	auctionDescriptionBytesLimit, _ := molecule.Bytes2GoU32(configCellSecondaryMarket.AuctionDescriptionBytesLimit().RawData())
	auctionCellBasicCapacity, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.AuctionCellBasicCapacity().RawData())
	auctionCellPreparedFeeCapacity, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.AuctionCellPreparedFeeCapacity().RawData())
	offerMinPrice, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.OfferMinPrice().RawData())
	offerCellBasicCapacity, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.OfferCellBasicCapacity().RawData())
	offerCellPreparedFeeCapacity, _ := molecule.Bytes2GoU64(configCellSecondaryMarket.OfferCellPreparedFeeCapacity().RawData())
	offerMessageBytesLimit, _ := molecule.Bytes2GoU32(configCellSecondaryMarket.OfferMessageBytesLimit().RawData())

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellSecondaryMarket.AsSlice())),
		"name":         "ConfigCellSecondaryMarket",
		"data": map[string]interface{}{
			"common_fee":                          ConvertCapacity(commonFee),
			"sale_min_price":                      ConvertCapacity(saleMinPrice),
			"sale_expiration_limit":               saleExpirationLimit,
			"sale_description_bytes_limit":        saleDescriptionBytesLimit,
			"sale_cell_basic_capacity":            ConvertCapacity(saleCellBasicCapacity),
			"sale_cell_prepared_fee_capacity":     ConvertCapacity(saleCellPreparedFeeCapacity),
			"auction_max_extendable_duration":     auctionMaxExtendableDuration,
			"auction_duration_increment_each_bid": auctionDurationIncrementEachBid,
			"auction_min_opening_price":           ConvertCapacity(auctionMinOpeningPrice),
			"auction_min_increment_rate_each_bid": auctionMinIncrementRateEachBid,
			"auction_description_bytes_limit":     auctionDescriptionBytesLimit,
			"auction_cell_basic_capacity":         ConvertCapacity(auctionCellBasicCapacity),
			"auction_cell_prepared_fee_capacity":  ConvertCapacity(auctionCellPreparedFeeCapacity),
			"offer_min_price":                     ConvertCapacity(offerMinPrice),
			"offer_cell_basic_capacity":           ConvertCapacity(offerCellBasicCapacity),
			"offer_cell_prepared_fee_capacity":    ConvertCapacity(offerCellPreparedFeeCapacity),
			"offer_message_bytes_limit":           offerMessageBytesLimit,
		},
	}
}

func ParserConfigCellReverseRecord(witnessByte []byte) interface{} {
	configCellReverseRecord, _ := molecule.ConfigCellReverseResolutionFromSlice(witnessByte[common.WitnessDasTableTypeEndIndex:], false)
	if configCellReverseRecord == nil {
		return parserDefaultWitness(witnessByte)
	}

	commonFee, _ := molecule.Bytes2GoU64(configCellReverseRecord.CommonFee().RawData())
	recordPreparedFeeCapacity, _ := molecule.Bytes2GoU64(configCellReverseRecord.RecordPreparedFeeCapacity().RawData())
	recordBasicCapacity, _ := molecule.Bytes2GoU64(configCellReverseRecord.RecordBasicCapacity().RawData())
	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(configCellReverseRecord.AsSlice())),
		"name":         "ConfigCellReverseRecord",
		"data": map[string]interface{}{
			"common_fee":                   ConvertCapacity(commonFee),
			"record_prepared_fee_capacity": ConvertCapacity(recordPreparedFeeCapacity),
			"record_basic_capacity":        ConvertCapacity(recordBasicCapacity),
		},
	}
}

func ParserConfigCellTypeArgsCharSetEmoji(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"name":         "ConfigCellTypeArgsCharSetEmoji",
		"data": map[string]interface{}{
			"length":            dataLength,
			"config_cell_emoji": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetDigit(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"name":         "ConfigCellTypeArgsCharSetDigit",
		"data": map[string]interface{}{
			"length":            dataLength,
			"config_cell_digit": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetEn(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"name":         "ConfigCellTypeArgsCharSetEn",
		"data": map[string]interface{}{
			"length":         dataLength,
			"config_cell_en": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetHanS(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"name":         "ConfigCellTypeArgsCharSetHanS",
		"data": map[string]interface{}{
			"length":            dataLength,
			"config_cell_han_s": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ParserConfigCellTypeArgsCharSetHanT(witnessByte []byte) interface{} {
	slice := witnessByte[common.WitnessDasTableTypeEndIndex:]
	dataLength, err := molecule.Bytes2GoU32(slice[:4])
	if err != nil {
		return parserDefaultWitness(witnessByte)
	}

	return map[string]interface{}{
		"witness":      common.Bytes2Hex(witnessByte),
		"witness_hash": common.Bytes2Hex(common.Blake2b(slice)),
		"name":         "ConfigCellTypeArgsCharSetHanT",
		"data": map[string]interface{}{
			"length":            dataLength,
			"config_cell_han_t": strings.Split(string(slice[4:dataLength]), string([]byte{0x00})),
		},
	}
}

func ConvertMinute(minute uint32) string {
	return fmt.Sprintf("%d (%d minutes)", minute, minute/60)
}

func ConvertDay(day uint32) string {
	return fmt.Sprintf("%d (%d days)", day, day/60/60/24)
}

func ConvertTimestamp(timestamp int64) string {
	return fmt.Sprintf("%d (%s)", timestamp, time.Unix(timestamp, 0).Format("2006-01-02 15:04:05"))
}

func ConvertDollar(dollar uint64) string {
	capacityDec, _ := decimal.NewFromString(fmt.Sprintf("%d", dollar))
	return fmt.Sprintf("%d ($%s)", dollar, capacityDec.DivRound(decimal.NewFromInt(1000000), 6))
}

func ConvertCapacity(capacity uint64) string {
	capacityDec, _ := decimal.NewFromString(fmt.Sprintf("%d", capacity))
	return fmt.Sprintf("%d (%s CKB)", capacity, capacityDec.DivRound(decimal.NewFromInt(100000000), 8))
}

func ConvertRate(rate uint32) string {
	return fmt.Sprintf("%d (%d%%)", rate, rate/100)
}
