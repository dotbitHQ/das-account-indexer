package handle

import (
	"context"
	"das-account-indexer/block_parser"
	"encoding/json"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RespServerInfo struct {
	IsLatestBlockNumber bool   `json:"is_latest_block_number"`
	CurrentBlockNumber  uint64 `json:"current_block_number"`
	Chain               string `json:"chain"`
}

func (h *HttpHandle) JsonRpcServerInfo(p json.RawMessage, apiResp *http_api.ApiResp) {
	if err := h.doServerInfo(h.Ctx, apiResp); err != nil {
		log.Error("doServerInfo err:", err.Error())
	}
}

func (h *HttpHandle) ServerInfo(ctx *gin.Context) {
	var (
		funcName = "ServerInfo"
		apiResp  http_api.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	log.Info("ApiReq:", ctx.Request.Host, funcName, clientIp, ctx.Request.Context())

	if err = h.doServerInfo(ctx.Request.Context(), &apiResp); err != nil {
		log.Error("doServerInfo err:", err.Error(), funcName, ctx.Request.Context())
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doServerInfo(ctx context.Context, apiResp *http_api.ApiResp) error {
	var resp RespServerInfo

	resp.IsLatestBlockNumber = block_parser.IsLatestBlockNumber
	resp.CurrentBlockNumber = block_parser.CurrentBlockNumber

	if h.DasCore.NetType() == common.DasNetTypeMainNet {
		resp.Chain = "mainnet"
	} else {
		resp.Chain = "testnet"
	}

	apiResp.ApiRespOK(resp)
	return nil
}
