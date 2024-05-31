package handle

import (
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

	if err = h.doBatchReverseRecordV2(&req[0], apiResp); err != nil {
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
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req))

	if err = h.doBatchReverseRecordV2(&req, &apiResp); err != nil {
		log.Error("doBatchReverseRecordV2 err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doBatchReverseRecordV2(req *ReqBatchReverseRecordV2, apiResp *http_api.ApiResp) error {
	var resp RespBatchReverseRecordV2

	resp.List = make([]BatchReverseRecordV2, 0)

	if count := len(req.BatchKeyInfo); count == 0 || count > 100 {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "Invalid number of key info")
		return nil
	}

	// check params
	var listKeyInfo []*core.DasAddressHex
	for _, v := range req.BatchKeyInfo {
		addrHex, err := v.FormatChainTypeAddress(config.Cfg.Server.Net, false)
		if err != nil {
			log.Warn("FormatChainTypeAddress err: %s", err.Error())
			listKeyInfo = append(listKeyInfo, nil)
			continue
		} else if addrHex.DasAlgorithmId == common.DasAlgorithmIdAnyLock {
			anyLockAddrHex, err := addrHex.FormatAnyLock()
			if err != nil {
				log.Warn("FormatAnyLock err: %s", err.Error())
				listKeyInfo = append(listKeyInfo, nil)
				continue
			}
			listKeyInfo = append(listKeyInfo, anyLockAddrHex)
		} else {
			listKeyInfo = append(listKeyInfo, addrHex)
		}
	}

	// get reverse
	for i, v := range listKeyInfo {
		var tmp BatchReverseRecordV2
		if v == nil {
			tmp.ErrMsg = "address is invalid"
		} else {
			account, errMsg := h.checkReverseV2(v.ChainType, v.AddressHex, req.BatchKeyInfo[i].KeyInfo.Key, apiResp)
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

func (h *HttpHandle) checkReverseV2(chainType common.ChainType, addressHex, reqKey string, apiResp *http_api.ApiResp) (account, errMsg string) {
	// reverse
	reverse, err := h.DbDao.FindLatestReverseRecord(chainType, addressHex)
	if err != nil {
		log.Error("FindAccountInfoByAccountId err: ", err.Error(), reverse.Account)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse record err")
		return
	} else if reverse.Id == 0 {
		errMsg = "reverse does not exit"
		return
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
	log.Info("owner manager:", owner, manager)

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
