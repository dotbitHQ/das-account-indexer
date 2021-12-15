package utils

// define some useful const
// https://github.com/nervosnetwork/ckb/blob/35392279150fe4e61b7904516be91bda18c46f05/test/src/utils.rs#L24
const (
	FlagSinceRelative    = 0x8000000000000000
	FlagSinceEpochNumber = 0x2000000000000000
	FlagSinceBlockNumber = 0x0
	FlagSinceTimestamp   = 0x4000000000000000
)

func SinceFromRelativeBlockNumber(blockNumber uint64) uint64 {
	return FlagSinceRelative | FlagSinceBlockNumber | blockNumber
}

func SinceFromAbsoluteBlockNumber(blockNumber uint64) uint64 {
	return FlagSinceBlockNumber | blockNumber
}

func SinceFromRelativeEpochNumber(epochNumber uint64) uint64 {
	return FlagSinceRelative | FlagSinceEpochNumber | epochNumber
}

func SinceFromAbsoluteEpochNumber(epochNumber uint64) uint64 {
	return FlagSinceEpochNumber | epochNumber
}

func SinceFromRelativeTimestamp(timestamp uint64) uint64 {
	return FlagSinceRelative | FlagSinceTimestamp | timestamp
}

func SinceFromAbsoluteTimestamp(timestamp uint64) uint64 {
	return FlagSinceTimestamp | timestamp
}
