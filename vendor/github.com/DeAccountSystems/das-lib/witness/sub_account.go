package witness

import (
	"fmt"
	"github.com/DeAccountSystems/das-lib/common"
	"github.com/DeAccountSystems/das-lib/molecule"
	"github.com/nervosnetwork/ckb-sdk-go/crypto/blake2b"
	"github.com/nervosnetwork/ckb-sdk-go/types"
)

const (
	SubAccountCurrentVersion = common.GoDataEntityVersion1
)

type SubAccountBuilder struct {
	Signature         []byte
	SignRole          []byte
	PrevRoot          []byte
	CurrentRoot       []byte
	Proof             []byte
	Version           uint32
	SubAccount        *SubAccount
	EditKey           []byte
	EditValue         []byte
	Account           string
	CurrentSubAccount *SubAccount
}

type SubAccountParam struct {
	Signature      []byte
	SignRole       []byte
	PrevRoot       []byte
	CurrentRoot    []byte
	Proof          []byte
	SubAccount     *SubAccount
	EditKey        string
	EditLockArgs   []byte
	EditRecords    []SubAccountRecord
	RenewExpiredAt uint64
}

type SubAccount struct {
	Lock                 *types.Script           `json:"lock"`
	AccountId            string                  `json:"account_id"`
	AccountCharSet       []common.AccountCharSet `json:"account_char_set"`
	Suffix               string                  `json:"suffix"`
	RegisteredAt         uint64                  `json:"registered_at"`
	ExpiredAt            uint64                  `json:"expired_at"`
	Status               uint8                   `json:"status"`
	Records              []SubAccountRecord      `json:"records"`
	Nonce                uint64                  `json:"nonce"`
	EnableSubAccount     uint8                   `json:"enable_sub_account"`
	RenewSubAccountPrice uint64                  `json:"renew_sub_account_price"`
}

type SubAccountRecord struct {
	Key   string `json:"key"`
	Type  string `json:"type"`
	Label string `json:"label"`
	Value string `json:"value"`
	TTL   uint32 `json:"ttl"`
}

type SubAccountEditValue struct {
	LockArgs  string             `json:"lock_args"`
	Records   []SubAccountRecord `json:"records"`
	ExpiredAt uint64             `json:"expired_at"`
}

func SubAccountBuilderFromTx(tx *types.Transaction) (*SubAccountBuilder, error) {
	respMap, err := SubAccountBuilderMapFromTx(tx)
	if err != nil {
		return nil, err
	}
	for k, _ := range respMap {
		return respMap[k], nil
	}
	return nil, fmt.Errorf("not exist sub account")
}

func SubAccountBuilderMapFromTx(tx *types.Transaction) (map[string]*SubAccountBuilder, error) {
	var respMap = make(map[string]*SubAccountBuilder)

	err := GetWitnessDataFromTx(tx, func(actionDataType common.ActionDataType, dataBys []byte) (bool, error) {
		switch actionDataType {
		case common.ActionDataTypeSubAccount:
			builder, err := SubAccountBuilderFromBytes(dataBys)
			if err != nil {
				return false, err
			}

			currentSubAccount := *builder.SubAccount
			builder.CurrentSubAccount = &currentSubAccount

			editKey := string(builder.EditKey)
			if editKey != "" {
				builder.CurrentSubAccount.Nonce++
			}
			switch editKey {
			case common.EditKeyOwner:
				builder.CurrentSubAccount.Lock = &types.Script{
					CodeHash: builder.SubAccount.Lock.CodeHash,
					HashType: builder.SubAccount.Lock.HashType,
					Args:     builder.EditValue,
				}
				builder.CurrentSubAccount.Records = nil
			case common.EditKeyManager:
				builder.CurrentSubAccount.Lock = &types.Script{
					CodeHash: builder.SubAccount.Lock.CodeHash,
					HashType: builder.SubAccount.Lock.HashType,
					Args:     builder.EditValue,
				}
			case common.EditKeyRecords:
				records := builder.ConvertEditValueToRecords()
				builder.CurrentSubAccount.Records = ConvertToSubAccountRecords(records)
			case common.EditKeyExpiredAt:
				expiredAt := builder.ConvertEditValueToExpiredAt()
				builder.CurrentSubAccount.ExpiredAt, _ = molecule.Bytes2GoU64(expiredAt.RawData())
			}

			respMap[builder.SubAccount.AccountId] = builder
		}
		return true, nil
	})

	if err != nil {
		return nil, fmt.Errorf("GetWitnessDataFromTx err: %s", err.Error())
	}
	if len(respMap) == 0 {
		return nil, fmt.Errorf("not exist sub account")
	}
	return respMap, nil
}

