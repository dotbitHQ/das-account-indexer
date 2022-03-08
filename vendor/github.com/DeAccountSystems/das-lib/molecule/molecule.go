package molecule

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

func GoU8ToMoleculeU8(i uint8) Uint8 {
	byteBuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(byteBuf, binary.LittleEndian, i)
	return *Uint8FromSliceUnchecked(byteBuf.Bytes())
}

func GoU32ToMoleculeU32(i uint32) Uint32 {
	byteBuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(byteBuf, binary.LittleEndian, i)
	return *Uint32FromSliceUnchecked(byteBuf.Bytes())
}

func GoU64ToMoleculeU64(i uint64) Uint64 {
	byteBuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(byteBuf, binary.LittleEndian, i)
	return *Uint64FromSliceUnchecked(byteBuf.Bytes())
}

func Bytes2GoU8(bys []byte) (uint8, error) {
	var t uint8
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func Bytes2GoU16(bys []byte) (uint16, error) {
	var t uint16
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func Bytes2GoU32(bys []byte) (uint32, error) {
	var t uint32
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func Bytes2GoU64(bys []byte) (uint64, error) {
	var t uint64
	bytesBuffer := bytes.NewBuffer(bys)
	if err := binary.Read(bytesBuffer, binary.LittleEndian, &t); err != nil {
		return 0, err
	}
	return t, nil
}

func GoBytes2MoleculeBytes(bys []byte) Bytes {
	_bytesBuilder := NewBytesBuilder()
	for _, bye := range bys {
		_bytesBuilder.Push(*ByteFromSliceUnchecked([]byte{bye}))
	}
	return _bytesBuilder.Build()
}

func GoString2MoleculeBytes(str string) Bytes {
	if str == "" {
		return BytesDefault()
	}
	strBytes := []byte(str)
	return GoBytes2MoleculeBytes(strBytes)
}

func CkbScript2MoleculeScript(script *types.Script) Script {
	// data 0x00 ，type 0x01
	ht := 0
	if script.HashType == types.HashTypeType {
		ht = 1
	}
	argBytes := BytesDefault()
	if script.Args != nil {
		argBytes = GoBytes2MoleculeBytes(script.Args)
	}
	return NewScriptBuilder().
		CodeHash(GoHexToMoleculeHash(script.CodeHash.String())).
		HashType(NewByte(byte(ht))).
		Args(argBytes).
		Build()
}

func MoleculeScript2CkbScript(script *Script) *types.Script {
	if script == nil {
		return nil
	}
	tmp := &types.Script{
		CodeHash: types.BytesToHash(script.CodeHash().RawData()),
		Args:     script.Args().RawData(),
		HashType: types.HashTypeType,
	}
	return tmp
}

func Has0xPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func GoHexToMoleculeHash(hexStr string) Hash {
	if Has0xPrefix(hexStr) {
		hexStr = hexStr[2:]
	}
	bys, _ := hex.DecodeString(hexStr)
	byteArr := [32]Byte{}
	size := len(bys)
	for i := 0; i < size; i++ {
		byteArr[i] = *ByteFromSliceUnchecked([]byte{bys[i]})
	}
	return NewHashBuilder().Set(byteArr).Build()
}

func GoTimeUnixToMoleculeBytes(timeSec int64) [8]Byte {
	byteBuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(byteBuf, binary.LittleEndian, timeSec)
	timestampByteArr := [8]Byte{}
	tmpBytes := byteBuf.Bytes()
	size := len(tmpBytes)
	for i := 0; i < size; i++ {
		timestampByteArr[i] = *ByteFromSliceUnchecked([]byte{tmpBytes[i]})
	}
	return timestampByteArr
}

func Go64ToBytes(i int64) []byte {
	byteBuf := bytes.NewBuffer([]byte{})
	_ = binary.Write(byteBuf, binary.LittleEndian, i)
	return byteBuf.Bytes()
}

func GoU32ToBytes(data uint32) []byte {
	buffer := bytes.NewBuffer([]byte{})
	_ = binary.Write(buffer, binary.LittleEndian, data)
	return buffer.Bytes()
}
