package handle

import (
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqReverseRecordV2 struct {
	core.ChainTypeAddress
}

type RespReverseRecordV2 struct {
	Account      string `json:"account"`
	AccountAlias string `json:"account_alias"`
	DisplayName  string `json:"display_name"`
}

func (h *HttpHandle) JsonRpcReverseRecordV2(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqReverseRecordV2
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is :", len(req))
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doReverseRecordV2(&req[0], apiResp); err != nil {
		log.Error("doReverseRecordV2 err:", err.Error())
	}
}

func (h *HttpHandle) ReverseRecordV2(ctx *gin.Context) {
	var (
		funcName = "ReverseRecordV2"
		req      ReqReverseRecordV2
		apiResp  http_api.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req))

	if err = h.doReverseRecordV2(&req, &apiResp); err != nil {
		log.Error("doReverseRecordV2 err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doReverseRecordV2(req *ReqReverseRecordV2, apiResp *http_api.ApiResp) error {
	var resp RespReverseRecordV2
	var chainType common.ChainType
	var addressHex string
	var btcAddr string

	addrHex, err := req.FormatChainTypeAddress(h.DasCore.NetType(), false)
	if err != nil {
		log.Error("FormatChainTypeAddress err:", req.KeyInfo.Key)
		apiResp.ApiRespOK(resp)
		return nil
	}

	switch addrHex.DasAlgorithmId {
	case common.DasAlgorithmIdAnyLock:
		res, err := addrHex.FormatAnyLock()
		if err != nil {
			log.Error("addrHex.FormatAnyLock err: %s", err.Error())
			apiResp.ApiRespOK(resp)
			return nil
		}
		chainType = res.ChainType
		addressHex = res.AddressHex
	case common.DasAlgorithmIdEth, common.DasAlgorithmIdTron,
		common.DasAlgorithmIdDogeChain, common.DasAlgorithmIdWebauthn:
		chainType = addrHex.ChainType
		addressHex = addrHex.AddressHex
	case common.DasAlgorithmIdBitcoin:
		log.Info("doReverseRecordV2:", addrHex.DasAlgorithmId, addrHex.DasSubAlgorithmId, addrHex.AddressHex)
		switch addrHex.DasSubAlgorithmId {
		case common.DasSubAlgorithmIdBitcoinP2PKH, common.DasSubAlgorithmIdBitcoinP2WPKH:
			chainType = addrHex.ChainType
			addressHex = addrHex.AddressHex
		default:
			chainType = addrHex.ChainType
			addressHex = addrHex.AddressHex
			btcAddr = req.KeyInfo.Key
		}
	default:
		log.Error("default address invalid")
		apiResp.ApiRespOK(resp)
		return nil
	}

	log.Info("doReverseRecordV2:", chainType, addressHex, req.KeyInfo.Key, btcAddr)

	// reverse
	reverse, err := h.DbDao.FindLatestReverseRecord(chainType, addressHex, btcAddr)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse record err")
		return fmt.Errorf("FindLatestReverseRecord err: %s", err.Error())
	} else if reverse.Id == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}
	if btcAddr != "" {
		addressHex = reverse.Address
	}

	// check account
	var owner, manager string
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(reverse.Account))
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse record account err")
		return fmt.Errorf("FindAccountInfoByAccountId err: %s", err.Error())
	} else if accountInfo.Id == 0 {
		apiResp.ApiRespOK(resp)
		return fmt.Errorf("account not exist")
	} else if accountInfo.Status == tables.AccountStatusOnLock {
		apiResp.ApiRespOK(resp)
		return fmt.Errorf("account on lock")
	} else if accountInfo.Status == tables.AccountStatusOnUpgrade {
		// did cell
		didAnyLock, err := h.getAnyLockAddressHex(accountId)
		if err != nil {
			log.Warn("getAnyLockAddressHex err: %s", err)
		} else {
			owner = didAnyLock.AddressHex
			manager = didAnyLock.AddressHex
		}
	} else {
		owner = accountInfo.Owner
		manager = accountInfo.Manager
	}
	log.Info("owner manager:", owner, manager)

	if strings.EqualFold(addressHex, owner) || strings.EqualFold(addressHex, manager) {
		resp.Account = accountInfo.Account
	} else {
		record, err := h.DbDao.FindRecordByAccountIdAddressValue(accountInfo.AccountId, req.KeyInfo.Key)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse record account record err")
			return fmt.Errorf("FindRecordByAccountIdAddressValue err: %s", err.Error())
		} else if record.Id > 0 {
			resp.Account = accountInfo.Account
		}
	}
	if resp.Account != "" {
		resp.AccountAlias = FormatDotToSharp(resp.Account)
		resp.DisplayName = FormatDisplayName(resp.Account)
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func (h *HttpHandle) getAnyLockAddressHex(accountId string) (*core.DasAddressHex, error) {
	// did cell
	didCell, err := h.DbDao.GetDidCellByAccountId(accountId)
	if err != nil {
		return nil, fmt.Errorf("GetDidCellInfoByAccountId err: %s", err.Error())
	} else if didCell.Id == 0 {
		return nil, fmt.Errorf("didCell is nil")
	}
	didHex := core.DasAddressHex{
		DasAlgorithmId: common.DasAlgorithmIdAnyLock,
		ParsedAddress: &address.ParsedAddress{
			Script: &types.Script{
				CodeHash: types.HexToHash(didCell.LockCodeHash),
				HashType: types.HashTypeType,
				Args:     common.Hex2Bytes(didCell.Args),
			},
		},
	}
	didAnyLock, err := didHex.FormatAnyLock()
	if err != nil {
		return nil, fmt.Errorf("FormatAnyLock err: %s", err.Error())
	}
	return didAnyLock, nil
}
