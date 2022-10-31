package handle

import (
	"das-account-indexer/http_server/code"
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqSubAccountList struct {
	Pagination
	Account string `json:"account"`
}

type RespSubAccountList struct {
	Account          string                  `json:"account"`
	AccountIdHex     string                  `json:"account_id_hex"`
	EnableSubAccount tables.EnableSubAccount `json:"enable_sub_account"`
	SubAccountTotal  int64                   `json:"sub_account_total"`
	SubAccountList   []SubAccountInfo        `json:"sub_account_list"`
}

type SubAccountInfo struct {
	Account            string                `json:"account"`
	AccountIdHex       string                `json:"account_id_hex"`
	CreateAtUnix       uint64                `json:"create_at_unix"`
	ExpiredAtUnix      uint64                `json:"expired_at_unix"`
	OwnerAlgorithmId   common.DasAlgorithmId `json:"owner_algorithm_id"`
	OwnerKey           string                `json:"owner_key"`
	ManagerAlgorithmId common.DasAlgorithmId `json:"manager_algorithm_id"`
	ManagerKey         string                `json:"manager_key"`
}

func (h *HttpHandle) JsonRpcSubAccountList(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqSubAccountList
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

	if err = h.doSubAccountList(&req[0], apiResp); err != nil {
		log.Error("doSubAccountList err:", err.Error())
	}
}

func (h *HttpHandle) SubAccountList(ctx *gin.Context) {
	var (
		funcName = "SubAccountList"
		req      ReqSubAccountList
		apiResp  code.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName)
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req))

	if err = h.doSubAccountList(&req, &apiResp); err != nil {
		log.Error("doSubAccountList err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doSubAccountList(req *ReqSubAccountList, apiResp *code.ApiResp) error {
	var resp RespSubAccountList
	resp.SubAccountList = make([]SubAccountInfo, 0)

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	accountInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, "find account info err")
		return fmt.Errorf("FindAccountInfoByAccountId err: %s", err.Error())
	} else if accountInfo.Id == 0 {
		apiResp.ApiRespErr(code.ApiCodeAccountNotExist, "account not exist")
		return nil
	} else if accountInfo.ParentAccountId != "" {
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "parameter [account] is invalid")
		return nil
	}
	resp.Account = accountInfo.Account
	resp.AccountIdHex = accountInfo.AccountId
	resp.EnableSubAccount = accountInfo.EnableSubAccount

	if accountInfo.EnableSubAccount == tables.AccountEnableStatusOn {
		list, err := h.DbDao.GetSubAccountListByParentAccountId(accountId, req.GetLimit(), req.GetOffset())
		if err != nil {
			apiResp.ApiRespErr(code.ApiCodeDbError, "find sub-account list err")
			return fmt.Errorf("GetSubAccountListByParentAccountId err: %s", err.Error())
		}
		for _, v := range list {
			ownerHex := core.DasAddressHex{
				DasAlgorithmId: v.OwnerAlgorithmId,
				AddressHex:     v.Owner,
				IsMulti:        false,
				ChainType:      v.OwnerChainType,
			}
			managerHex := core.DasAddressHex{
				DasAlgorithmId: v.ManagerAlgorithmId,
				AddressHex:     v.Manager,
				IsMulti:        false,
				ChainType:      v.ManagerChainType,
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
			resp.SubAccountList = append(resp.SubAccountList, SubAccountInfo{
				Account:            v.Account,
				AccountIdHex:       v.AccountId,
				CreateAtUnix:       v.RegisteredAt,
				ExpiredAtUnix:      v.ExpiredAt,
				OwnerAlgorithmId:   v.OwnerAlgorithmId,
				OwnerKey:           ownerNormal.AddressNormal,
				ManagerAlgorithmId: v.ManagerAlgorithmId,
				ManagerKey:         managerNormal.AddressNormal,
			})
		}
	}
	count, err := h.DbDao.GetSubAccountListCountByParentAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, "find sub-account count err")
		return fmt.Errorf("GetSubAccountListCountByParentAccountId err: %s", err.Error())
	}
	resp.SubAccountTotal = count

	apiResp.ApiRespOK(resp)
	return nil
}
