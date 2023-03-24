package handle

import (
	"das-account-indexer/http_server/code"
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
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
	ErrMsg       string `json:"err_msg"`
}

func (h *HttpHandle) JsonRpcBatchReverseRecord(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqBatchReverseRecord
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is :", len(req))
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
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
		apiResp  code.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req))

	if err = h.doBatchReverseRecord(&req, &apiResp); err != nil {
		log.Error("doBatchReverseRecord err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doBatchReverseRecord(req *ReqBatchReverseRecord, apiResp *code.ApiResp) error {
	var resp RespBatchReverseRecord
	resp.List = make([]BatchReverseRecord, 0)

	if count := len(req.BatchKeyInfo); count == 0 || count > 100 {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "Invalid number of key info")
		return nil
	}

	// check params
	var listKeyInfo []*core.DasAddressHex
	for i, _ := range req.BatchKeyInfo {
		res := checkReqKeyInfo(h.DasCore.Daf(), &req.BatchKeyInfo[i], apiResp)
		if apiResp.ErrNo != code.ApiCodeSuccess {
			log.Error("checkReqReverseRecord:", apiResp.ErrMsg)
			return nil
		}
		listKeyInfo = append(listKeyInfo, res)
	}

	// get reverse
	for i := range listKeyInfo {
		account, errMsg := h.checkReverse(listKeyInfo[i], apiResp)
		if apiResp.ErrNo != code.ApiCodeSuccess {
			return nil
		}
		resp.List = append(resp.List, BatchReverseRecord{
			Account:      account,
			AccountAlias: FormatDotToSharp(account),
			ErrMsg:       errMsg,
		})
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func (h *HttpHandle) checkReverse(dasAddrHex *core.DasAddressHex, apiResp *code.ApiResp) (account, errMsg string) {
	reverse, err := h.DbDao.FindLatestReverseRecord(dasAddrHex.ChainType, dasAddrHex.AddressHex)
	if err != nil {
		log.Error("FindLatestReverseRecord err: ", err.Error(), dasAddrHex.AddressHex)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record err")
		return
	} else if reverse.Id == 0 {
		errMsg = "reverse does not exit"
		return
	}

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(reverse.Account))
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		log.Error("FindAccountInfoByAccountId err: ", err.Error(), reverse.Account)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record err")
		return
	} else if accountInfo.Id == 0 {
		errMsg = fmt.Sprintf("reverse account[%s] does not exit", reverse.Account)
		return
	} else if accountInfo.Status == tables.AccountStatusOnLock {
		errMsg = fmt.Sprintf("reverse account[%s] cross-chain", reverse.Account)
		return
	}

	if accountInfo.OwnerChainType == dasAddrHex.ChainType && strings.EqualFold(accountInfo.Owner, dasAddrHex.AddressHex) {
		account = accountInfo.Account
	} else if accountInfo.ManagerChainType == dasAddrHex.ChainType && strings.EqualFold(accountInfo.Manager, dasAddrHex.AddressHex) {
		account = accountInfo.Account
	} else {
		addrNormal, err := h.DasCore.Daf().HexToNormal(*dasAddrHex)
		if err != nil {
			apiResp.ApiRespErr(code.ApiCodeParamsInvalid, err.Error())
			return
		}
		record, err := h.DbDao.FindRecordByAccountIdAddressValue(accountInfo.AccountId, addrNormal.AddressNormal)
		if err != nil {
			log.Error("FindRecordByAccountIdAddressValue err: ", err.Error(), accountInfo.Account, addrNormal.AddressNormal)
			apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse account record err")
			return
		} else if record.Id > 0 {
			account = accountInfo.Account
		}
	}
	return
}
