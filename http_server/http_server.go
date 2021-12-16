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
	Address        string
	AddressIndexer string
	AddressReverse string
	H              *handle.HttpHandle

	engine        *gin.Engine
	engineIndexer *gin.Engine
	engineReverse *gin.Engine

	srv        *http.Server
	srvIndexer *http.Server
	srvReverse *http.Server
}

func (h *HttpServer) Run() {
	h.engine = gin.New()
	h.engineIndexer = gin.New()
	h.engineReverse = gin.New()

	h.initRouter()

	h.srv = &http.Server{
		Addr:    h.Address,
		Handler: h.engine,
	}
	h.srvIndexer = &http.Server{
		Addr:    h.AddressIndexer,
		Handler: h.engineIndexer,
	}
	h.srvReverse = &http.Server{
		Addr:    h.AddressReverse,
		Handler: h.engineReverse,
	}

	go func() {
		if err := h.srv.ListenAndServe(); err != nil {
			log.Error("http_server old api run err:", err)
		}
	}()

	go func() {
		if err := h.srvIndexer.ListenAndServe(); err != nil {
			log.Error("http_server indexer api run err:", err)
		}
	}()

	go func() {
		if err := h.srvReverse.ListenAndServe(); err != nil {
			log.Error("http_server reverse api run err:", err)
		}
	}()
}

func (h *HttpServer) Shutdown() {
	if h.srv != nil {
		log.Warn("http server Shutdown ... ")
		if err := h.srv.Shutdown(h.Ctx); err != nil {
			log.Error("http server Shutdown err:", err.Error())
		}
		if err := h.srvIndexer.Shutdown(h.Ctx); err != nil {
			log.Error("http server Shutdown err:", err.Error())
		}
		if err := h.srvReverse.Shutdown(h.Ctx); err != nil {
			log.Error("http server Shutdown err:", err.Error())
		}
	}
}
