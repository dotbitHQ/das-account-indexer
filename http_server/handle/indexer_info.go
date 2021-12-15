package handle

import (
	"das-account-indexer/block_parser"
	"das-account-indexer/http_server/code"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RespIndexerInfo struct {
	IsLatestBlockNumber bool   `json:"is_latest_block_number"`
	CurrentBlockNumber  uint64 `json:"current_block_number"`
}

func (h *HttpHandle) JsonRpcIndexerInfo(p json.RawMessage, apiResp *code.ApiResp) {
	if err := h.doIndexerInfo(apiResp); err != nil {
		log.Error("doIndexerInfo err:", err.Error())
	}
}

func (h *HttpHandle) IndexerInfo(ctx *gin.Context) {
	var (
		funcName = "IndexerInfo"
		apiResp  code.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	log.Info("ApiReq:", funcName, clientIp)

	if err = h.doIndexerInfo(&apiResp); err != nil {
		log.Error("doIndexerInfo err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doIndexerInfo(apiResp *code.ApiResp) error {
	var resp RespIndexerInfo

	resp.IsLatestBlockNumber = block_parser.IsLatestBlockNumber
	resp.CurrentBlockNumber = block_parser.CurrentBlockNumber

	apiResp.ApiRespOK(resp)
	return nil
}
