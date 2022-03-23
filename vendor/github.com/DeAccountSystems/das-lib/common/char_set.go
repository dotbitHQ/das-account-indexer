package common

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/molecule"
	"strings"
)

type AccountCharType uint32

const (
	AccountCharTypeEmoji  AccountCharType = 0
	AccountCharTypeNumber AccountCharType = 1
	AccountCharTypeEn     AccountCharType = 2
)

var CharSetTypeEmoji = ""

const (
	CharSetTypeNumber = "0123456789"
	CharSetTypeEn     = "abcdefghijklmnopqrstuvwxyz"
)

type AccountCharSet struct {
	CharSetName AccountCharType `json:"char_set_name"`
	Char        string          `json:"char"`
}

func AccountCharsToAccount(accountChars *molecule.AccountChars) string {
	index := uint(0)
	var accountRawBytes []byte
	accountCharsSize := accountChars.ItemCount()
	for ; index < accountCharsSize; index++ {
		char := accountChars.Get(index)
		accountRawBytes = append(accountRawBytes, char.Bytes().RawData()...)
	}
	accountStr := string(accountRawBytes)
	if accountStr != "" && !strings.HasSuffix(accountStr, DasAccountSuffix) {
		accountStr = accountStr + DasAccountSuffix
	}
	return accountStr
}

func AccountToAccountChars(account string) ([]AccountCharSet, error) {
	if strings.HasSuffix(account, DasAccountSuffix) {
		account = strings.TrimSuffix(account, DasAccountSuffix)
	}

	var list []AccountCharSet
	for _, v := range account {
		char := string(v)
		var charSetName AccountCharType
		if strings.Contains(CharSetTypeEmoji, char) {
			charSetName = AccountCharTypeEmoji
		} else if strings.Contains(CharSetTypeNumber, char) {
			charSetName = AccountCharTypeNumber
		} else if strings.Contains(CharSetTypeEn, char) {
			charSetName = AccountCharTypeEn
		} else {
			return nil, fmt.Errorf("invilid char type")
		}
		list = append(list, AccountCharSet{
			CharSetName: charSetName,
			Char:        char,
		})
	}
	return list, nil
}

func InitEmoji(emojis []string) {
	for _, v := range emojis {
		CharSetTypeEmoji = fmt.Sprintf("%v%v", CharSetTypeEmoji, v)
	}
	//fmt.Println(CharSetTypeEmoji)
}
