package types

import (
	"bytes"
	"encoding/binary"
	"reflect"
)

const u32Size uint = 4

type Serializer interface {
	Serialize() ([]byte, error)
}

func SerializeUint(n uint) []byte {
	b := make([]byte, u32Size)
	binary.LittleEndian.PutUint32(b, uint32(n))

	return b
}

func SerializeUint64(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)

	return b
}

// SerializeArray serialize array
func SerializeArray(items []Serializer) ([][]byte, error) {
	ret := make([][]byte, len(items))
	for i := 0; i < len(items); i++ {
		data, err := items[i].Serialize()
		if err != nil {
			return nil, err
		}

		ret[i] = data
	}

	return ret, nil
}

// SerializeStruct serialize struct
func SerializeStruct(fields [][]byte) []byte {
	b := new(bytes.Buffer)

	for i := 0; i < len(fields); i++ {
		b.Write(fields[i])
	}

	return b.Bytes()
}

// SerializeFixVec serialize bytes
// There are two steps of serializing a bytes:
//   Serialize the length as a 32 bit unsigned integer in little-endian.
//   Serialize all items in it.
func SerializeBytes(items []byte) []byte {
	// Empty fix vector bytes
	if len(items) == 0 {
		return []byte{00, 00, 00, 00}
	}

	l := SerializeUint(uint(len(items)))

	b := new(bytes.Buffer)

	b.Write(l)
	b.Write(items)

	return b.Bytes()
}

// SerializeFixVec serialize fixvec vector
// There are two steps of serializing a fixvec:
//   Serialize the length as a 32 bit unsigned integer in little-endian.
//   Serialize all items in it.
func SerializeFixVec(items [][]byte) []byte {
	// Empty fix vector bytes
	if len(items) == 0 {
		return []byte{00, 00, 00, 00}
	}

	l := SerializeUint(uint(len(items)))

	b := new(bytes.Buffer)

	b.Write(l)

	for i := 0; i < len(items); i++ {
		b.Write(items[i])
	}

	return b.Bytes()
}

// SerializeDynVec serialize dynvec
// There are three steps of serializing a dynvec:
//    Serialize the full size in bytes as a 32 bit unsigned integer in little-endian.
//    Serialize all offset of items as 32 bit unsigned integer in little-endian.
//    Serialize all items in it.
func SerializeDynVec(items [][]byte) []byte {
	// Start with u32Size
	size := u32Size

	// Empty dyn vector, just return size's bytes
	if len(items) == 0 {
		return SerializeUint(size)
	}

	offsets := make([]uint, len(items))

	// Calculate first offset then loop for rest items offsets
	offsets[0] = size + u32Size*uint(len(items))
	for i := 0; i < len(items); i++ {
		size += u32Size + uint(len(items[i]))

		if i != 0 {
			offsets[i] = offsets[i-1] + uint(len(items[i-1]))
		}
	}

	b := new(bytes.Buffer)

	b.Write(SerializeUint(size))

	for i := 0; i < len(items); i++ {
		b.Write(SerializeUint(offsets[i]))
	}

	for i := 0; i < len(items); i++ {
		b.Write(items[i])
	}

	return b.Bytes()
}

// SerializeTable serialize table
// The serializing steps are same as table:
//    Serialize the full size in bytes as a 32 bit unsigned integer in little-endian.
//    Serialize all offset of fields as 32 bit unsigned integer in little-endian.
//    Serialize all fields in it in the order they are declared.
func SerializeTable(fields [][]byte) []byte {
	size := u32Size
	offsets := make([]uint, len(fields))

	// Calculate first offset then loop for rest items offsets
	offsets[0] = u32Size + u32Size*uint(len(fields))
	for i := 0; i < len(fields); i++ {
		size += u32Size + uint(len(fields[i]))

		if i != 0 {
			offsets[i] = offsets[i-1] + uint(len(fields[i-1]))
		}
	}

	b := new(bytes.Buffer)

	b.Write(SerializeUint(size))

	for i := 0; i < len(fields); i++ {
		b.Write(SerializeUint(offsets[i]))
	}

	for i := 0; i < len(fields); i++ {
		b.Write(fields[i])
	}

	return b.Bytes()
}

// SerializeOption serialize option
func SerializeOption(o Serializer) ([]byte, error) {
	if o == nil || reflect.ValueOf(o).IsNil() {
		return []byte{}, nil
	}

	return o.Serialize()
}

// SerializeOption serialize option
func SerializeOptionBytes(o []byte) ([]byte, error) {
	if o == nil || reflect.ValueOf(o).IsNil() {
		return []byte{}, nil
	}

	return SerializeBytes(o), nil
}
