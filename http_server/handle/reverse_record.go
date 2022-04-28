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
	"regexp"
	"strings"
)

type ReqReverseRecord struct {
	ReqKeyInfo
}

type ReqKeyInfo struct {
	Type    string `json:"type"` // blockchain
	KeyInfo struct {
		CoinType code.CoinType `json:"coin_type"`
		ChainId  code.ChainId  `json:"chain_id"`
		Key      string        `json:"key"`
	} `json:"key_info"`
}

type RespReverseRecord struct {
	Account      string `json:"account"`
	AccountAlias string `json:"account_alias"`
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

func checkReqKeyInfo(daf *core.DasAddressFormat, req *ReqKeyInfo, apiResp *code.ApiResp) *core.DasAddressHex {
	if req.Type != "blockchain" {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, fmt.Sprintf("type [%s] is invalid", req.Type))
		return nil
	}
	if req.KeyInfo.Key == "" {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "key is invalid")
		return nil
	}
	dasChainType := code.FormatCoinTypeToDasChainType(req.KeyInfo.CoinType)
	if dasChainType == -1 {
		dasChainType = code.FormatChainIdToDasChainType(config.Cfg.Server.Net, req.KeyInfo.ChainId)
	}
	if dasChainType == -1 {
		if strings.HasPrefix(req.KeyInfo.Key, "0x") {
			if ok, err := regexp.MatchString("^0x[0-9a-fA-F]{40}$", req.KeyInfo.Key); err != nil {
				apiResp.ApiRespErr(code.ApiCodeParamsInvalid, err.Error())
				return nil
			} else if ok {
				dasChainType = common.ChainTypeEth
			} else if ok, err = regexp.MatchString("^0x[0-9a-fA-F]{64}$", req.KeyInfo.Key); err != nil {
				apiResp.ApiRespErr(code.ApiCodeParamsInvalid, err.Error())
				return nil
			} else if ok {
				dasChainType = common.ChainTypeMixin
			} else {
				apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "key is invalid")
				return nil
			}
		} else {
			apiResp.ApiRespErr(code.ApiCodeParamsInvalid, fmt.Sprintf("coin_type [%s] and chain_id [%s] is invalid", req.KeyInfo.CoinType, req.KeyInfo.ChainId))
			return nil
		}
	}
	addrHex, err := daf.NormalToHex(core.DasAddressNormal{
		ChainType:     dasChainType,
		AddressNormal: req.KeyInfo.Key,
		Is712:         true,
	})
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, err.Error())
		return nil
	}
	return &addrHex
}

func (h *HttpHandle) doReverseRecord(req *ReqReverseRecord, apiResp *code.ApiResp) error {
	var resp RespReverseRecord
	res := checkReqKeyInfo(h.DasCore.Daf(), &req.ReqKeyInfo, apiResp)
	if apiResp.ErrNo != code.ApiCodeSuccess {
		log.Error("checkReqReverseRecord:", apiResp.ErrMsg)
		return nil
	}

	reverse, err := h.DbDao.FindLatestReverseRecord(res.ChainType, res.AddressHex)
	if err != nil {
		log.Error("FindLatestReverseRecord err:", err.Error(), res.ChainType, res.AddressHex)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record err")
		return nil
	} else if reverse.Id == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(reverse.Account))
	account, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		log.Error("FindAccountInfoByAccountName err:", err.Error(), res.ChainType, res.AddressHex, reverse.Account)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record account err")
		return nil
	}

	if account.OwnerChainType == res.ChainType && strings.EqualFold(account.Owner, res.AddressHex) {
		resp.Account = account.Account
	} else if account.ManagerChainType == res.ChainType && strings.EqualFold(account.Manager, res.AddressHex) {
		resp.Account = account.Account
	} else {
		record, err := h.DbDao.FindRecordByAccountIdAddressValue(account.AccountId, res.AddressHex)
		if err != nil {
			log.Error("FindRecordByAccountAddressValue err:", err.Error(), res.ChainType, res.AddressHex, reverse.Account)
			apiResp.ApiRespErr(code.ApiCodeDbError, "find reverse record account record err")
			return nil
		} else if record.Id > 0 {
			resp.Account = account.Account
		}
	}

	resp.AccountAlias = FormatDotToSharp(resp.Account)

	apiResp.ApiRespOK(resp)
	return nil
}

func FormatDotToSharp(account string) string {
	countDot := strings.Count(account, ".")
	countSharp := strings.Count(account, "#")
	if countDot == 2 && countSharp == 0 {
		list := strings.Split(account, ".")
		return list[1] + "#" + list[0] + ".bit"
	}
	return account
}

func FormatSharpToDot(account string) string {
	countDot := strings.Count(account, ".")
	countSharp := strings.Count(account, "#")
	if countDot == 1 && countSharp == 1 {
		indexSharp := strings.Index(account, "#")
		indexDot := strings.Index(account, ".")

		return account[indexSharp+1:indexDot] + "." + account[:indexSharp] + ".bit"
	}
	return account
}
