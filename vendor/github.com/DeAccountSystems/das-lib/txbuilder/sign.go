package txbuilder

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"strings"
)

func GenerateMultiSignWitnessArgs(firstN uint8, signatures [][]byte, sortArgsList [][]byte) *types.WitnessArgs {
	var lock []byte

	lock = append(lock, 0)

	// first n
	lock = append(lock, firstN)

	// signature num
	sigNum := uint8(len(signatures))
	lock = append(lock, sigNum)

	// address num
	addressNum := uint8(len(sortArgsList))
	lock = append(lock, addressNum)

	// args_of_all_addresses
	for _, v := range sortArgsList {
		lock = append(lock, v...)
	}

	// signatures
	for _, v := range signatures {
		if len(v) > 0 {
			lock = append(lock, v...)
		} else {
			lock = append(lock, transaction.SignaturePlaceholder...)
		}
	}

	wa := types.WitnessArgs{
		Lock:       lock,
		InputType:  nil,
		OutputType: nil,
	}

	return &wa
}

func (d *DasTxBuilder) AddMultiSignatureForTx(group []int, firstN uint8, signatures [][]byte, sortArgsList [][]byte) error {
	wa := GenerateMultiSignWitnessArgs(firstN, signatures, sortArgsList)

	wab, err := wa.Serialize()
	if err != nil {
		return err
	}

	d.Transaction.Witnesses[group[0]] = wab
	return nil
}

func (d *DasTxBuilder) AddSignatureForTx(signData []SignData) error {
	if signData == nil || len(signData) == 0 {
		return fmt.Errorf("signData is nil")
	}
	d.fixDasSignature(signData)
	tmpMapForGroup, err := d.getGroupsFromTx()
	if err != nil {
		return fmt.Errorf("getGroupsFromTx err: %s", err.Error())
	}
	index := 0
	for _, group := range tmpMapForGroup {
		sig := signData[index].SignMsg
		index = index + 1
		if sig == "" {
			continue
		}
		wa := &types.WitnessArgs{
			Lock:       common.Hex2Bytes(sig),
			InputType:  nil,
			OutputType: nil,
		}
		wab, err := wa.Serialize()
		if err != nil {
			return err
		}

		d.Transaction.Witnesses[group[0]] = wab
	}
	return nil
}

func (d *DasTxBuilder) fixDasSignature(signData []SignData) {
	for i, v := range signData {
		log.Info("fixDasSignature:", v.SignMsg)
		switch v.SignType {
		case common.DasAlgorithmIdCkb:
		case common.DasAlgorithmIdEth, common.DasAlgorithmIdEth712:
			if len(v.SignMsg) >= 132 && v.SignMsg[130:132] == "1b" {
				signData[i].SignMsg = v.SignMsg[0:130] + "00" + v.SignMsg[132:len(v.SignMsg)]
			}
			if len(v.SignMsg) >= 132 && v.SignMsg[130:132] == "1c" {
				signData[i].SignMsg = v.SignMsg[0:130] + "01" + v.SignMsg[132:len(v.SignMsg)]
			}
		case common.DasAlgorithmIdTron:
			if strings.HasSuffix(v.SignMsg, "1b") {
				signData[i].SignMsg = v.SignMsg[0:len(v.SignMsg)-2] + "00"
			}
			if strings.HasSuffix(v.SignMsg, "1c") {
				signData[i].SignMsg = v.SignMsg[0:len(v.SignMsg)-2] + "01"
			}
		default:
			log.Warn("unknown sign type:", v.SignType)
		}
	}
}

func (d *DasTxBuilder) serverSignTx() error {
	if len(d.ServerSignGroup) == 0 {
		return nil
	}
	if d.handleServerSign == nil {
		return fmt.Errorf("handleRemoteSign is nil")
	}
	if digest, err := d.generateDigestByGroup(d.ServerSignGroup, []int{}); err != nil {
		return fmt.Errorf("generateDigestByGroup err: %s", err.Error())
	} else {
		sig, err := d.handleServerSign(digest.SignMsg)
		if err != nil {
			return fmt.Errorf("handleServerSign err: %s", err.Error())
		}

		wa := &types.WitnessArgs{
			Lock:       sig,
			InputType:  nil,
			OutputType: nil,
		}
		wab, err := wa.Serialize()
		if err != nil {
			return err
		}

		d.Transaction.Witnesses[d.ServerSignGroup[0]] = wab
		return nil
	}
}
