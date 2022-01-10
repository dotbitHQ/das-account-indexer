package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

type ConfigCellDataBuilder struct {
	ConfigCellAccount           *molecule.ConfigCellAccount
	ConfigCellPrice             *molecule.ConfigCellPrice
	PriceConfigMap              map[uint8]*molecule.PriceConfig
	ConfigCellSecondaryMarket   *molecule.ConfigCellSecondaryMarket
	ConfigCellIncome            *molecule.ConfigCellIncome
	ConfigCellProfitRate        *molecule.ConfigCellProfitRate
	ConfigCellMain              *molecule.ConfigCellMain
	ConfigCellReverseResolution *molecule.ConfigCellReverseResolution
	ConfigCellProposal          *molecule.ConfigCellProposal
	ConfigCellApply             *molecule.ConfigCellApply
	ConfigCellRelease           *molecule.ConfigCellRelease
	ConfigCellRecordKeys        []string
	ConfigCellEmojis            []string
}

func ConfigCellDataBuilderRefByTypeArgs(builder *ConfigCellDataBuilder, tx *types.Transaction, configCellTypeArgs common.ConfigCellTypeArgs) error {
	var configCellDataBys []byte
	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		if actionDataType == configCellTypeArgs {
			configCellDataBys = dataBys
			return false, nil
		}
		return true, nil
	})
	if err != nil {
		return fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}

	switch configCellTypeArgs {
	case common.ConfigCellTypeArgsAccount:
		ConfigCellAccount, err := molecule.ConfigCellAccountFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellAccountFromSlice err: %s", err.Error())
		}
		builder.ConfigCellAccount = ConfigCellAccount
	case common.ConfigCellTypeArgsPrice:
		ConfigCellPrice, err := molecule.ConfigCellPriceFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellPriceFromSlice err: %s", err.Error())
		}
		builder.ConfigCellPrice = ConfigCellPrice
		builder.PriceConfigMap = make(map[uint8]*molecule.PriceConfig)
		prices := builder.ConfigCellPrice.Prices()
		for i, count := uint(0), prices.Len(); i < count; i++ {
			price, err := molecule.PriceConfigFromSlice(prices.Get(i).AsSlice(), false)
			if err != nil {
				return fmt.Errorf("PriceConfigFromSlice err: %s", err.Error())
			}
			length, err := molecule.Bytes2GoU8(price.Length().RawData())
			if err != nil {
				return fmt.Errorf("price.Length() err: %s", err.Error())
			}
			builder.PriceConfigMap[length] = price
		}
	case common.ConfigCellTypeArgsApply:
		ConfigCellApply, err := molecule.ConfigCellApplyFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
		}
		builder.ConfigCellApply = ConfigCellApply
	case common.ConfigCellTypeArgsRelease:
		ConfigCellRelease, err := molecule.ConfigCellReleaseFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
		}
		builder.ConfigCellRelease = ConfigCellRelease
	case common.ConfigCellTypeArgsSecondaryMarket:
		ConfigCellSecondaryMarket, err := molecule.ConfigCellSecondaryMarketFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellSecondaryMarketFromSlice err: %s", err.Error())
		}
		builder.ConfigCellSecondaryMarket = ConfigCellSecondaryMarket
	case common.ConfigCellTypeArgsIncome:
		ConfigCellIncome, err := molecule.ConfigCellIncomeFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellIncomeFromSlice err: %s", err.Error())
		}
		builder.ConfigCellIncome = ConfigCellIncome
	case common.ConfigCellTypeArgsProfitRate:
		ConfigCellProfitRate, err := molecule.ConfigCellProfitRateFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellProfitRateFromSlice err: %s", err.Error())
		}
		builder.ConfigCellProfitRate = ConfigCellProfitRate
	case common.ConfigCellTypeArgsMain:
		ConfigCellMain, err := molecule.ConfigCellMainFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellMainFromSlice err: %s", err.Error())
		}
		builder.ConfigCellMain = ConfigCellMain
	case common.ConfigCellTypeArgsReverseRecord:
		ConfigCellReverseResolution, err := molecule.ConfigCellReverseResolutionFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellReverseResolutionFromSlice err: %s", err.Error())
		}
		builder.ConfigCellReverseResolution = ConfigCellReverseResolution
	case common.ConfigCellTypeArgsProposal:
		ConfigCellProposal, err := molecule.ConfigCellProposalFromSlice(configCellDataBys, false)
		if err != nil {
			return fmt.Errorf("ConfigCellProposalFromSlice err: %s", err.Error())
		}
		builder.ConfigCellProposal = ConfigCellProposal
	case common.ConfigCellTypeArgsRecordNamespace:
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("key name space len err: %s", err.Error())
		}
		builder.ConfigCellRecordKeys = strings.Split(string(configCellDataBys[4:dataLength]), string([]byte{0x00}))
	case common.ConfigCellTypeArgsCharSetEmoji:
		dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
		if err != nil {
			return fmt.Errorf("key name space len err: %s", err.Error())
		}
		builder.ConfigCellEmojis = strings.Split(string(configCellDataBys[4:dataLength]), string([]byte{0x00}))
	case common.ConfigCellTypeArgsUnavailable:
		fmt.Println(string(configCellDataBys))
	//dataLength, err := molecule.Bytes2GoU32(configCellDataBys[:4])
	//if err != nil {
	//	return fmt.Errorf("key name space len err: %s", err.Error())
	//}
	//fmt.Println(string(configCellDataBys[4:dataLength]))
	case common.ConfigCellTypeArgsPreservedAccount00:
		fmt.Println(string(configCellDataBys))
	}
	return nil
}

