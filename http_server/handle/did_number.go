package handle

import (
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqDidNumber struct {
}

type RespDidNumber struct {
	Total int64 `json:"total"`
	TLDid int64 `json:"tl_did"`
	SLDid int64 `json:"sl_did"`
	Dobs  int64 `json:"dobs"`
}

func (h *HttpHandle) JsonRpcDidNumber(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqDidNumber
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

	if err = h.doDidNumber(&req[0], apiResp); err != nil {
		log.Error("doDidNumber err:", err.Error())
	}
}

func (h *HttpHandle) DidNumber(ctx *gin.Context) {
	var (
		funcName = "DidNumber"
		req      ReqDidNumber
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

	if err = h.doDidNumber(&req, &apiResp); err != nil {
		log.Error("doDidNumber err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doDidNumber(req *ReqDidNumber, apiResp *http_api.ApiResp) error {
	var resp RespDidNumber

	tldid, err := h.DbDao.GetTotalTLDid()
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
		return fmt.Errorf("GetTotalTLDid err: %s", err.Error())
	}
	sldid, err := h.DbDao.GetTotalSLDid()
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
		return fmt.Errorf("GetTotalSLDid err: %s", err.Error())
	}
	dobs, err := h.DbDao.GetTotalDobs()
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
		return fmt.Errorf("GetTotalDobs err: %s", err.Error())
	}

	resp.TLDid = tldid
	resp.SLDid = sldid
	resp.Total = tldid + sldid
	resp.Dobs = dobs

	apiResp.ApiRespOK(resp)
	return nil
}
