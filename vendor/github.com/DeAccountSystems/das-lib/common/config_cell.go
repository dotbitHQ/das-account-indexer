package common

type ConfigCellTypeArgs = string

const (
	ConfigCellTypeArgsAccount         ConfigCellTypeArgs = "0x64000000"
	ConfigCellTypeArgsApply           ConfigCellTypeArgs = "0x65000000"
	ConfigCellTypeArgsIncome          ConfigCellTypeArgs = "0x67000000"
	ConfigCellTypeArgsMain            ConfigCellTypeArgs = "0x68000000"
	ConfigCellTypeArgsPrice           ConfigCellTypeArgs = "0x69000000"
	ConfigCellTypeArgsProposal        ConfigCellTypeArgs = "0x6a000000"
	ConfigCellTypeArgsProfitRate      ConfigCellTypeArgs = "0x6b000000"
	ConfigCellTypeArgsRecordNamespace ConfigCellTypeArgs = "0x6c000000"
	ConfigCellTypeArgsRelease         ConfigCellTypeArgs = "0x6d000000"
	ConfigCellTypeArgsUnavailable     ConfigCellTypeArgs = "0x6e000000"
	ConfigCellTypeArgsSecondaryMarket ConfigCellTypeArgs = "0x6f000000"
	ConfigCellTypeArgsReverseRecord   ConfigCellTypeArgs = "0x70000000"
	ConfigCellTypeArgsSubAccount      ConfigCellTypeArgs = "0x71000000"

	ConfigCellTypeArgsPreservedAccount00 ConfigCellTypeArgs = "0x10270000"
	ConfigCellTypeArgsPreservedAccount01 ConfigCellTypeArgs = "0x11270000"
	ConfigCellTypeArgsPreservedAccount02 ConfigCellTypeArgs = "0x12270000"
	ConfigCellTypeArgsPreservedAccount03 ConfigCellTypeArgs = "0x13270000"
	ConfigCellTypeArgsPreservedAccount04 ConfigCellTypeArgs = "0x14270000"
	ConfigCellTypeArgsPreservedAccount05 ConfigCellTypeArgs = "0x15270000"
	ConfigCellTypeArgsPreservedAccount06 ConfigCellTypeArgs = "0x16270000"
	ConfigCellTypeArgsPreservedAccount07 ConfigCellTypeArgs = "0x17270000"
	ConfigCellTypeArgsPreservedAccount08 ConfigCellTypeArgs = "0x18270000"
	ConfigCellTypeArgsPreservedAccount09 ConfigCellTypeArgs = "0x19270000"
	ConfigCellTypeArgsPreservedAccount10 ConfigCellTypeArgs = "0x1a270000"
	ConfigCellTypeArgsPreservedAccount11 ConfigCellTypeArgs = "0x1b270000"
	ConfigCellTypeArgsPreservedAccount12 ConfigCellTypeArgs = "0x1c270000"
	ConfigCellTypeArgsPreservedAccount13 ConfigCellTypeArgs = "0x1d270000"
	ConfigCellTypeArgsPreservedAccount14 ConfigCellTypeArgs = "0x1e270000"
	ConfigCellTypeArgsPreservedAccount15 ConfigCellTypeArgs = "0x1f270000"
	ConfigCellTypeArgsPreservedAccount16 ConfigCellTypeArgs = "0x20270000"
	ConfigCellTypeArgsPreservedAccount17 ConfigCellTypeArgs = "0x21270000"
	ConfigCellTypeArgsPreservedAccount18 ConfigCellTypeArgs = "0x22270000"
	ConfigCellTypeArgsPreservedAccount19 ConfigCellTypeArgs = "0x23270000"

	ConfigCellTypeArgsCharSetEmoji ConfigCellTypeArgs = "0xa0860100"
	ConfigCellTypeArgsCharSetDigit ConfigCellTypeArgs = "0xa1860100"
	ConfigCellTypeArgsCharSetEn    ConfigCellTypeArgs = "0xa2860100"
	ConfigCellTypeArgsCharSetHanS  ConfigCellTypeArgs = "0xa3860100"
	ConfigCellTypeArgsCharSetHanT  ConfigCellTypeArgs = "0xa4860100"
)

func GetConfigCellTypeArgsPreservedAccountByIndex(index uint32) ConfigCellTypeArgs {
	switch index {
	case 0:
		return ConfigCellTypeArgsPreservedAccount00
	case 1:
		return ConfigCellTypeArgsPreservedAccount01
	case 2:
		return ConfigCellTypeArgsPreservedAccount02
	case 3:
		return ConfigCellTypeArgsPreservedAccount03
	case 4:
		return ConfigCellTypeArgsPreservedAccount04
	case 5:
		return ConfigCellTypeArgsPreservedAccount05
	case 6:
		return ConfigCellTypeArgsPreservedAccount06
	case 7:
		return ConfigCellTypeArgsPreservedAccount07
	case 8:
		return ConfigCellTypeArgsPreservedAccount08
	case 9:
		return ConfigCellTypeArgsPreservedAccount09
	case 10:
		return ConfigCellTypeArgsPreservedAccount10
	case 11:
		return ConfigCellTypeArgsPreservedAccount11
	case 12:
		return ConfigCellTypeArgsPreservedAccount12
	case 13:
		return ConfigCellTypeArgsPreservedAccount13
	case 14:
		return ConfigCellTypeArgsPreservedAccount14
	case 15:
		return ConfigCellTypeArgsPreservedAccount15
	case 16:
		return ConfigCellTypeArgsPreservedAccount16
	case 17:
		return ConfigCellTypeArgsPreservedAccount17
	case 18:
		return ConfigCellTypeArgsPreservedAccount18
	case 19:
		return ConfigCellTypeArgsPreservedAccount19
	}
	return ""
}
