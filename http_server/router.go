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

	h.engine.POST("", cacheHandle, h.H.Query)

	v1 := h.engine.Group("v1")
	{
		v1.POST("/version", cacheHandle, h.H.Version)
		v1.POST("/search/account", cacheHandle, h.H.SearchAccount)
		v1.POST("/address/account", cacheHandle, h.H.AddressAccount)

		v1.POST("/indexer/info", h.H.IndexerInfo)
		v1.POST("/account/info", cacheHandle, h.H.AccountInfo)
		v1.POST("/account/records", cacheHandle, h.H.AccountRecords)
		v1.POST("/reverse/record", cacheHandle, h.H.ReverseRecord)
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
