package handle

import (
	"context"
	"das-account-indexer/config"
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqBatchReverseRecordV2 struct {
	BatchKeyInfo []core.ChainTypeAddress `json:"batch_key_info"`
}

type RespBatchReverseRecordV2 struct {
	List []BatchReverseRecordV2 `json:"list"`
}

type BatchReverseRecordV2 struct {
	Account      string `json:"account"`
	AccountAlias string `json:"account_alias"`
	DisplayName  string `json:"display_name"`
	ErrMsg       string `json:"err_msg"`
}

func (h *HttpHandle) JsonRpcBatchReverseRecordV2(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqBatchReverseRecordV2
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

	if err = h.doBatchReverseRecordV2(h.Ctx, &req[0], apiResp); err != nil {
		log.Error("doBatchReverseRecordV2 err:", err.Error())
	}
}

func (h *HttpHandle) BatchReverseRecordV2(ctx *gin.Context) {
	var (
		funcName = "BatchReverseRecordV2"
		req      ReqBatchReverseRecordV2
		apiResp  http_api.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, ctx.Request.Context())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req), ctx.Request.Context())

	if err = h.doBatchReverseRecordV2(ctx.Request.Context(), &req, &apiResp); err != nil {
		log.Error("doBatchReverseRecordV2 err:", err.Error(), funcName, ctx.Request.Context())
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doBatchReverseRecordV2(ctx context.Context, req *ReqBatchReverseRecordV2, apiResp *http_api.ApiResp) error {
	var resp RespBatchReverseRecordV2

	resp.List = make([]BatchReverseRecordV2, 0)

	if count := len(req.BatchKeyInfo); count == 0 || count > 100 {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "Invalid number of key info")
		return nil
	}

	// check params
	var listKeyInfo []*core.DasAddressHex
	var listBtcAddr []string
	for _, v := range req.BatchKeyInfo {
		addrHex, err := v.FormatChainTypeAddress(config.Cfg.Server.Net, false)
		if err != nil {
			log.Warn(ctx, "FormatChainTypeAddress err: %s", err.Error())
			listKeyInfo = append(listKeyInfo, nil)
			listBtcAddr = append(listBtcAddr, "")
			continue
		}
		switch addrHex.DasAlgorithmId {
		case common.DasAlgorithmIdAnyLock:
			anyLockAddrHex, err := addrHex.FormatAnyLock()
			if err != nil {
				log.Warn(ctx, "FormatAnyLock err: %s", err.Error())
				listKeyInfo = append(listKeyInfo, nil)
				listBtcAddr = append(listBtcAddr, "")
				continue
			}
			listKeyInfo = append(listKeyInfo, anyLockAddrHex)
			listBtcAddr = append(listBtcAddr, "")
		case common.DasAlgorithmIdEth, common.DasAlgorithmIdTron,
			common.DasAlgorithmIdDogeChain, common.DasAlgorithmIdWebauthn:
			listKeyInfo = append(listKeyInfo, addrHex)
			listBtcAddr = append(listBtcAddr, "")
		case common.DasAlgorithmIdBitcoin:
			log.Info(ctx, "doReverseInfoV2:", addrHex.DasAlgorithmId, addrHex.DasSubAlgorithmId, addrHex.AddressHex)
			switch addrHex.DasSubAlgorithmId {
			case common.DasSubAlgorithmIdBitcoinP2PKH, common.DasSubAlgorithmIdBitcoinP2WPKH:
				listKeyInfo = append(listKeyInfo, addrHex)
				listBtcAddr = append(listBtcAddr, "")
			default:
				listKeyInfo = append(listKeyInfo, addrHex)
				listBtcAddr = append(listBtcAddr, v.KeyInfo.Key)
			}
		default:
			listKeyInfo = append(listKeyInfo, nil)
			listBtcAddr = append(listBtcAddr, "")
		}
	}

	// get reverse
	for i, v := range listKeyInfo {
		var tmp BatchReverseRecordV2
		if v == nil {
			tmp.ErrMsg = "address is invalid"
		} else {
			account, errMsg := h.checkReverseV2(v.ChainType, v.AddressHex, req.BatchKeyInfo[i].KeyInfo.Key, listBtcAddr[i], apiResp)
			if apiResp.ErrNo != http_api.ApiCodeSuccess {
				return nil
			}
			tmp = BatchReverseRecordV2{
				Account:      account,
				AccountAlias: FormatDotToSharp(account),
				DisplayName:  FormatDisplayName(account),
				ErrMsg:       errMsg,
			}
		}
		resp.List = append(resp.List, tmp)
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func (h *HttpHandle) checkReverseV2(chainType common.ChainType, addressHex, reqKey, btcAddr string, apiResp *http_api.ApiResp) (account, errMsg string) {
	// reverse
	reverse, err := h.DbDao.FindLatestReverseRecord(chainType, addressHex, btcAddr)
	if err != nil {
		log.Error("FindAccountInfoByAccountId err: ", err.Error(), reverse.Account)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse record err")
		return
	} else if reverse.Id == 0 {
		errMsg = "reverse does not exit"
		return
	}

	if btcAddr != "" {
		addressHex = reverse.Address
	}

	// check account
	var owner, manager string
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(reverse.Account))
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		log.Error("FindAccountInfoByAccountId err: ", err.Error(), reverse.Account)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse record err")
		return
	} else if accountInfo.Id == 0 {
		errMsg = fmt.Sprintf("reverse account[%s] does not exit", reverse.Account)
		return
	} else if accountInfo.Status == tables.AccountStatusOnLock {
		errMsg = fmt.Sprintf("reverse account[%s] cross-chain", reverse.Account)
		return
	} else if accountInfo.Status == tables.AccountStatusOnUpgrade {
		// did cell
		didAnyLock, err := h.getAnyLockAddressHex(accountId)
		if err != nil {
			log.Warn("getAnyLockAddressHex err: %s", err.Error())
		} else {
			owner = didAnyLock.AddressHex
			manager = didAnyLock.AddressHex
		}
	} else {
		owner = accountInfo.Owner
		manager = accountInfo.Manager
	}
	log.Info("owner manager:", owner, manager, addressHex)

	if strings.EqualFold(addressHex, owner) || strings.EqualFold(addressHex, manager) {
		account = accountInfo.Account
	} else {
		record, err := h.DbDao.FindRecordByAccountIdAddressValue(accountInfo.AccountId, reqKey)
		if err != nil {
			log.Error("FindRecordByAccountIdAddressValue err: ", err.Error(), accountInfo.Account, addressHex)
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse account record err")
			return
		} else if record.Id > 0 {
			account = accountInfo.Account
		}
	}
	return
}
