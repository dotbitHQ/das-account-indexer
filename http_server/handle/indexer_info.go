package handle

import (
	"das-account-indexer/block_parser"
	"das-account-indexer/http_server/code"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RespServerInfo struct {
	IsLatestBlockNumber bool   `json:"is_latest_block_number"`
	CurrentBlockNumber  uint64 `json:"current_block_number"`
}

func (h *HttpHandle) JsonRpcServerInfo(p json.RawMessage, apiResp *code.ApiResp) {
	if err := h.doServerInfo(apiResp); err != nil {
		log.Error("doServerInfo err:", err.Error())
	}
}

func (h *HttpHandle) ServerInfo(ctx *gin.Context) {
	var (
		funcName = "ServerInfo"
		apiResp  code.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	log.Info("ApiReq:", funcName, clientIp)

	if err = h.doServerInfo(&apiResp); err != nil {
		log.Error("doServerInfo err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doServerInfo(apiResp *code.ApiResp) error {
	var resp RespServerInfo

	resp.IsLatestBlockNumber = block_parser.IsLatestBlockNumber
	resp.CurrentBlockNumber = block_parser.CurrentBlockNumber

	apiResp.ApiRespOK(resp)
	return nil
}
