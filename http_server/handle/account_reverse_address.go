package handle

import (
	"encoding/json"
	"fmt"
	"github.com/dotbitHQ/das-lib/common"
	"github.com/dotbitHQ/das-lib/core"
	code "github.com/dotbitHQ/das-lib/http_api"
	"github.com/gin-gonic/gin"
	"github.com/scorpiotzh/toolib"
	"net/http"
	"strings"
)

type ReqAccountReverseAddress struct {
	Account string `json:"account"`
}

type RespAccountReverseAddress struct {
	List []core.ChainTypeAddress `json:"list"`
}

type AccountReverseAddress struct {
}

func (h *HttpHandle) JsonRpcAccountReverseAddress(p json.RawMessage, apiResp *code.ApiResp) {
	var req []ReqAccountReverseAddress
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

	if err = h.doAccountReverseAddress(&req[0], apiResp); err != nil {
		log.Error("doAccountReverseAddress err:", err.Error())
	}
}

func (h *HttpHandle) AccountReverseAddress(ctx *gin.Context) {
	var (
		funcName = "AccountReverseAddress"
		req      ReqAccountReverseAddress
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

	if err = h.doAccountReverseAddress(&req, &apiResp); err != nil {
		log.Error("doAccountReverseAddress err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAccountReverseAddress(req *ReqAccountReverseAddress, apiResp *code.ApiResp) error {
	var resp RespAccountReverseAddress
	resp.List = make([]core.ChainTypeAddress, 0)

	accountId := common.Bytes2Hex(common.GetAccountIdByAccount(req.Account))
	accInfo, err := h.DbDao.FindAccountInfoByAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, "Failed to get account info")
		return fmt.Errorf("FindAccountInfoByAccountId err: %s", err.Error())
	} else if accInfo.Id == 0 {
		apiResp.ApiRespErr(code.ApiCodeIndexerAccountNotExist, "Account does not exist")
		return nil
	}

	list, err := h.DbDao.GetReverseListByAccount(req.Account)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, "Failed to get reverse list")
		return fmt.Errorf("GetReverseListByAccount err: %s", err.Error())
	}
	log.Info("doAccountReverseAddress:", len(list), req.Account, accountId)
	if len(list) == 0 {
		apiResp.ApiRespOK(resp)
		return nil
	}
	records, err := h.DbDao.FindAccountRecordsByAccountId(accountId)
	if err != nil {
		apiResp.ApiRespErr(code.ApiCodeDbError, "Failed to get records")
		return fmt.Errorf("FindAccountRecordsByAccountId err: %s", err.Error())
	}
	var recordMap = make(map[string]struct{})
	for _, v := range records {
		if v.Type != "address" {
			continue
		}
		key := strings.ToLower(fmt.Sprintf("%s%s", v.Key, v.Value))
		recordMap[key] = struct{}{}
	}
	for _, v := range list {
		accId := common.Bytes2Hex(common.GetAccountIdByAccount(v.Account))
		if accId != accountId {
			continue
		}
		addrNormal, err := h.DasCore.Daf().HexToNormal(core.DasAddressHex{
			DasAlgorithmId:    v.AlgorithmId,
			DasSubAlgorithmId: 0,
			AddressHex:        v.Address,
			AddressPayload:    nil,
			IsMulti:           false,
			ChainType:         0,
		})
		if err != nil {
			log.Error("HexToNormal err:", err.Error(), v.AlgorithmId, v.Address)
			continue
		}
		if v.AlgorithmId == accInfo.OwnerAlgorithmId && strings.EqualFold(v.Address, accInfo.Owner) {
			resp.List = append(resp.List, core.ChainTypeAddress{
				Type: "blockchain",
				KeyInfo: core.KeyInfo{
					CoinType: v.AlgorithmId.ToCoinType(),
					ChainId:  "",
					Key:      addrNormal.AddressNormal,
				},
			})
			continue
		}
		if v.AlgorithmId == accInfo.ManagerAlgorithmId && strings.EqualFold(v.Address, accInfo.Manager) {
			resp.List = append(resp.List, core.ChainTypeAddress{
				Type: "blockchain",
				KeyInfo: core.KeyInfo{
					CoinType: v.AlgorithmId.ToCoinType(),
					ChainId:  "",
					Key:      addrNormal.AddressNormal,
				},
			})
			continue
		}
		// records
		key := strings.ToLower(fmt.Sprintf("%s%s", v.AlgorithmId.ToCoinType(), addrNormal.AddressNormal))
		if _, ok := recordMap[key]; ok {
			resp.List = append(resp.List, core.ChainTypeAddress{
				Type: "blockchain",
				KeyInfo: core.KeyInfo{
					CoinType: v.AlgorithmId.ToCoinType(),
					ChainId:  "",
					Key:      addrNormal.AddressNormal,
				},
			})
			continue
		}
	}

	apiResp.ApiRespOK(resp)
	return nil
}