func SubAccountBuilderFromBytes(dataBys []byte) (*SubAccountBuilder, error) {
	var resp SubAccountBuilder
	index, length := uint32(0), uint32(4)

	signatureLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.Signature = dataBys[index : index+signatureLen]
	index += signatureLen

	signRoleLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.SignRole = dataBys[index : index+signRoleLen]
	index += signRoleLen

	prevRootLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.PrevRoot = dataBys[index : index+prevRootLen]
	index += prevRootLen

	currentRootLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.CurrentRoot = dataBys[index : index+currentRootLen]
	index += currentRootLen

	proofLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.Proof = dataBys[index : index+proofLen]
	index += proofLen

	versionLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.Version, _ = molecule.Bytes2GoU32(dataBys[index : index+versionLen])
	index += versionLen

	subAccountLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	subAccountBys := dataBys[index : index+subAccountLen]
	index += subAccountLen

	keyLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.EditKey = dataBys[index : index+keyLen]
	index += keyLen

	valueLen, _ := molecule.Bytes2GoU32(dataBys[index : index+length])
	index += length
	resp.EditValue = dataBys[index : index+valueLen]
	index += valueLen

	switch resp.Version {
	case common.GoDataEntityVersion1:
		subAccount, err := ConvertToSubAccount(subAccountBys)
		if err != nil {
			return nil, fmt.Errorf("ConvertToSubAccount err: %s", err.Error())
		}
		resp.SubAccount = subAccount
		resp.Account = subAccount.Account()
		return &resp, nil
	default:
		return nil, fmt.Errorf("sub account version: %d", resp.Version)
	}
}

func (s *SubAccountBuilder) ConvertToEditValue() (*SubAccountEditValue, error) {
	var editValue SubAccountEditValue
	editKey := string(s.EditKey)
	switch editKey {
	case common.EditKeyOwner, common.EditKeyManager:
		editValue.LockArgs = common.Bytes2Hex(s.EditValue)
	case common.EditKeyRecords:
		records := s.ConvertEditValueToRecords()
		editValue.Records = ConvertToSubAccountRecords(records)
	case common.EditKeyExpiredAt:
		expiredAt := s.ConvertEditValueToExpiredAt()
		editValue.ExpiredAt, _ = molecule.Bytes2GoU64(expiredAt.RawData())
	default:
		return nil, fmt.Errorf("not support edit key[%s]", editKey)
	}
	return &editValue, nil
}

func (s *SubAccountBuilder) ConvertEditValueToExpiredAt() *molecule.Uint64 {
	expiredAt, _ := molecule.Uint64FromSlice(s.EditValue, false)
	return expiredAt
}

func (s *SubAccountBuilder) ConvertEditValueToRecords() *molecule.Records {
	records, _ := molecule.RecordsFromSlice(s.EditValue, false)
	return records
}

func ConvertToSubAccountRecords(records *molecule.Records) []SubAccountRecord {
	var subAccountRecords []SubAccountRecord
	for index, lenRecords := uint(0), records.Len(); index < lenRecords; index++ {
		record := records.Get(index)
		ttl, _ := molecule.Bytes2GoU32(record.RecordTtl().RawData())
		subAccountRecords = append(subAccountRecords, SubAccountRecord{
			Key:   string(record.RecordKey().RawData()),
			Type:  string(record.RecordType().RawData()),
			Label: string(record.RecordLabel().RawData()),
			Value: string(record.RecordValue().RawData()),
			TTL:   ttl,
		})
	}
	return subAccountRecords
}

