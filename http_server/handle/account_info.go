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

type ReqAccountInfo struct {
	Account string `json:"account"`
}

type RespAccountInfo struct {
	OutPoint    *types.OutPoint `json:"out_point"`
	AccountInfo AccountInfo     `json:"account_info"`
}

type AccountInfo struct {
	Account            string                `json:"account"`
	AccountIdHex       string                `json:"account_id_hex"`
	NextAccountIdHex   string                `json:"next_account_id_hex"`
	CreateAtUnix       uint64                `json:"create_at_unix"`
	ExpiredAtUnix      uint64                `json:"expired_at_unix"`
	Status             tables.AccountStatus  `json:"status"`
	DasLockArgHex      string                `json:"das_lock_arg_hex"`
	OwnerAlgorithmId   common.DasAlgorithmId `json:"owner_algorithm_id"`
	OwnerKey           string                `json:"owner_key"`
	ManagerAlgorithmId common.DasAlgorithmId `json:"manager_algorithm_id"`
	ManagerKey         string                `json:"manager_key"`
}

func (h *HttpHandle) JsonRpcAccountInfo(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqAccountInfo
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is:", len(req))
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
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
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
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

	req.Account = strings.TrimSpace(req.Account)
	if err := checkAccount(req.Account, apiResp); err != nil {
		log.Error("checkAccount err: ", err.Error())
		return nil
	}

	accountInfo, err := h.DbDao.FindAccountInfoByAccountName(req.Account)
	if err != nil {
		log.Error("FindAccountInfoByAccountName err:", err.Error(), req.Account)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find account info err")
		return nil
	} else if accountInfo.Id == 0 {
		apiResp.ApiRespErr(code.ApiCodeAccountNotExist, "account not exist")
		return nil
	}

	resp.OutPoint = common.String2OutPointStruct(accountInfo.Outpoint)
	dasLockArgs := core.FormatOwnerManagerAddressToArgs(accountInfo.OwnerChainType, accountInfo.ManagerChainType, accountInfo.Owner, accountInfo.Manager)
	resp.AccountInfo = AccountInfo{
		Account:            accountInfo.Account,
		AccountIdHex:       accountInfo.AccountId,
		NextAccountIdHex:   accountInfo.NextAccountId,
		CreateAtUnix:       accountInfo.RegisteredAt,
		ExpiredAtUnix:      accountInfo.ExpiredAt,
		Status:             accountInfo.Status,
		DasLockArgHex:      common.Bytes2Hex(dasLockArgs),
		OwnerAlgorithmId:   accountInfo.OwnerAlgorithmId,
		OwnerKey:           core.FormatHexAddressToNormal(accountInfo.OwnerChainType, accountInfo.Owner),
		ManagerAlgorithmId: accountInfo.ManagerAlgorithmId,
		ManagerKey:         core.FormatHexAddressToNormal(accountInfo.ManagerChainType, accountInfo.Manager),
	}

	apiResp.ApiRespOK(resp)
	return nil
}

func checkAccount(account string, apiResp *code.ApiResp) error {
	if account == "" || !strings.HasSuffix(account, common.DasAccountSuffix) ||
		strings.Contains(account, " ") || strings.Contains(account, "_") {
		apiResp.ApiRespErr(code.ApiCodeAccountFormatInvalid, "account invalid")
		return fmt.Errorf("account invalid: [%s]", account)
	}
	return nil
}
