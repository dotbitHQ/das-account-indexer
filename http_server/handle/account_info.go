package handle

import (
	"das-account-indexer/config"
	"das-account-indexer/http_server/code"
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/address"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqAccountInfo struct {
	Account   string `json:"account"`
	AccountId string `json:"account_id"`
}

type RespAccountInfo struct {
	OutPoint    *types.OutPoint `json:"out_point"`
	AccountInfo AccountInfo     `json:"account_info"`
}

type AccountInfo struct {
	Account            string                   `json:"account"`
	AccountAlias       string                   `json:"account_alias"`
	AccountIdHex       string                   `json:"account_id_hex"`
	NextAccountIdHex   string                   `json:"next_account_id_hex"`
	CreateAtUnix       uint64                   `json:"create_at_unix"`
	ExpiredAtUnix      uint64                   `json:"expired_at_unix"`
	Status             tables.AccountStatus     `json:"status"`
	DasLockArgHex      string                   `json:"das_lock_arg_hex"`
	OwnerAlgorithmId   common.DasAlgorithmId    `json:"owner_algorithm_id"`
	OwnerSubAid        common.DasSubAlgorithmId `json:"owner_sub_aid"`
	OwnerKey           string                   `json:"owner_key"`
	ManagerAlgorithmId common.DasAlgorithmId    `json:"manager_algorithm_id"`
	ManagerSubAid      common.DasSubAlgorithmId `json:"manager_sub_aid"`
	ManagerKey         string                   `json:"manager_key"`
	EnableSubAccount   tables.EnableSubAccount  `json:"enable_sub_account"`
	DisplayName        string                   `json:"display_name"`
}

func (h *HttpHandle) JsonRpcAccountInfo(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqAccountInfo
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is:", len(req))
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doAccountInfo(&req[0], apiResp); err != nil {
		log.Error("doAccountInfo err:", err.Error())
	}
}

func (h *HttpHandle) AccountInfo(ctx *gin.Context) {
	var (
		funcName = "AccountInfo"
		req      ReqAccountInfo
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

	if err = h.doAccountInfo(&req, &apiResp); err != nil {
		log.Error("doAccountInfo err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAccountInfo(req *ReqAccountInfo, apiResp *code.ApiResp) error {
	var resp RespAccountInfo

	accountId := req.AccountId
	if accountId == "" {
		req.Account = strings.TrimSpace(req.Account)
		req.Account = FormatSharpToDot(req.Account)
		if err := checkAccount(req.Account, apiResp); err != nil {
			log.Error("checkAccount err: ", err.Error())
			return nil
		}
		accountId = common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	}
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		log.Error("FindAccountInfoByAccountName err:", err.Error(), req.Account)
		apiResp.ApiRespErr(http_api.ApiCodeDbError, "find account info err")
		return nil
	} else if accountInfo.Id == 0 {
		apiResp.ApiRespErr(http_api.ApiCodeAccountNotExist, "account not exist")
		return nil
	}

	resp.OutPoint = common.String2OutPointStruct(accountInfo.Outpoint)
	resp.AccountInfo = AccountInfo{
		Account:          accountInfo.Account,
		AccountAlias:     FormatDotToSharp(accountInfo.Account),
		AccountIdHex:     accountInfo.AccountId,
		NextAccountIdHex: accountInfo.NextAccountId,
		CreateAtUnix:     accountInfo.RegisteredAt,
		ExpiredAtUnix:    accountInfo.ExpiredAt,
		Status:           accountInfo.Status,
		//DasLockArgHex:      common.Bytes2Hex(dasLockArgs),
		//OwnerAlgorithmId: accountInfo.OwnerAlgorithmId,
		//OwnerSubAid:      accountInfo.OwnerSubAid,
		//OwnerKey:           ownerNormal.AddressNormal,
		ManagerAlgorithmId: accountInfo.ManagerAlgorithmId,
		ManagerSubAid:      accountInfo.ManagerSubAid,
		//ManagerKey:         managerNormal.AddressNormal,
		EnableSubAccount: accountInfo.EnableSubAccount,
		DisplayName:      FormatDisplayName(accountInfo.Account),
	}

	if accountInfo.Status == tables.AccountStatusOnLock {
		didCell, err := h.DbDao.GetDidCellByAccountId(accountId)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "find did cell info err")
			return fmt.Errorf("GetDidCellByAccountId err: %s", err.Error())
		} else if didCell.Id == 0 {
			apiResp.ApiRespErr(http_api.ApiCodeAccountNotExist, "did cell not exist")
			return nil
		}

		mode := address.Mainnet
		if config.Cfg.Server.Net != common.DasNetTypeMainNet {
			mode = address.Testnet
		}
		addrOwner, err := didCell.ToAnyLockAddr(mode)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeError500, "Failed to get did cell addr")
			return fmt.Errorf("ConvertScriptToAddress err: %s", err.Error())
		}
		resp.AccountInfo.OwnerKey = addrOwner
		resp.AccountInfo.OwnerAlgorithmId = common.DasAlgorithmIdAnyLock
		resp.AccountInfo.ManagerKey = addrOwner
		resp.AccountInfo.ManagerAlgorithmId = common.DasAlgorithmIdAnyLock
	} else {
		ownerHex := core.DasAddressHex{
			DasAlgorithmId:    accountInfo.OwnerAlgorithmId,
			DasSubAlgorithmId: accountInfo.OwnerSubAid,
			AddressHex:        accountInfo.Owner,
			IsMulti:           false,
			ChainType:         accountInfo.OwnerChainType,
		}
		managerHex := core.DasAddressHex{
			DasAlgorithmId:    accountInfo.ManagerAlgorithmId,
			DasSubAlgorithmId: accountInfo.ManagerSubAid,
			AddressHex:        accountInfo.Manager,
			IsMulti:           false,
			ChainType:         accountInfo.ManagerChainType,
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
		resp.AccountInfo.DasLockArgHex = common.Bytes2Hex(dasLockArgs)
		resp.AccountInfo.OwnerAlgorithmId = accountInfo.OwnerAlgorithmId
		resp.AccountInfo.OwnerSubAid = accountInfo.OwnerSubAid
		resp.AccountInfo.OwnerKey = ownerNormal.AddressNormal
		resp.AccountInfo.ManagerKey = managerNormal.AddressNormal
		resp.AccountInfo.ManagerAlgorithmId = accountInfo.ManagerAlgorithmId
		resp.AccountInfo.ManagerSubAid = accountInfo.ManagerSubAid
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func checkAccount(account string, apiResp *code.ApiResp) error {
	if account == "" || !strings.HasSuffix(account, common.DasAccountSuffix) ||
		strings.Contains(account, " ") || strings.Contains(account, "_") {
		apiResp.ApiRespErr(http_api.ApiCodeAccountFormatInvalid, "account invalid")
		return fmt.Errorf("account invalid: [%s]", account)
	}
	return nil
}
