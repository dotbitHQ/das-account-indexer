package handle

import (
	"das-account-indexer/http_server/code"
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/core"
	"github.com/gin-gonic/gin"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqSearchAccount struct {
	Account string `json:"account"`
}

type RespSearchAccount struct {
	OutPoint    *types.OutPoint `json:"out_point"`
	AccountData AccountData     `json:"account_data"`
}

type AccountData struct {
	Account             string               `json:"account"`
	AccountAlias        string               `json:"account_alias"`
	AccountIdHex        string               `json:"account_id_hex"`
	NextAccountIdHex    string               `json:"next_account_id_hex"`
	CreateAtUnix        uint64               `json:"create_at_unix"`
	ExpiredAtUnix       uint64               `json:"expired_at_unix"`
	Status              tables.AccountStatus `json:"status"`
	DasLockArgHex       string               `json:"das_lock_arg_hex"`
	OwnerAddressChain   string               `json:"owner_address_chain"`
	OwnerLockArgsHex    string               `json:"owner_lock_args_hex"`
	OwnerAddress        string               `json:"owner_address"`
	ManagerAddressChain string               `json:"manager_address_chain"`
	ManagerAddress      string               `json:"manager_address"`
	ManagerLockArgsHex  string               `json:"manager_lock_args_hex"`
	Records             []DataRecord         `json:"records"`
}

func (h *HttpHandle) JsonRpcSearchAccount(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqSearchAccount
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Warn("json.Unmarshal err:", err.Error())
		var reqOld []string
		if err = json.Unmarshal(p, &reqOld); err != nil {
			log.Error("json.Unmarshal old req err:", err.Error())
			apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
			return
		} else if len(reqOld) == 1 {
			req[0] = ReqSearchAccount{Account: reqOld[0]}
		}
	}
	if len(req) != 1 {
		log.Error("len(req) is:", len(req))
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doSearchAccount(&req[0], apiResp); err != nil {
		log.Error("doSearchAccount err:", err.Error())
	}
}

func (h *HttpHandle) SearchAccount(ctx *gin.Context) {
	var (
		funcName = "SearchAccount"
		req      ReqSearchAccount
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

	if err = h.doSearchAccount(&req, &apiResp); err != nil {
		log.Error("doSearchAccount err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doSearchAccount(req *ReqSearchAccount, apiResp *code.ApiResp) error {
	var resp RespSearchAccount

	req.Account = strings.TrimSpace(req.Account)
	req.Account = FormatSharpToDot(req.Account)
	if err := checkAccount(req.Account, apiResp); err != nil {
		log.Error("checkAccount err: ", err.Error())
		return nil
	}

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		log.Error("FindAccountInfoByAccountName err:", err.Error(), req.Account)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find account info err")
		return nil
	} else if accountInfo.Id == 0 {
		apiResp.ApiRespErr(code.ApiCodeAccountNotExist, "account not exist")
		return nil
	}

	resp.OutPoint = common.String2OutPointStruct(accountInfo.Outpoint)

	ownerHex := core.DasAddressHex{
		DasAlgorithmId: accountInfo.OwnerAlgorithmId,
		AddressHex:     accountInfo.Owner,
		IsMulti:        false,
		ChainType:      accountInfo.OwnerChainType,
	}
	managerHex := core.DasAddressHex{
		DasAlgorithmId: accountInfo.ManagerAlgorithmId,
		AddressHex:     accountInfo.Manager,
		IsMulti:        false,
		ChainType:      accountInfo.ManagerChainType,
	}
	dasLockArgs, err := h.DasCore.Daf().HexToArgs(ownerHex, managerHex)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeError500, err.Error())
		return fmt.Errorf("HexToArgs err: %s", err.Error())
	}
	ownerNormal, err := h.DasCore.Daf().HexToNormal(ownerHex)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeError500, err.Error())
		return fmt.Errorf("owner HexToNormal err: %s", err.Error())
	}
	managerNormal, err := h.DasCore.Daf().HexToNormal(managerHex)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeError500, err.Error())
		return fmt.Errorf("manager HexToNormal err: %s", err.Error())
	}
	resp.AccountData = AccountData{
		Account:             accountInfo.Account,
		AccountAlias:        FormatDotToSharp(accountInfo.Account),
		AccountIdHex:        accountInfo.AccountId,
		NextAccountIdHex:    accountInfo.NextAccountId,
		CreateAtUnix:        accountInfo.RegisteredAt,
		ExpiredAtUnix:       accountInfo.ExpiredAt,
		Status:              accountInfo.Status,
		DasLockArgHex:       common.Bytes2Hex(dasLockArgs),
		OwnerAddressChain:   accountInfo.OwnerChainType.ToString(),
		OwnerLockArgsHex:    common.Bytes2Hex(dasLockArgs[:len(dasLockArgs)/2]),
		OwnerAddress:        ownerNormal.AddressNormal,
		ManagerAddressChain: accountInfo.ManagerChainType.ToString(),
		ManagerAddress:      managerNormal.AddressNormal,
		ManagerLockArgsHex:  common.Bytes2Hex(dasLockArgs[len(dasLockArgs)/2:]),
		Records:             make([]DataRecord, 0),
	}

	// records
	list, err := h.DbDao.FindAccountRecordsByAccountId(accountId)
	if err != nil {
		log.Error("FindAccountRecords err:", err.Error(), req.Account)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find records info err")
		return nil
	}
	for _, v := range list {
		key := fmt.Sprintf("%s.%s", v.Type, v.Key)
		resp.AccountData.Records = append(resp.AccountData.Records, DataRecord{
			Key:   key,
			Label: v.Label,
			Value: v.Value,
			TTL:   v.Ttl,
		})
	}

	apiResp.ApiRespOK(resp)
	return nil
}
