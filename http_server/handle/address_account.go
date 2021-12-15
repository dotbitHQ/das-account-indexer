package handle

import (
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
			apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
			return
		} else if len(reqOld) == 1 {
			req[0] = ReqAddressAccount{Address: reqOld[0]}
		}
	}
	if len(req) != 1 {
		log.Error("len(req) is :", len(req))
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
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
		apiResp.ApiRespErr(code.ApiCodeParamsInvalid, "params invalid")
		ctx.JSON(http.StatusOK, apiResp)
		return
	}
	log.Info("ApiReq:", funcName, clientIp, toolib.JsonString(req))

	if err = h.doAddressAccount(&req, &apiResp); err != nil {
		log.Error("doAddressAccount err:", err.Error(), funcName)
	}

	ctx.JSON(http.StatusOK, apiResp)
}

func (h *HttpHandle) doAddressAccount(req *ReqAddressAccount, apiResp *code.ApiResp) error {
	var resp = make([]RespSearchAccount, 0)

	chainType, address := formatAddress(req.Address)
	log.Info("formatAddress:", req.Address, chainType, address)

	list, err := h.DbDao.FindAccountListByAddress(chainType, address)
	if err != nil {
		log.Error("FindAccountListByAddress err:", err.Error(), req.Address)
		apiResp.ApiRespErr(code.ApiCodeDbError, "find account list err")
		return nil
	}
	var accounts []string
	var mapAccountIndex = make(map[string]int)
	for i, v := range list {
		dasLockArgs := core.FormatOwnerManagerAddressToArgs(v.OwnerChainType, v.ManagerChainType, v.Owner, v.Manager)
		tmp := RespSearchAccount{
			OutPoint: common.String2OutPointStruct(v.Outpoint),
			AccountData: AccountData{
				Account:             v.Account,
				AccountIdHex:        v.AccountId,
				NextAccountIdHex:    v.NextAccountId,
				CreateAtUnix:        v.RegisteredAt,
				ExpiredAtUnix:       v.ExpiredAt,
				Status:              v.Status,
				DasLockArgHex:       common.Bytes2Hex(dasLockArgs),
				OwnerAddressChain:   v.OwnerChainType.String(),
				OwnerLockArgsHex:    common.Bytes2Hex(dasLockArgs[:len(dasLockArgs)/2]),
				OwnerAddress:        core.FormatHexAddressToNormal(v.OwnerChainType, v.Owner),
				ManagerAddressChain: v.ManagerChainType.String(),
				ManagerAddress:      core.FormatHexAddressToNormal(v.ManagerChainType, v.Manager),
				ManagerLockArgsHex:  common.Bytes2Hex(dasLockArgs[len(dasLockArgs)/2:]),
				Records:             make([]DataRecord, 0),
			},
		}
		resp = append(resp, tmp)
		accounts = append(accounts, v.Account)
		mapAccountIndex[v.Account] = i
	}

	// records
	if len(accounts) > 0 {
		records, err := h.DbDao.FindRecordsByAccounts(accounts)
		if err != nil {
			log.Error("FindRecordsByAccounts err:", err.Error(), req.Address)
			apiResp.ApiRespErr(code.ApiCodeDbError, "find records info err")
			return nil
		}
		for _, v := range records {
			key := fmt.Sprintf("%s.%s", v.Type, v.Key)
			if index, ok := mapAccountIndex[v.Account]; ok {
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

func formatAddress(address string) (common.ChainType, string) {
	if strings.HasPrefix(address, common.TronBase58PreFix) {
		return common.ChainTypeTron, core.FormatAddressToHex(common.ChainTypeTron, address)
	} else if strings.HasPrefix(address, common.HexPreFix) {
		return common.ChainTypeEth, address
	} else {
		return common.ChainTypeEth, address
	}
}
