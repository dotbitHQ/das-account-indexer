package handle

import (
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqBatchAccountRecords struct {
	Accounts []string `json:"accounts"`
}

type RespBatchAccountRecords struct {
	List []BatchAccountRecord `json:"list"`
}

type BatchAccountRecord struct {
	Account   string       `json:"account"`
	AccountId string       `json:"account_id"`
	Records   []DataRecord `json:"records"`
	ErrMsg    string       `json:"err_msg"`
}

func (h *HttpHandle) JsonRpcBatchAccountRecords(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqBatchAccountRecords
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

	if err = h.doBatchAccountRecords(&req[0], apiResp); err != nil {
		log.Error("doBatchAccountRecords err:", err.Error())
	}
}

func (h *HttpHandle) BatchAccountRecords(ctx *gin.Context) {
	var (
		funcName = "BatchAccountRecords"
		req      ReqBatchAccountRecords
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

	if err = h.doBatchAccountRecords(&req, &apiResp); err != nil {
		log.Error("doBatchAccountRecords err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doBatchAccountRecords(req *ReqBatchAccountRecords, apiResp *http_api.ApiResp) error {
	var resp RespBatchAccountRecords
	resp.List = make([]BatchAccountRecord, 0)

	if count := len(req.Accounts); count == 0 || count > 100 {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "Invalid number of accounts")
		return nil
	}

	var accountIds []string
	for i, v := range req.Accounts {
		req.Accounts[i] = strings.TrimSpace(req.Accounts[i])
		req.Accounts[i] = FormatSharpToDot(req.Accounts[i])
		if err := checkAccount(req.Accounts[i], apiResp); err != nil {
			log.Error("checkAccount err: ", err.Error(), v)
			return nil
		}
		accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Accounts[i]))
		accountIds = append(accountIds, accountId)
		resp.List = append(resp.List, BatchAccountRecord{
			Account:   req.Accounts[i],
			AccountId: accountId,
			Records:   make([]DataRecord, 0),
		})
	}

	// get accounts
	list, err := h.DbDao.FindAccountInfoListByAccountIds(accountIds)
	if err != nil {
		log.Error("FindAccountInfoListByAccountIds err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find accounts err")
		return nil
	}

	var mapAcc = make(map[string]tables.TableAccountInfo)
	for i, v := range list {
		mapAcc[v.AccountId] = list[i]
	}

	// check accounts
	var okIds []string
	for i, v := range resp.List {
		acc, ok := mapAcc[v.AccountId]
		if !ok {
			resp.List[i].ErrMsg = fmt.Sprintf("account[%s] does not exist", v.Account)
		} else if acc.Status == tables.AccountStatusOnLock {
			resp.List[i].ErrMsg = fmt.Sprintf("account[%s] cross-chain", v.Account)
		} else {
			okIds = append(okIds, v.AccountId)
		}
	}
	if len(okIds) == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}

	// get records
	records, err := h.DbDao.FindRecordsByAccountIds(okIds)
	if err != nil {
		log.Error("FindRecordsByAccountIds err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find records err")
		return nil
	}
	var mapRecords = make(map[string][]DataRecord)
	for _, v := range records {
		key := fmt.Sprintf("%s.%s", v.Type, v.Key)
		mapRecords[v.AccountId] = append(mapRecords[v.AccountId], DataRecord{
			Key:   common.ConvertRecordsAddressKey(key),
			Label: v.Label,
			Value: v.Value,
			TTL:   v.Ttl,
		})
	}
	for i, v := range resp.List {
		if rs, ok := mapRecords[v.AccountId]; ok {
			resp.List[i].Records = rs
		}
	}

	apiResp.ApiRespOK(resp)
	return nil
}
