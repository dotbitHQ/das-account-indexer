package handle

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqReverseRecord struct {
	DasType common.ChainType `json:"das_type"`
	Address string           `json:"address"`
}

type RespReverseRecord struct {
	Account string `json:"account"`
}

func (h *HttpHandle) JsonRpcReverseRecord(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqReverseRecord
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

	if err = h.doReverseRecord(&req[0], apiResp); err != nil {
		log.Error("doReverseRecord err:", err.Error())
	}
}

func (h *HttpHandle) ReverseRecord(ctx *gin.Context) {
	var (
		funcName = "ReverseRecord"
		req      ReqReverseRecord
		apiResp  code.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, toolib.JsonString(req))

	if err = h.doReverseRecord(&req, &apiResp); err != nil {
		log.Error("doReverseRecord err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doReverseRecord(req *ReqReverseRecord, apiResp *code.ApiResp) error {
	var resp RespReverseRecord
	req.Address = core.FormatAddressToHex(req.DasType, req.Address)

	reverse, err := h.DbDao.FindLatestReverseRecord(req.DasType, req.Address)
	if err != nil {
		log.Error("FindLatestReverseRecord err:", err.Error(), req.DasType, req.Address)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record err")
		return nil
	} else if reverse.Id == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}

	account, err := h.DbDao.FindAccountInfoByAccountName(reverse.Account)
	if err != nil {
		log.Error("FindAccountInfoByAccountName err:", err.Error(), req.DasType, req.Address, reverse.Account)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record account err")
		return nil
	}

	if account.OwnerChainType == req.DasType && strings.EqualFold(account.Owner, req.Address) {
		resp.Account = account.Account
	} else if account.ManagerChainType == req.DasType && strings.EqualFold(account.Manager, req.Address) {
		resp.Account = account.Account
	} else {
		record, err := h.DbDao.FindRecordByAccountAddressValue(account.Account, req.Address)
		if err != nil {
			log.Error("FindRecordByAccountAddressValue err:", err.Error(), req.DasType, req.Address, reverse.Account)
			apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record account record err")
			return nil
		} else if record.Id > 0 {
			resp.Account = account.Account
		}
	}

	apiResp.ApiRespOK(resp)
	return nil
}
