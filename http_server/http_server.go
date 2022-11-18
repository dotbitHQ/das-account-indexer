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
	Ctx context.Context
	//Address        string
	AddressIndexer string
	//AddressReverse string
	H *handle.HttpHandle

	//engine        *gin.Engine
	engineIndexer *gin.Engine
	//engineReverse *gin.Engine

	//srv        *http.Server
	srvIndexer *http.Server
	//srvReverse *http.Server
}

func (h *HttpServer) Run() {
	//if h.Address != "" {
	//	h.engine = gin.New()
	//}
	if h.AddressIndexer != "" {
		h.engineIndexer = gin.New()
	}
	//if h.AddressReverse != "" {
	//	h.engineReverse = gin.New()
	//}

	h.initRouter()

	//if h.Address != "" {
	//	h.srv = &http.Server{
	//		Addr:    h.Address,
	//		Handler: h.engine,
	//	}
	//	go func() {
	//		if err := h.srv.ListenAndServe(); err != nil {
	//			log.Error("http_server old api run err:", err)
	//		}
	//	}()
	//}

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

	//if h.AddressReverse != "" {
	//	h.srvReverse = &http.Server{
	//		Addr:    h.AddressReverse,
	//		Handler: h.engineReverse,
	//	}
	//
	//	go func() {
	//		if err := h.srvReverse.ListenAndServe(); err != nil {
	//			log.Error("http_server reverse api run err:", err)
	//		}
	//	}()
	//}

}

func (h *HttpServer) Shutdown() {
	log.Warn("http server Shutdown ... ")
	//if err := h.srv.Shutdown(h.Ctx); err != nil {
	//	log.Error("http server Shutdown err:", err.Error())
	//}

	if h.srvIndexer != nil {
		if err := h.srvIndexer.Shutdown(h.Ctx); err != nil {
			log.Error("http server Shutdown err:", err.Error())
		}
	}

	//if err := h.srvReverse.Shutdown(h.Ctx); err != nil {
	//	log.Error("http server Shutdown err:", err.Error())
	//}
}
