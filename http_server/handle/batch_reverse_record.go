package handle

import (
	"das-account-indexer/http_server/code"
	"das-account-indexer/tables"
	"encoding/json"
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
	Account string `json:"account"`
	ErrMsg  string `json:"err_msg"`
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

	if count := len(req.BatchKeyInfo); count == 0 || count > 20 {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "Invalid number of accounts")
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
	for _, v := range listKeyInfo {
		account, errMsg := h.checkReverse(v.ChainType, v.AddressHex)
		resp.List = append(resp.List, BatchReverseRecord{
			Account: account,
			ErrMsg:  errMsg,
		})
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func (h *HttpHandle) checkReverse(chainType common.ChainType, addressHex string) (account, errMsg string) {
	reverse, err := h.DbDao.FindLatestReverseRecord(chainType, addressHex)
	if err != nil {
		errMsg = "find reverse record err"
		return
	} else if reverse.Id == 0 {
		errMsg = "reverse does not exit"
		return
	}

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(reverse.Account))
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		errMsg = "find reverse record account err"
		return
	} else if accountInfo.Id == 0 {
		errMsg = "reverse account does not exit"
		return
	} else if accountInfo.Status == tables.AccountStatusOnLock {
		errMsg = "reverse account cross-chain"
		return
	}

	if accountInfo.OwnerChainType == chainType && strings.EqualFold(accountInfo.Owner, addressHex) {
		account = accountInfo.Account
	} else if accountInfo.ManagerChainType == chainType && strings.EqualFold(accountInfo.Manager, addressHex) {
		account = accountInfo.Account
	} else {
		record, err := h.DbDao.FindRecordByAccountIdAddressValue(accountInfo.AccountId, addressHex)
		if err != nil {
			account = ""
			errMsg = "find reverse account record err"
			return
		} else if record.Id > 0 {
			account = accountInfo.Account
		}
	}
	return
}
