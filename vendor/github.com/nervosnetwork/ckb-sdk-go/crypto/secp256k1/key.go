package secp256k1

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"

	"github.com/nervosnetwork/ckb-sdk-go/crypto"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
	"github.com/nervosnetwork/ckb-sdk-go/utils"
)

var (
	secp256k1N, _  = new(big.Int).SetString("fffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141", 16)
	secp256k1halfN = new(big.Int).Div(secp256k1N, big.NewInt(2))
)

type Secp256k1Key struct {
	PrivateKey *ecdsa.PrivateKey
}

func (k *Secp256k1Key) Bytes() []byte {
	return math.PaddedBigBytes(k.PrivateKey.D, k.PrivateKey.Params().BitSize/8)
}

func (k *Secp256k1Key) Sign(data []byte) ([]byte, error) {
	seckey := k.Bytes()
	defer crypto.ZeroBytes(seckey)

	return secp256k1.Sign(data, seckey)
}

func (k *Secp256k1Key) Script(systemScripts *utils.SystemScripts) (*types.Script, error) {
	pub := k.PubKey()

	args, err := blake2b.Blake160(pub)
	if err != nil {
		return nil, err
	}

	return &types.Script{
		CodeHash: systemScripts.SecpSingleSigCell.CellHash,
		HashType: types.HashTypeType,
		Args:     args,
	}, nil
}

func (k *Secp256k1Key) PubKey() []byte {
	pub := &k.PrivateKey.PublicKey
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}

	return secp256k1.CompressPubkey(pub.X, pub.Y)
}

func RandomNew() (*Secp256k1Key, error) {
	randBytes := make([]byte, 64)
	_, err := rand.Read(randBytes)
	if err != nil {
		return nil, errors.New("key generation: could not read from random source: " + err.Error())
	}
	reader := bytes.NewReader(randBytes)
	priv, err := ecdsa.GenerateKey(secp256k1.S256(), reader)
	if err != nil {
		return nil, errors.New("key generation: ecdsa.GenerateKey failed: " + err.Error())
	}

	return &Secp256k1Key{PrivateKey: priv}, nil
}

func HexToKey(hexkey string) (*Secp256k1Key, error) {
	b, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, errors.New("invalid hex string")
	}
	return ToKey(b)
}

func ToKey(d []byte) (*Secp256k1Key, error) {
	return toKey(d, true)
}

func toKey(d []byte, strict bool) (*Secp256k1Key, error) {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = secp256k1.S256()
	if strict && 8*len(d) != priv.Params().BitSize {
		return nil, fmt.Errorf("invalid length, need %d bits", priv.Params().BitSize)
	}
	priv.D = new(big.Int).SetBytes(d)

	// The priv.D must < N
	if priv.D.Cmp(secp256k1N) >= 0 {
		return nil, errors.New("invalid private key, >=N")
	}
	// The priv.D must not be zero or negative.
	if priv.D.Sign() <= 0 {
		return nil, errors.New("invalid private key, zero or negative")
	}

	priv.PublicKey.X, priv.PublicKey.Y = priv.PublicKey.Curve.ScalarBaseMult(d)
	if priv.PublicKey.X == nil {
		return nil, errors.New("invalid private key")
	}
	return &Secp256k1Key{PrivateKey: priv}, nil
}
