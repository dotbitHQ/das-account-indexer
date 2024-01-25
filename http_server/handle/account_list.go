package handle

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqAccountList struct {
	core.ChainTypeAddress
	Role string `json:"role"` // owner,manager
}

type RespAccountList struct {
	AccountList []RespAddressAccount `json:"account_list"`
}

type RespAddressAccount struct {
	Account      string `json:"account"`
	AccountAlias string `json:"account_alias"`
	DisplayName  string `json:"display_name"`
	RegisteredAt uint64 `json:"registered_at"`
	ExpiredAt    uint64 `json:"expired_at"`
}

func (h *HttpHandle) JsonRpcAccountList(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqAccountList
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

	if err = h.doAccountList(&req[0], apiResp); err != nil {
		log.Error("doAccountList err:", err.Error())
	}
}

func (h *HttpHandle) AccountList(ctx *gin.Context) {
	var (
		funcName = "AccountList"
		req      ReqAccountList
		apiResp  code.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req))

	if err = h.doAccountList(&req, &apiResp); err != nil {
		log.Error("doAccountList err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAccountList(req *ReqAccountList, apiResp *code.ApiResp) error {
	var resp RespAccountList
	resp.AccountList = make([]RespAddressAccount, 0)

	res := checkReqKeyInfo(h.DasCore.Daf(), &req.ChainTypeAddress, apiResp)
	if apiResp.ErrNo != http_api.ApiCodeSuccess {
		log.Error("checkReqReverseRecord:", apiResp.ErrMsg)
		return nil
	}

	log.Info("doAccountList:", res.ChainType, res.AddressHex)

	list, err := h.DbDao.FindAccountNameListByAddress(res.ChainType, res.AddressHex, req.Role)
	if err != nil {
		log.Error("FindAccountListByAddress err:", err.Error(), req.KeyInfo)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find account list err")
		return fmt.Errorf("FindAccountListByAddress err: %s", err.Error())
	}

	for _, v := range list {
		accountAlias := FormatDotToSharp(v.Account)
		tmp := RespAddressAccount{
			Account:      v.Account,
			AccountAlias: accountAlias,
			DisplayName:  FormatDisplayName(v.Account),
			RegisteredAt: v.RegisteredAt,
			ExpiredAt:    v.ExpiredAt,
		}
		resp.AccountList = append(resp.AccountList, tmp)
	}

	apiResp.ApiRespOK(resp)
	return nil
}