func ConfigCellDataBuilderByTypeArgs(tx *types.Transaction, configCellTypeArgs common.ConfigCellTypeArgs) (*ConfigCellDataBuilder, error) {
	var resp ConfigCellDataBuilder

	err := ConfigCellDataBuilderRefByTypeArgs(&resp, tx, configCellTypeArgs)
	if err != nil {
		return nil, fmt.Errorf("ConfigCellDataBuilderRefByTypeArgs err: %s", err.Error())
	}

	return &resp, nil
}

func (c *ConfigCellDataBuilder) PriceInvitedDiscount() (uint32, error) {
	if c.ConfigCellPrice != nil {
		return molecule.Bytes2GoU32(c.ConfigCellPrice.Discount().InvitedDiscount().RawData())
	}
	return 0, fmt.Errorf("ConfigCellPrice is nil")
}

func (c *ConfigCellDataBuilder) RecordBasicCapacity() (uint64, error) {
	if c.ConfigCellReverseResolution != nil {
		return molecule.Bytes2GoU64(c.ConfigCellReverseResolution.RecordBasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellReverseResolution is nil")
}

func (c *ConfigCellDataBuilder) RecordPreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellReverseResolution != nil {
		return molecule.Bytes2GoU64(c.ConfigCellReverseResolution.RecordPreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellReverseResolution is nil")
}

func (c *ConfigCellDataBuilder) RecordCommonFee() (uint64, error) {
	if c.ConfigCellReverseResolution != nil {
		return molecule.Bytes2GoU64(c.ConfigCellReverseResolution.CommonFee().RawData())
	}
	return 0, fmt.Errorf("ConfigCellReverseResolution is nil")
}

func (c *ConfigCellDataBuilder) BasicCapacity() (uint64, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellAccount.BasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) TransferAccountThrottle() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.TransferAccountThrottle().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) EditRecordsThrottle() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.EditRecordsThrottle().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) RecordMinTtl() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.RecordMinTtl().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) ExpirationGracePeriod() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.ExpirationGracePeriod().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) EditManagerThrottle() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.EditManagerThrottle().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) MaxLength() (uint32, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU32(c.ConfigCellAccount.MaxLength().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) AccountCommonFee() (uint64, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellAccount.CommonFee().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) PreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellAccount != nil {
		return molecule.Bytes2GoU64(c.ConfigCellAccount.PreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellAccount is nil")
}

func (c *ConfigCellDataBuilder) AccountPrice(length uint8) (uint64, uint64, error) {
	if length > 5 {
		length = 5
	}
	if c.PriceConfigMap != nil {
		price, ok := c.PriceConfigMap[length]
		if ok {
			newPrice, err := molecule.Bytes2GoU64(price.New().RawData())
			if err != nil {
				return 0, 0, fmt.Errorf("price.New() err: %s", err.Error())
			}
			renewPrice, err := molecule.Bytes2GoU64(price.Renew().RawData())
			if err != nil {
				return 0, 0, fmt.Errorf("price.Renew() err: %s", err.Error())
			}
			return newPrice, renewPrice, nil
		}
	}
	return 0, 0, fmt.Errorf("not exist price of length[%d]", length)
}

func (c *ConfigCellDataBuilder) PriceConfig(length uint8) *molecule.PriceConfig {
	if length > 5 {
		length = 5
	}
	if c.PriceConfigMap != nil {
		if price, ok := c.PriceConfigMap[length]; ok {
			return price
		}
	}
	return nil
}

func (c *ConfigCellDataBuilder) SaleCellBasicCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.SaleCellBasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) SaleMinPrice() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.SaleMinPrice().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) SaleCellPreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.SaleCellPreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) CommonFee() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.CommonFee().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferCellBasicCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.OfferCellBasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferMinPrice() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.OfferMinPrice().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferCellPreparedFeeCapacity() (uint64, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU64(c.ConfigCellSecondaryMarket.OfferCellPreparedFeeCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) OfferMessageBytesLimit() (uint32, error) {
	if c.ConfigCellSecondaryMarket != nil {
		return molecule.Bytes2GoU32(c.ConfigCellSecondaryMarket.OfferMessageBytesLimit().RawData())
	}
	return 0, fmt.Errorf("ConfigCellSecondaryMarket is nil")
}

func (c *ConfigCellDataBuilder) IncomeBasicCapacity() (uint64, error) {
	if c.ConfigCellIncome != nil {
		return molecule.Bytes2GoU64(c.ConfigCellIncome.BasicCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellIncome is nil")
}

func (c *ConfigCellDataBuilder) IncomeMinTransferCapacity() (uint64, error) {
	if c.ConfigCellIncome != nil {
		return molecule.Bytes2GoU64(c.ConfigCellIncome.MinTransferCapacity().RawData())
	}
	return 0, fmt.Errorf("ConfigCellIncome is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateChannel() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.Channel().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateProposalCreate() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.ProposalCreate().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateProposalConfirm() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.ProposalConfirm().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateIncomeConsolidate() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.IncomeConsolidate().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateSaleBuyerInviter() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.SaleBuyerInviter().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateSaleBuyerChannel() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.SaleBuyerChannel().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateSaleDas() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.SaleDas().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) ProfitRateInviter() (uint32, error) {
	if c.ConfigCellProfitRate != nil {
		return molecule.Bytes2GoU32(c.ConfigCellProfitRate.Inviter().RawData())
	}
	return 0, fmt.Errorf("ConfigCellProfitRate is nil")
}

func (c *ConfigCellDataBuilder) Status() (uint8, error) {
	if c.ConfigCellMain != nil {
		return molecule.Bytes2GoU8(c.ConfigCellMain.Status().RawData())
	}
	return 0, fmt.Errorf("ConfigCellMain is nil")
}
