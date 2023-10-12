package handle

import (
	"das-account-indexer/http_server/code"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqAddressAccount struct {
	Address string `json:"address"`
}

func (h *HttpHandle) JsonRpcAddressAccount(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqAddressAccount
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Warn("json.Unmarshal err:", err.Error())
		var reqOld []string
		if err = json.Unmarshal(p, &reqOld); err != nil {
			log.Error("json.Unmarshal old req err:", err.Error())
			apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
			return
		} else if len(reqOld) == 1 {
			req[0] = ReqAddressAccount{Address: reqOld[0]}
		}
	}
	if len(req) != 1 {
		log.Error("len(req) is :", len(req))
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doAddressAccount(&req[0], apiResp); err != nil {
		log.Error("doAddressAccount err:", err.Error())
	}
}

func (h *HttpHandle) AddressAccount(ctx *gin.Context) {
	var (
		funcName = "AddressAccount"
		req      ReqAddressAccount
		apiResp  code.ApiResp
		err      error
		clientIp = GetClientIp(ctx)
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", ctx.Request.Host, funcName, clientIp, toolib.JsonString(req))

	if err = h.doAddressAccount(&req, &apiResp); err != nil {
		log.Error("doAddressAccount err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAddressAccount(req *ReqAddressAccount, apiResp *code.ApiResp) error {
	var resp = make([]RespSearchAccount, 0)

	addrHex, err := formatAddress(h.DasCore.Daf(), req.Address)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, err.Error())
		return fmt.Errorf("formatAddress err: %s", err.Error())
	}
	log.Info("formatAddress:", req.Address, addrHex.ChainType, addrHex.AddressHex)

	list, err := h.DbDao.FindAccountListByAddress(addrHex.ChainType, addrHex.AddressHex)
	if err != nil {
		log.Error("FindAccountListByAddress err:", err.Error(), req.Address)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find account list err")
		return nil
	}
	var accountIds []string
	var mapAccountIndex = make(map[string]int)
	for i, v := range list {
		ownerHex := core.DasAddressHex{
			DasAlgorithmId:    v.OwnerAlgorithmId,
			DasSubAlgorithmId: v.OwnerSubAid,
			AddressHex:        v.Owner,
			IsMulti:           false,
			ChainType:         v.OwnerChainType,
		}
		managerHex := core.DasAddressHex{
			DasAlgorithmId:    v.ManagerAlgorithmId,
			DasSubAlgorithmId: v.ManagerSubAid,
			AddressHex:        v.Manager,
			IsMulti:           false,
			ChainType:         v.ManagerChainType,
		}
		dasLockArgs, err := h.DasCore.Daf().HexToArgs(ownerHex, managerHex)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
			return fmt.Errorf("HexToArgs err: %s", err.Error())
		}
		ownerNormal, err := h.DasCore.Daf().HexToNormal(ownerHex)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
			return fmt.Errorf("owner HexToNormal err: %s", err.Error())
		}
		managerNormal, err := h.DasCore.Daf().HexToNormal(managerHex)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeError500, err.Error())
			return fmt.Errorf("manager HexToNormal err: %s", err.Error())
		}
		tmp := RespSearchAccount{
			OutPoint: common.String2OutPointStruct(v.Outpoint),
			AccountData: AccountData{
				Account:             v.Account,
				AccountAlias:        FormatDotToSharp(v.Account),
				AccountIdHex:        v.AccountId,
				NextAccountIdHex:    v.NextAccountId,
				CreateAtUnix:        v.RegisteredAt,
				ExpiredAtUnix:       v.ExpiredAt,
				Status:              v.Status,
				DasLockArgHex:       common.Bytes2Hex(dasLockArgs),
				OwnerAddressChain:   v.OwnerChainType.ToString(),
				OwnerLockArgsHex:    common.Bytes2Hex(dasLockArgs[:len(dasLockArgs)/2]),
				OwnerAddress:        ownerNormal.AddressNormal,
				ManagerAddressChain: v.ManagerChainType.ToString(),
				ManagerAddress:      managerNormal.AddressNormal,
				ManagerLockArgsHex:  common.Bytes2Hex(dasLockArgs[len(dasLockArgs)/2:]),
				Records:             make([]DataRecord, 0),
				DisplayName:         FormatDisplayName(v.Account),
			},
		}
		resp = append(resp, tmp)
		accountId := common.Bytes2Hex(common.GetAccountIdByAccount(v.Account))
		accountIds = append(accountIds, accountId)
		mapAccountIndex[accountId] = i
	}

	// records
	if len(accountIds) > 0 {
		records, err := h.DbDao.FindRecordsByAccountIds(accountIds)
		if err != nil {
			log.Error("FindRecordsByAccounts err:", err.Error(), req.Address)
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "find records info err")
			return nil
		}
		for _, v := range records {
			key := fmt.Sprintf("%s.%s", v.Type, v.Key)
			if index, ok := mapAccountIndex[v.AccountId]; ok {
				resp[index].AccountData.Records = append(resp[index].AccountData.Records, DataRecord{
					Key:   key,
					Label: v.Label,
					Value: v.Value,
					TTL:   v.Ttl,
				})
			}
		}
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func formatAddress(daf *core.DasAddressFormat, addr string) (core.DasAddressHex, error) {
	chainType := common.ChainTypeEth
	if strings.HasPrefix(addr, common.TronBase58PreFix) {
		chainType = common.ChainTypeTron
		//return common.ChainTypeTron, core.FormatAddressToHex(common.ChainTypeTron, address)
	} else if strings.HasPrefix(addr, common.TronPreFix) {
		chainType = common.ChainTypeTron
	} else if strings.HasPrefix(addr, "ckt") || strings.HasPrefix(addr, "ckb") {
		chainType = common.ChainTypeCkbMulti
	} else if strings.HasPrefix(addr, common.HexPreFix) && len(addr) == 42 {
		chainType = common.ChainTypeEth
	} else if strings.HasPrefix(addr, common.HexPreFix) && len(addr) == 66 {
		chainType = common.ChainTypeMixin
	}
	return daf.NormalToHex(core.DasAddressNormal{
		ChainType:     chainType,
		AddressNormal: addr,
		Is712:         true,
	})
}
