package address

import (
	"encoding/binary"
	"errors"

	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/transaction"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

func GenerateSecp256k1MultisigScript(requireN, threshold int, publicKeys [][]byte) (*types.Script, []byte, error) {
	if requireN < 0 || requireN > 255 {
		return nil, nil, errors.New("requireN must ranging from 0 to 255")
	}
	if threshold < 0 || threshold > 255 {
		return nil, nil, errors.New("requireN must ranging from 0 to 255")
	}
	if len(publicKeys) > 255 {
		return nil, nil, errors.New("public keys size must be less than 256")
	}
	if len(publicKeys) < requireN || len(publicKeys) < threshold {
		return nil, nil, errors.New("public keys error")
	}

	var data []byte
	data = append(data, 0)

	b := make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(requireN))
	data = append(data, b[:1]...)

	b = make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(threshold))
	data = append(data, b[:1]...)

	b = make([]byte, 2)
	binary.LittleEndian.PutUint16(b, uint16(len(publicKeys)))
	data = append(data, b[:1]...)

	for _, pub := range publicKeys {
		hash, err := blake2b.Blake160(pub)
		if err != nil {
			return nil, nil, err
		}
		data = append(data, hash...)
	}

	args, err := blake2b.Blake160(data)
	if err != nil {
		return nil, nil, err
	}

	return &types.Script{
		CodeHash: types.HexToHash(transaction.SECP256K1_BLAKE160_MULTISIG_ALL_TYPE_HASH),
		HashType: types.HashTypeType,
		Args:     args,
	}, data, nil
}
