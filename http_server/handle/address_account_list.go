package handle

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqAddressAccountList struct {
	ReqKeyInfo
}

type RespAddressAccountList struct {
	AccountList []string `json:"account_list"`
}

func (h *HttpHandle) JsonRpcAddressAccountList(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqAddressAccountList
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

	if err = h.doAddressAccountList(&req[0], apiResp); err != nil {
		log.Error("doAddressAccountList err:", err.Error())
	}
}

func (h *HttpHandle) AddressAccountList(ctx *gin.Context) {
	var (
		funcName = "AddressAccountList"
		req      ReqAddressAccountList
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

	if err = h.doAddressAccountList(&req, &apiResp); err != nil {
		log.Error("doAddressAccountList err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAddressAccountList(req *ReqAddressAccountList, apiResp *code.ApiResp) error {
	var resp RespAddressAccountList
	resp.AccountList = make([]string, 0)

	res := checkReqKeyInfo(&req.ReqKeyInfo, apiResp)
	if apiResp.ErrNo != code.ApiCodeSuccess {
		log.Error("checkReqReverseRecord:", apiResp.ErrMsg)
		return nil
	}

	log.Info("doAddressAccountList:", res.ChainType, res.Address)

	list, err := h.DbDao.FindAccountNameListByAddress(res.ChainType, res.Address)
	if err != nil {
		log.Error("FindAccountListByAddress err:", err.Error(), req.KeyInfo)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find account list err")
		return fmt.Errorf("FindAccountListByAddress err: %s", err.Error())
	}

	for _, v := range list {
		resp.AccountList = append(resp.AccountList, v.Account)
	}

	apiResp.ApiRespOK(resp)
	return nil
}
