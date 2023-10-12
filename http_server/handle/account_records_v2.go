package handle

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

func (h *HttpHandle) JsonRpcAccountRecordsV2(p json.RawMessage, apiResp *code.ApiResp) {
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

	if err = h.doAccountRecords(&req[0], apiResp, common.ConvertRecordsAddressKey); err != nil {
		log.Error("doAccountRecords err:", err.Error())
	}
}

func (h *HttpHandle) AccountRecordsV2(ctx *gin.Context) {
	var (
		funcName = "AccountRecordsV2"
		req      ReqAccountRecords
		apiResp  code.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", ctx.Request.Host, funcName, clientIp, toolib.JsonString(req))

	if err = h.doAccountRecords(&req, &apiResp, common.ConvertRecordsAddressKey); err != nil {
		log.Error("doAccountRecords err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}
