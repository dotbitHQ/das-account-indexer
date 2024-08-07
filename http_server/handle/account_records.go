package handle

import (
	"context"
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

type ReqAccountRecords struct {
	Account string `json:"account"`
}

type RespAccountRecords struct {
	Account string       `json:"account"`
	Records []DataRecord `json:"records"`
}

type DataRecord struct {
	Key   string `json:"key"`
	Label string `json:"label"`
	Value string `json:"value"`
	TTL   string `json:"ttl"`
}

func (h *HttpHandle) JsonRpcAccountRecords(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqAccountRecords
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

	if err = h.doAccountRecords(h.Ctx, &req[0], apiResp, common.ConvertRecordsAddressCoinType); err != nil {
		log.Error("doAccountRecords err:", err.Error())
	}
}

func (h *HttpHandle) AccountRecords(ctx *gin.Context) {
	var (
		funcName = "AccountRecords"
		req      ReqAccountRecords
		apiResp  http_api.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, ctx.Request.Context())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", ctx.Request.Host, funcName, clientIp, toolib.JsonString(req), ctx.Request.Context())

	if err = h.doAccountRecords(ctx.Request.Context(), &req, &apiResp, common.ConvertRecordsAddressCoinType); err != nil {
		log.Error("doAccountRecords err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

type ConvertRecordsFunc func(string) string

func (h *HttpHandle) doAccountRecords(ctx context.Context, req *ReqAccountRecords, apiResp *http_api.ApiResp, convertRecordsFunc ConvertRecordsFunc) error {
	var resp RespAccountRecords
	resp.Records = make([]DataRecord, 0)

	req.Account = strings.TrimSpace(req.Account)
	req.Account = FormatSharpToDot(req.Account)
	if err := checkAccount(req.Account, apiResp); err != nil {
		log.Error(ctx, "checkAccount err: ", err.Error())
		return nil
	}

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		log.Error(ctx, "FindAccountInfoByAccountName err:", err.Error(), req.Account)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find account info err")
		return nil
	} else if accountInfo.Id == 0 {
		apiResp.ApiRespErr(http_api.ApiCodeAccountNotExist, "account not exist")
		return nil
	} else if accountInfo.Status == tables.AccountStatusOnLock {
		apiResp.ApiRespErr(http_api.ApiCodeAccountOnLock, "account cross-chain")
		return nil
	}

	resp.Account = req.Account

	list, err := h.DbDao.FindAccountRecordsByAccountId(accountId)
	if err != nil {
		log.Error(ctx, "FindAccountRecords err:", err.Error(), req.Account)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find records info err")
		return nil
	}
	for _, v := range list {
		key := fmt.Sprintf("%s.%s", v.Type, v.Key)
		resp.Records = append(resp.Records, DataRecord{
			Key:   convertRecordsFunc(key),
			Label: v.Label,
			Value: v.Value,
			TTL:   v.Ttl,
		})
	}

	apiResp.ApiRespOK(resp)
	return nil
}
