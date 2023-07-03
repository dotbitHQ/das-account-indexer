package handle

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqSubAccountVerify struct {
	Account    string `json:"account"`
	Address    string `json:"address"`
	VerifyType uint   `json:"verify_type"`
}

type RespSubAccountVerify struct {
	IsSubdid bool `json:"is_subdid"`
}

func (h *HttpHandle) JsonRpcSubAccountVerify(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqSubAccountVerify
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is:", len(req))
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if err = h.doSubAccountVerify(&req[0], apiResp); err != nil {
		log.Error("doSubAccountVerify err:", err.Error())
	}
}

func (h *HttpHandle) SubAccountVerify(ctx *gin.Context) {
	var (
		funcName = "SubAccountList"
		req      ReqSubAccountVerify
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

	if err = h.doSubAccountVerify(&req, &apiResp); err != nil {
		log.Error("doSubAccountList err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doSubAccountVerify(req *ReqSubAccountVerify, apiResp *code.ApiResp) error {
	var resp RespSubAccountVerify
	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	addrHex, err := formatAddress(h.DasCore.Daf(), req.Address)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, err.Error())
		return fmt.Errorf("formatAddress err: %s", err.Error())
	}
	log.Info("formatAddress:", req.Address, addrHex.ChainType, addrHex.AddressHex)

	res, err := h.DbDao.GetSubAccByParentAccountIdOfAddress(accountId, addrHex.AddressHex, req.VerifyType)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, "find account info err")
		return fmt.Errorf("GetSubAccByParentAccountIdOfAddress err: %s", err.Error())
	}
	if res > 0 {
		resp.IsSubdid = true
	}
	apiResp.ApiRespOK(resp)
	return nil
}