func ConvertToAccountCharSets(accountChars *molecule.AccountChars) []common.AccountCharSet {
	index := uint(0)
	var accountCharSets []common.AccountCharSet
	for ; index < accountChars.ItemCount(); index++ {
		char := accountChars.Get(index)
		charSetName, _ := molecule.Bytes2GoU32(char.CharSetName().RawData())
		accountCharSets = append(accountCharSets, common.AccountCharSet{
			CharSetName: common.AccountCharType(charSetName),
			Char:        string(char.Bytes().RawData()),
		})
	}
	return accountCharSets
}

/****************************************** Parting Line ******************************************/

func ConvertToSubAccount(subAccountBys []byte) (*SubAccount, error) {
	subAccount, err := molecule.SubAccountFromSlice(subAccountBys, false)
	if err != nil {
		return nil, fmt.Errorf("SubAccountDataFromSlice err: %s", err.Error())
	}
	var tmp SubAccount
	tmp.Lock = molecule.MoleculeScript2CkbScript(subAccount.Lock())
	tmp.AccountId = common.Bytes2Hex(subAccount.Id().RawData())
	tmp.AccountCharSet = ConvertToAccountCharSets(subAccount.Account())
	tmp.Suffix = string(subAccount.Suffix().RawData())
	tmp.RegisteredAt, _ = molecule.Bytes2GoU64(subAccount.RegisteredAt().RawData())
	tmp.ExpiredAt, _ = molecule.Bytes2GoU64(subAccount.ExpiredAt().RawData())
	tmp.Status, _ = molecule.Bytes2GoU8(subAccount.Status().RawData())
	tmp.Records = ConvertToSubAccountRecords(subAccount.Records())
	tmp.Nonce, _ = molecule.Bytes2GoU64(subAccount.Nonce().RawData())
	tmp.EnableSubAccount, _ = molecule.Bytes2GoU8(subAccount.EnableSubAccount().RawData())
	tmp.RenewSubAccountPrice, _ = molecule.Bytes2GoU64(subAccount.RenewSubAccountPrice().RawData())

	return &tmp, nil
}

func ConvertToAccountChars(accountCharSet []common.AccountCharSet) *molecule.AccountChars {
	accountCharsBuilder := molecule.NewAccountCharsBuilder()
	for _, item := range accountCharSet {
		if item.Char == "." {
			break
		}
		accountChar := molecule.NewAccountCharBuilder().
			CharSetName(molecule.GoU32ToMoleculeU32(uint32(item.CharSetName))).
			Bytes(molecule.GoBytes2MoleculeBytes([]byte(item.Char))).Build()
		accountCharsBuilder.Push(accountChar)
	}
	accountChars := accountCharsBuilder.Build()
	return &accountChars
}

func ConvertToRecords(subAccountRecords []SubAccountRecord) *molecule.Records {
	recordsBuilder := molecule.NewRecordsBuilder()
	for _, v := range subAccountRecords {
		record := molecule.RecordDefault()
		recordBuilder := record.AsBuilder()
		recordBuilder.RecordKey(molecule.GoString2MoleculeBytes(v.Key)).
			RecordType(molecule.GoString2MoleculeBytes(v.Type)).
			RecordLabel(molecule.GoString2MoleculeBytes(v.Label)).
			RecordValue(molecule.GoString2MoleculeBytes(v.Value)).
			RecordTtl(molecule.GoU32ToMoleculeU32(v.TTL))
		recordsBuilder.Push(recordBuilder.Build())
	}
	records := recordsBuilder.Build()
	return &records
}

