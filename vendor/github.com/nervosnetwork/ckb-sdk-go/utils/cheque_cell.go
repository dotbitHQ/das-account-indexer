package utils

import "github.com/nervosnetwork/ckb-sdk-go/types"

func ChequeCellArgs(senderLock, receiverLock *types.Script) ([]byte, error) {
	senderLockHash, err := senderLock.Hash()
	if err != nil {
		return []byte{}, err
	}
	receiverLockHash, err := receiverLock.Hash()
	if err != nil {
		return []byte{}, err
	}
	return append(receiverLockHash.Bytes()[0:20], senderLockHash.Bytes()[0:20]...), nil
}

func IsChequeCell(o *types.CellOutput, systemScripts *SystemScripts) bool {
	return o.Lock.CodeHash == systemScripts.ChequeCell.CellHash
}
