package handle

import (
	"context"
	"das-account-indexer/tables"
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	"github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
)

type ReqAccountList struct {
	core.ChainTypeAddress
	Role string `json:"role"` // owner,manager
	Pagination
}

type RespAccountList struct {
	Total       int64                `json:"total"`
	AccountList []RespAddressAccount `json:"account_list"`
}

type RespAddressAccount struct {
	AccountId    string `json:"account_id"`
	Account      string `json:"account"`
	DisplayName  string `json:"display_name"`
	RegisteredAt uint64 `json:"registered_at"`
	ExpiredAt    uint64 `json:"expired_at"`
}

func (h *HttpHandle) JsonRpcAccountList(p json.RawMessage, apiResp *http_api.ApiResp) {
	var req []ReqAccountList
	err := json.Unmarshal(p, &req)
	if err != nil {
		log.Error("json.Unmarshal err:", err.Error())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}
	if len(req) != 1 {
		log.Error("len(req) is :", len(req))
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		return
	}

	if err = h.doAccountList(h.Ctx, &req[0], apiResp); err != nil {
		log.Error("doAccountList err:", err.Error())
	}
}

func (h *HttpHandle) AccountList(ctx *gin.Context) {
	var (
		funcName = "AccountList"
		req      ReqAccountList
		apiResp  http_api.ApiResp
		err      error
	)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Error("ShouldBindJSON err: ", err.Error(), funcName, ctx.Request.Context())
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, toolib.JsonString(req))

	if err = h.doAccountList(ctx.Request.Context(), &req, &apiResp); err != nil {
		log.Error("doAccountList err:", err.Error(), funcName, ctx.Request.Context())
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAccountList(ctx context.Context, req *ReqAccountList, apiResp *http_api.ApiResp) error {
	var resp RespAccountList
	resp.AccountList = make([]RespAddressAccount, 0)

	//res := checkReqKeyInfo(h.DasCore.Daf(), &req.ChainTypeAddress, apiResp)
	//if apiResp.ErrNo != http_api.ApiCodeSuccess {
	//	log.Error("checkReqReverseRecord:", apiResp.ErrMsg)
	//	return nil
	//}

	addrHex, err := req.FormatChainTypeAddress(h.DasCore.NetType(), true)
	if err != nil {
		apiResp.ApiRespErr(http_api.ApiCodeParamsInvalid, "address invalid")
		return fmt.Errorf("FormatChainTypeAddress err: %s", err.Error())
	}
	log.Info(ctx, "doAccountList:", addrHex.ChainType, addrHex.AddressHex)

	if addrHex.DasAlgorithmId == common.DasAlgorithmIdAnyLock {
		didCells, err := h.DbDao.QueryDidCell(addrHex.AddressHex, tables.DidCellStatusNormal, req.GetLimit(), req.GetOffset())
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "find did cell list err")
			return fmt.Errorf("QueryDidCell err: %s", err.Error())
		}
		var accIds []string
		for _, v := range didCells {
			tmp := RespAddressAccount{
				AccountId:   v.AccountId,
				Account:     v.Account,
				DisplayName: FormatDisplayName(v.Account),
				//RegisteredAt: v.RegisteredAt,
				ExpiredAt: v.ExpiredAt,
			}
			resp.AccountList = append(resp.AccountList, tmp)
			accIds = append(accIds, v.AccountId)
		}
		accs, err := h.DbDao.GetAccountByAccIds(accIds)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "find account list err")
			return fmt.Errorf("GetAccountByAccIds err: %s", err.Error())
		}
		var accRegAtMap = make(map[string]uint64)
		for _, v := range accs {
			accRegAtMap[v.AccountId] = v.RegisteredAt
		}
		for i, v := range resp.AccountList {
			if regAt, ok := accRegAtMap[v.AccountId]; ok {
				resp.AccountList[i].RegisteredAt = regAt
			}
		}

		total, err := h.DbDao.QueryDidCellTotal(addrHex.AddressHex, tables.DidCellStatusNormal)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "search did cell total err")
			return fmt.Errorf("QueryDidCellTotal err: %s", err.Error())
		}
		resp.Total = total
	} else {
		list, err := h.DbDao.FindAccountNameListByAddress(addrHex.ChainType, addrHex.AddressHex, req.Role, req.GetLimit(), req.GetOffset())
		if err != nil {
			log.Error(ctx, "FindAccountListByAddress err:", err.Error(), req.KeyInfo)
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "find account list err")
			return fmt.Errorf("FindAccountListByAddress err: %s", err.Error())
		}

		for _, v := range list {
			tmp := RespAddressAccount{
				Account:      v.Account,
				DisplayName:  FormatDisplayName(v.Account),
				RegisteredAt: v.RegisteredAt,
				ExpiredAt:    v.ExpiredAt,
			}
			resp.AccountList = append(resp.AccountList, tmp)
		}

		total, err := h.DbDao.FindTotalAccountNameListByAddress(addrHex.ChainType, addrHex.AddressHex, req.Role)
		if err != nil {
			apiResp.ApiRespErr(http_api.ApiCodeDbError, "search account total err")
			return fmt.Errorf("FindTotalAccountNameListByAddress err: %s", err.Error())
		}
		resp.Total = total
	}

	apiResp.ApiRespOK(resp)
	return nil
}
