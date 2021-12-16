package handle

import (
	"das-account-indexer/http_server/code"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

func (h *HttpHandle) Query(ctx *gin.Context) {
	var (
		req      code.JsonRequest
		resp     code.JsonResponse
		apiResp  code.ApiResp
		clientIp = GetClientIp(ctx)
	)
	resp.Result = &apiResp

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Error("ShouldBindJSON err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
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
		apiResp.ApiRespErr(code.ApiCodeMethodNotExist, fmt.Sprintf("method [%s] not exits", req.Method))
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

func (h *HttpHandle) QueryIndexer(ctx *gin.Context) {
	var (
		req      code.JsonRequest
		resp     code.JsonResponse
		apiResp  code.ApiResp
		clientIp = GetClientIp(ctx)
	)
	resp.Result = &apiResp

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Error("ShouldBindJSON err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.ID, resp.JsonRpc = req.ID, req.JsonRpc
	log.Info("QueryIndexer:", req.Method, clientIp, toolib.JsonString(req))

	switch req.Method {
	case code.MethodServerInfo:
		h.JsonRpcServerInfo(req.Params, &apiResp)
	case code.MethodAccountInfo:
		h.JsonRpcAccountInfo(req.Params, &apiResp)
	case code.MethodAccountRecords:
		h.JsonRpcAccountRecords(req.Params, &apiResp)
	case code.MethodReverseRecord:
		h.JsonRpcReverseRecord(req.Params, &apiResp)
	default:
		log.Error("method not exist:", req.Method)
		apiResp.ApiRespErr(code.ApiCodeMethodNotExist, fmt.Sprintf("method [%s] not exits", req.Method))
	}

	ctx.JSON(http.StatusOK, resp)
	return
}

func (h *HttpHandle) QueryReverse(ctx *gin.Context) {
	var (
		req      code.JsonRequest
		resp     code.JsonResponse
		apiResp  code.ApiResp
		clientIp = GetClientIp(ctx)
	)
	resp.Result = &apiResp

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Error("ShouldBindJSON err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, resp)
		return
	}

	resp.ID, resp.JsonRpc = req.ID, req.JsonRpc
	log.Info("QueryReverse:", req.Method, clientIp, toolib.JsonString(req))

	switch req.Method {
	case code.MethodServerInfo:
		h.JsonRpcServerInfo(req.Params, &apiResp)
	case code.MethodReverseRecord:
		h.JsonRpcReverseRecord(req.Params, &apiResp)
	default:
		log.Error("method not exist:", req.Method)
		apiResp.ApiRespErr(code.ApiCodeMethodNotExist, fmt.Sprintf("method [%s] not exits", req.Method))
	}

	ctx.JSON(http.StatusOK, resp)
	return
}
