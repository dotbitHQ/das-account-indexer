package types

type EpochParams struct {
	Length uint64
	Index  uint64
	Number uint64
}

func ParseEpoch(epoch uint64) *EpochParams {
	length := (epoch >> 40) & 0xFFFF
	index := (epoch >> 24) & 0xFFFF
	number := epoch & 0xFFFFFF
	return &EpochParams{
		Length: length,
		Index:  index,
		Number: number,
	}
}

func (ep *EpochParams) Uint64() uint64 {
	return (32 << 56) + (ep.Length << 40) + (ep.Index << 24) + ep.Number
}
