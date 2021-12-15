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
	Ctx     context.Context
	Address string
	H       *handle.HttpHandle

	engine *gin.Engine
	srv    *http.Server
}

func (h *HttpServer) Run() {
	h.engine = gin.New()
	h.initRouter()
	h.srv = &http.Server{
		Addr:    h.Address,
		Handler: h.engine,
	}
	go func() {
		if err := h.srv.ListenAndServe(); err != nil {
			log.Error("http_server run err:", err)
		}
	}()
}

func (h *HttpServer) Shutdown() {
	if h.srv != nil {
		log.Warn("http server Shutdown ... ")
		if err := h.srv.Shutdown(h.Ctx); err != nil {
			log.Error("http server Shutdown err:", err.Error())
		}
	}
}
