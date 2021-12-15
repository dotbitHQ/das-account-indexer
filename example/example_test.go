package example

import (
	"das-account-indexer/http_server/code"
	"das-account-indexer/http_server/handle"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"testing"
)

func searchAccount(account string) error {
	url := "http://127.0.0.1:8121/v1/search/account"
	var req handle.ReqSearchAccount
	req.Account = account
	var data handle.RespSearchAccount
	var resp code.ApiResp
	resp.Data = &data
	_, _, errs := gorequest.New().Post(url).SendStruct(&req).EndStruct(&resp)
	if errs != nil {
		return fmt.Errorf("%v", errs)
	}
	fmt.Println(data.AccountData.Account)
	return nil
}

// go test -v -run TestSearchAccount
func TestSearchAccount(t *testing.T) {
	if err := searchAccount("duzhihongyi.bit"); err != nil {
		t.Error("searchAccount err:", err.Error())
	}
}

// go test -benchtime=6s -bench=BenchmarkSearchAccount
func BenchmarkSearchAccount(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := searchAccount("duzhihongyi.bit"); err != nil {
			b.Error("searchAccount err:", err.Error())
		}
	}
	b.StopTimer()
}
