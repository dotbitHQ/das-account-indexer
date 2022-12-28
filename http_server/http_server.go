package http_server

import (
	"context"
	"das-account-indexer/http_server/handle"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/mylog"
	"net/http"
)

var (
	log = mylog.NewLogger("http_server", mylog.LevelDebug)
)

type HttpServer struct {
	Ctx            context.Context
	AddressIndexer string
	H              *handle.HttpHandle

	engineIndexer *gin.Engine
	srvIndexer    *http.Server
}

func (h *HttpServer) Run() {
	if h.AddressIndexer != "" {
		h.engineIndexer = gin.New()
	}

	h.initRouter()

	if h.AddressIndexer != "" {
		h.srvIndexer = &http.Server{
			Addr:    h.AddressIndexer,
			Handler: h.engineIndexer,
		}
		go func() {
			if err := h.srvIndexer.ListenAndServe(); err != nil {
				log.Error("http_server indexer api run err:", err)
			}
		}()
	}

}

func (h *HttpServer) Shutdown() {
	log.Warn("http server Shutdown ... ")
	if h.srvIndexer != nil {
		if err := h.srvIndexer.Shutdown(h.Ctx); err != nil {
			log.Error("http server Shutdown err:", err.Error())
		}
	}
}
