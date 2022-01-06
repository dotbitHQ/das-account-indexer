package http_server

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"time"
)

func (h *HttpServer) initRouter() {
	shortDataTime, lockTime, shortExpireTime := time.Minute, time.Second*30, time.Second*5
	cacheHandle := toolib.MiddlewareCacheByRedis(h.H.Red, false, shortDataTime, lockTime, shortExpireTime, respHandle)

	if h.Address != "" {
		// old api
		h.engine.Use(toolib.MiddlewareCors())
		h.engine.POST("", cacheHandle, h.H.Query)
		v1 := h.engine.Group("v1")
		{
			v1.POST("/server/info", cacheHandle, h.H.ServerInfo)
			v1.POST("/search/account", cacheHandle, h.H.SearchAccount)
			v1.POST("/address/account", cacheHandle, h.H.AddressAccount)
		}
	}
	if h.AddressIndexer != "" {
		// indexer api
		h.engineIndexer.Use(toolib.MiddlewareCors())
		h.engineIndexer.POST("", cacheHandle, h.H.QueryIndexer)
		v1Indexer := h.engineIndexer.Group("v1")
		{
			v1Indexer.POST("/server/info", cacheHandle, h.H.ServerInfo)
			v1Indexer.POST("/account/info", cacheHandle, h.H.AccountInfo)
			v1Indexer.POST("/account/list", cacheHandle, h.H.AccountList)
			v1Indexer.POST("/account/records", cacheHandle, h.H.AccountRecords)
			v1Indexer.POST("/reverse/record", cacheHandle, h.H.ReverseRecord)

		}
	}

	if h.AddressReverse != "" {
		// reverse api
		h.engineReverse.Use(toolib.MiddlewareCors())
		h.engineReverse.POST("", cacheHandle, h.H.QueryReverse)
		v1Reverse := h.engineReverse.Group("v1")
		{
			v1Reverse.POST("/server/info", cacheHandle, h.H.ServerInfo)
			v1Reverse.POST("/reverse/record", cacheHandle, h.H.ReverseRecord)
			v1Reverse.POST("/account/info", cacheHandle, h.H.AccountInfo)
			v1Reverse.POST("/account/list", cacheHandle, h.H.AccountList)
		}
	}

}

func respHandle(c *gin.Context, res string, err error) {
	if err != nil {
		log.Error("respHandle err:", err.Error())
		c.AbortWithStatusJSON(http.StatusOK, code.ApiRespErr(code.ApiCodeError500, err.Error()))
	} else if res != "" {
		var respMap map[string]interface{}
		_ = json.Unmarshal([]byte(res), &respMap)
		c.AbortWithStatusJSON(http.StatusOK, respMap)
	}
}
