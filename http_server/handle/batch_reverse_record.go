package handle

import (
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

type ReqBatchReverseRecord struct {
	BatchKeyInfo []core.ChainTypeAddress `json:"batch_key_info"`
}

type RespBatchReverseRecord struct {
	List []BatchReverseRecord `json:"list"`
}

type BatchReverseRecord struct {
	Account      string `json:"account"`
	AccountAlias string `json:"account_alias"`
	DisplayName  string `json:"display_name"`
	ErrMsg       string `json:"err_msg"`
}

func (h *HttpHandle) JsonRpcBatchReverseRecord(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqBatchReverseRecord
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

	if err = h.doBatchReverseRecord(&req[0], apiResp); err != nil {
		log.Error("doBatchReverseRecord err:", err.Error())
	}
}

func (h *HttpHandle) BatchReverseRecord(ctx *gin.Context) {
	var (
		funcName = "BatchReverseRecord"
		req      ReqBatchReverseRecord
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

	if err = h.doBatchReverseRecord(&req, &apiResp); err != nil {
		log.Error("doBatchReverseRecord err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doBatchReverseRecord(req *ReqBatchReverseRecord, apiResp *http_api.ApiResp) error {
	var resp RespBatchReverseRecord
	resp.List = make([]BatchReverseRecord, 0)

	if count := len(req.BatchKeyInfo); count == 0 || count > 100 {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "Invalid number of key info")
		return nil
	}

	// check params
	var listKeyInfo []*core.DasAddressHex
	for i, _ := range req.BatchKeyInfo {
		res := checkReqKeyInfo(h.DasCore.Daf(), &req.BatchKeyInfo[i], apiResp)
		if apiResp.ErrNo != http_api.ApiCodeSuccess {
			log.Error("checkReqReverseRecord:", apiResp.ErrMsg)
			return nil
		}
		listKeyInfo = append(listKeyInfo, res)
	}

	// get reverse
	for _, v := range listKeyInfo {
		account, errMsg := h.checkReverse(v.ChainType, v.AddressHex, apiResp)
		if apiResp.ErrNo != http_api.ApiCodeSuccess {
			return nil
		}
		resp.List = append(resp.List, BatchReverseRecord{
			Account:      account,
			AccountAlias: FormatDotToSharp(account),
			DisplayName:  FormatDisplayName(account),
			ErrMsg:       errMsg,
		})
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func (h *HttpHandle) checkReverse(chainType common.ChainType, addressHex string, apiResp *http_api.ApiResp) (account, errMsg string) {
	reverse, err := h.DbDao.FindLatestReverseRecord(chainType, addressHex, "")
	if err != nil {
		log.Error("FindLatestReverseRecord err: ", err.Error(), addressHex)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find reverse record err")
		return
	} else if reverse.Id == 0 {
		errMsg = "reverse does not exit"
		return
	}

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
	}

	if accountInfo.OwnerChainType == chainType && strings.EqualFold(accountInfo.Owner, addressHex) {
		account = accountInfo.Account
	} else if accountInfo.ManagerChainType == chainType && strings.EqualFold(accountInfo.Manager, addressHex) {
		account = accountInfo.Account
	} else {
		record, err := h.DbDao.FindRecordByAccountIdAddressValue(accountInfo.AccountId, addressHex)
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
