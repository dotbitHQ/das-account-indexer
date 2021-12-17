package handle

import (
	"das-account-indexer/config"
	"das-account-indexer/http_server/code"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqReverseRecord struct {
	Type    string `json:"type"` // blockchain
	KeyInfo struct {
		CoinType code.CoinType `json:"coin_type"`
		ChainId  code.ChainId  `json:"chain_id"`
		Key      string        `json:"key"`
	} `json:"key_info"`
}

type formatReqReverseRecord struct {
	ChainType common.ChainType
	Address   string
}

type RespReverseRecord struct {
	Account string `json:"account"`
}

func (h *HttpHandle) JsonRpcReverseRecord(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqReverseRecord
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is :", len(req))
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doReverseRecord(&req[0], apiResp); err != nil {
		log.Error("doReverseRecord err:", err.Error())
	}
}

func (h *HttpHandle) ReverseRecord(ctx *gin.Context) {
	var (
		funcName = "ReverseRecord"
		req      ReqReverseRecord
		apiResp  code.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", ctx.Request.Host, funcName, clientIp, toolib.JsonString(req))

	if err = h.doReverseRecord(&req, &apiResp); err != nil {
		log.Error("doReverseRecord err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func checkReqReverseRecord(req *ReqReverseRecord, apiResp *code.ApiResp) *formatReqReverseRecord {
	var res formatReqReverseRecord
	if req.Type != "blockchain" {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, fmt.Sprintf("type [%s] is invalid", req.Type))
		return &res
	}
	dasChainType := code.FormatCoinTypeToDasChainType(req.KeyInfo.CoinType)
	if dasChainType == -1 {
		dasChainType = code.FormatChainIdToDasChainType(config.Cfg.Server.Net, req.KeyInfo.ChainId)
	}
	if dasChainType == -1 {
		if req.KeyInfo.ChainId != "" && strings.HasPrefix(req.KeyInfo.Key, "0x") {
			dasChainType = common.ChainTypeEth
		} else {
			apiResp.ApiRespErr(code.ApiCodeParamsInvalid, fmt.Sprintf("coin_type [%s] and chain_id [%s] is invalid", req.KeyInfo.CoinType, req.KeyInfo.ChainId))
			return &res
		}
	}
	if req.KeyInfo.Key == "" {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "key is invalid")
		return &res
	}
	res.ChainType = dasChainType
	res.Address = core.FormatAddressToHex(dasChainType, req.KeyInfo.Key)
	return &res
}

func (h *HttpHandle) doReverseRecord(req *ReqReverseRecord, apiResp *code.ApiResp) error {
	var resp RespReverseRecord
	res := checkReqReverseRecord(req, apiResp)
	if apiResp.ErrNo != code.ApiCodeSuccess {
		log.Error("checkReqReverseRecord:", apiResp.ErrMsg)
		return nil
	}

	reverse, err := h.DbDao.FindLatestReverseRecord(res.ChainType, res.Address)
	if err != nil {
		log.Error("FindLatestReverseRecord err:", err.Error(), res.ChainType, res.Address)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record err")
		return nil
	} else if reverse.Id == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}

	account, err := h.DbDao.FindAccountInfoByAccountName(reverse.Account)
	if err != nil {
		log.Error("FindAccountInfoByAccountName err:", err.Error(), res.ChainType, res.Address, reverse.Account)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record account err")
		return nil
	}

	if account.OwnerChainType == res.ChainType && strings.EqualFold(account.Owner, res.Address) {
		resp.Account = account.Account
	} else if account.ManagerChainType == res.ChainType && strings.EqualFold(account.Manager, res.Address) {
		resp.Account = account.Account
	} else {
		record, err := h.DbDao.FindRecordByAccountAddressValue(account.Account, res.Address)
		if err != nil {
			log.Error("FindRecordByAccountAddressValue err:", err.Error(), res.ChainType, res.Address, reverse.Account)
			apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record account record err")
			return nil
		} else if record.Id > 0 {
			resp.Account = account.Account
		}
	}

	apiResp.ApiRespOK(resp)
	return nil
}