func (s *SubAccount) ConvertToMoleculeSubAccount() *molecule.SubAccount {
	lock := molecule.CkbScript2MoleculeScript(s.Lock)
	accountChars := ConvertToAccountChars(s.AccountCharSet)
	accountId, _ := molecule.AccountIdFromSlice(common.Hex2Bytes(s.AccountId), false)
	suffix := molecule.GoBytes2MoleculeBytes([]byte(s.Suffix))
	registeredAt := molecule.GoU64ToMoleculeU64(s.RegisteredAt)
	expiredAt := molecule.GoU64ToMoleculeU64(s.ExpiredAt)
	status := molecule.GoU8ToMoleculeU8(s.Status)
	records := ConvertToRecords(s.Records)
	nonce := molecule.GoU64ToMoleculeU64(s.Nonce)
	enableSubAccount := molecule.GoU8ToMoleculeU8(s.EnableSubAccount)
	renewSubAccountPrice := molecule.GoU64ToMoleculeU64(s.RenewSubAccountPrice)

	moleculeSubAccount := molecule.NewSubAccountBuilder().
		Lock(lock).
		Id(*accountId).
		Account(*accountChars).
		Suffix(suffix).
		RegisteredAt(registeredAt).
		ExpiredAt(expiredAt).
		Status(status).
		Records(*records).
		Nonce(nonce).
		EnableSubAccount(enableSubAccount).
		RenewSubAccountPrice(renewSubAccountPrice).
		Build()
	return &moleculeSubAccount
}

func (s *SubAccount) Account() string {
	var account string
	for _, v := range s.AccountCharSet {
		account += v.Char
	}
	return account + s.Suffix
}

func (s *SubAccount) ToH256() []byte {
	moleculeSubAccount := s.ConvertToMoleculeSubAccount()
	bys, _ := blake2b.Blake256(moleculeSubAccount.AsSlice())
	return bys
}

func (p *SubAccountParam) GenSubAccountBytes() (bys []byte) {
	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.Signature)))...)
	bys = append(bys, p.Signature...)

	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.SignRole)))...)
	bys = append(bys, p.SignRole...)

	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.PrevRoot)))...)
	bys = append(bys, p.PrevRoot...)

	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.CurrentRoot)))...)
	bys = append(bys, p.CurrentRoot...)

	bys = append(bys, molecule.GoU32ToBytes(uint32(len(p.Proof)))...)
	bys = append(bys, p.Proof...)

	versionBys := molecule.GoU32ToMoleculeU32(SubAccountCurrentVersion)
	bys = append(bys, molecule.GoU32ToBytes(uint32(len(versionBys.RawData())))...)
	bys = append(bys, versionBys.RawData()...)

	subAccount := p.SubAccount.ConvertToMoleculeSubAccount()
	bys = append(bys, molecule.GoU32ToBytes(uint32(len(subAccount.AsSlice())))...)
	bys = append(bys, subAccount.AsSlice()...)

	bys = append(bys, molecule.GoU32ToBytes(uint32(len([]byte(p.EditKey))))...)
	bys = append(bys, p.EditKey...)

	var editValue []byte
	switch p.EditKey {
	case common.EditKeyOwner, common.EditKeyManager:
		editValue = p.EditLockArgs
	case common.EditKeyRecords:
		records := ConvertToRecords(p.EditRecords)
		editValue = records.AsSlice()
	case common.EditKeyExpiredAt:
		expiredAt := molecule.GoU64ToMoleculeU64(p.RenewExpiredAt)
		editValue = expiredAt.AsSlice()
	}

	bys = append(bys, molecule.GoU32ToBytes(uint32(len(editValue)))...)
	bys = append(bys, editValue...)
	return
}

func (p *SubAccountParam) NewSubAccountWitness() ([]byte, error) {
	bys := p.GenSubAccountBytes()
	witness := GenDasDataWitnessWithByte(common.ActionDataTypeSubAccount, bys)
	return witness, nil
}

func ConvertSubAccountCellOutputData(data []byte) (smtRoot []byte, profit uint64) {
	if len(data) == 32 {
		smtRoot = data
	} else if len(data) == 40 {
		smtRoot = data[:32]
		profit, _ = molecule.Bytes2GoU64(data[32:])
	}
	return
}

func BuildSubAccountCellOutputData(smtRoot []byte, profit uint64) []byte {
	data := molecule.GoU64ToMoleculeU64(profit)
	smtRoot = append(smtRoot, data.RawData()...)
	return smtRoot
}
