package http_server

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

func (h *HttpServer) initRouter() {
	shortDataTime, lockTime, shortExpireTime := time.Minute, time.Second*30, time.Second*5
	cacheHandle := toolib.MiddlewareCacheByRedis(h.H.Red, false, shortDataTime, lockTime, shortExpireTime, respHandle)

	if h.AddressIndexer != "" {
		// indexer api
		h.engineIndexer.Use(toolib.MiddlewareCors())
		h.engineIndexer.POST("", cacheHandle, h.H.QueryIndexer)
		v1Indexer := h.engineIndexer.Group("v1")
		{
			v1Indexer.POST("/search/account", code.DoMonitorLog(code.MethodSearchAccount), cacheHandle, h.H.SearchAccount)
			v1Indexer.POST("/address/account", code.DoMonitorLog(code.MethodAddressAccount), cacheHandle, h.H.AddressAccount)
			v1Indexer.POST("/server/info", code.DoMonitorLog(code.MethodServerInfo), cacheHandle, h.H.ServerInfo)
			v1Indexer.POST("/account/info", code.DoMonitorLog(code.MethodAccountInfo), cacheHandle, h.H.AccountInfo)
			v1Indexer.POST("/account/list", code.DoMonitorLog(code.MethodAccountList), cacheHandle, h.H.AccountList)
			v1Indexer.POST("/account/records", code.DoMonitorLog(code.MethodAccountRecords), cacheHandle, h.H.AccountRecords)
			v1Indexer.POST("/account/reverse/address", code.DoMonitorLog(code.MethodAccountReverseAddress), cacheHandle, h.H.AccountReverseAddress)
			v1Indexer.POST("/reverse/record", code.DoMonitorLog(code.MethodReverseRecord), cacheHandle, h.H.ReverseRecord)
			v1Indexer.POST("/sub/account/list", code.DoMonitorLog(code.MethodSubAccountList), cacheHandle, h.H.SubAccountList)
			v1Indexer.POST("/sub/account/verify", code.DoMonitorLog(code.MethodSubAccountVerify), cacheHandle, h.H.SubAccountVerify)

			v1Indexer.POST("/batch/account/records", code.DoMonitorLog(code.MethodBatchAccountRecords), cacheHandle, h.H.BatchAccountRecords)
			v1Indexer.POST("/batch/reverse/record", code.DoMonitorLog(code.MethodBatchReverseRecord), cacheHandle, h.H.BatchReverseRecord)
			v1Indexer.POST("/batch/register/info", code.DoMonitorLog(code.MethodBatchRegisterInfo), cacheHandle, h.H.BatchRegisterInfo)
		}
		v2Indexer := h.engineIndexer.Group("v2")
		{
			v2Indexer.POST("/account/records", code.DoMonitorLog(code.MethodAccountRecordsV2), cacheHandle, h.H.AccountRecordsV2)
		}
	}

}

func respHandle(c *gin.Context, res string, err error) {
	if err != nil {
		log.Error("respHandle err:", err.Error())
		c.AbortWithStatusJSON(http.StatusOK, http_api.ApiRespErr(http_api.ApiCodeError500, err.Error()))
	} else if res != "" {
		var respMap map[string]interface{}
		_ = json.Unmarshal([]byte(res), &respMap)
		c.AbortWithStatusJSON(http.StatusOK, respMap)
	}
}
