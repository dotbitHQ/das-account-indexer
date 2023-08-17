package handle

import (
	"das-account-indexer/http_server/code"
	"fmt"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

func (h *HttpHandle) Query(ctx *gin.Context) {
	var (
		req      http_api.JsonRequest
		resp     http_api.JsonResponse
		apiResp  http_api.ApiResp
		clientIp = GetClientIp(ctx)
	)
	resp.Result = &apiResp

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Error("ShouldBindJSON err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.ID, resp.JsonRpc = req.ID, req.JsonRpc
	log.Info("Query:", req.Method, clientIp, toolib.JsonString(req))

	switch req.Method {
	case code.MethodServerInfo:
		h.JsonRpcServerInfo(req.Params, &apiResp)
	case code.MethodSearchAccount:
		h.JsonRpcSearchAccount(req.Params, &apiResp)
	case code.MethodAddressAccount:
		h.JsonRpcAddressAccount(req.Params, &apiResp)
	default:
		log.Error("method not exist:", req.Method)
		apiResp.ApiRespErr(http_api.ApiCodeMethodNotExist, fmt.Sprintf("method [%s] not exits", req.Method))
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

func (h *HttpHandle) QueryIndexer(ctx *gin.Context) {
	var (
		req      http_api.JsonRequest
		resp     http_api.JsonResponse
		apiResp  http_api.ApiResp
		clientIp = GetClientIp(ctx)
	)
	resp.Result = &apiResp

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Error("ShouldBindJSON err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.ID, resp.JsonRpc = req.ID, req.JsonRpc
	log.Info("QueryIndexer:", req.Method, clientIp, toolib.JsonString(req))

	switch req.Method {
	case code.MethodSearchAccount:
		h.JsonRpcSearchAccount(req.Params, &apiResp)
	case code.MethodAddressAccount:
		h.JsonRpcAddressAccount(req.Params, &apiResp)
	case code.MethodServerInfo:
		h.JsonRpcServerInfo(req.Params, &apiResp)
	case code.MethodAccountInfo:
		h.JsonRpcAccountInfo(req.Params, &apiResp)
	case code.MethodAccountList:
		h.JsonRpcAccountList(req.Params, &apiResp)
	case code.MethodAccountRecords:
		h.JsonRpcAccountRecords(req.Params, &apiResp)
	case code.MethodAccountReverseAddress:
		h.JsonRpcAccountReverseAddress(req.Params, &apiResp)
	case code.MethodBatchAccountRecords:
		h.JsonRpcBatchAccountRecords(req.Params, &apiResp)
	case code.MethodAccountRecordsV2:
		h.JsonRpcAccountRecordsV2(req.Params, &apiResp)
	case code.MethodReverseRecord:
		h.JsonRpcReverseRecord(req.Params, &apiResp)
	case code.MethodBatchReverseRecord:
		h.JsonRpcBatchReverseRecord(req.Params, &apiResp)
	case code.MethodSubAccountList:
		h.JsonRpcSubAccountList(req.Params, &apiResp)
	case code.MethodSubAccountVerify:
		h.JsonRpcSubAccountVerify(req.Params, &apiResp)
	default:
		log.Error("method not exist:", req.Method)
		apiResp.ApiRespErr(http_api.ApiCodeMethodNotExist, fmt.Sprintf("method [%s] not exits", req.Method))
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

func (h *HttpHandle) QueryReverse(ctx *gin.Context) {
	var (
		req      http_api.JsonRequest
		resp     http_api.JsonResponse
		apiResp  http_api.ApiResp
		clientIp = GetClientIp(ctx)
	)
	resp.Result = &apiResp

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Error("ShouldBindJSON err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.ID, resp.JsonRpc = req.ID, req.JsonRpc
	log.Info("QueryReverse:", req.Method, clientIp, toolib.JsonString(req))

	switch req.Method {
	case code.MethodServerInfo:
		h.JsonRpcServerInfo(req.Params, &apiResp)
	case code.MethodAccountInfo:
		h.JsonRpcAccountInfo(req.Params, &apiResp)
	case code.MethodAccountList:
		h.JsonRpcAccountList(req.Params, &apiResp)
	case code.MethodReverseRecord:
		h.JsonRpcReverseRecord(req.Params, &apiResp)
	default:
		log.Error("method not exist:", req.Method)
		apiResp.ApiRespErr(http_api.ApiCodeMethodNotExist, fmt.Sprintf("method [%s] not exits", req.Method))
	}

	ctx.JSON(http.StatusOK, resp)
	return
}
