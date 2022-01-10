package types

import (
	"bytes"
	"errors"
)

func (h Hash) Serialize() ([]byte, error) {
	return h.Bytes(), nil
}

func (t ScriptHashType) Serialize() ([]byte, error) {
	if t == HashTypeData {
		return []byte{00}, nil
	} else if t == HashTypeType {
		return []byte{01}, nil
	} else if t == HashTypeData1 {
		return []byte{02}, nil
	}
	return nil, errors.New("invalid script hash type")
}

// Serialize dep type
func (t DepType) Serialize() ([]byte, error) {
	if t == DepTypeCode {
		return []byte{00}, nil
	} else if t == DepTypeDepGroup {
		return []byte{01}, nil
	}
	return nil, errors.New("invalid dep group")
}

// Serialize script
func (script *Script) Serialize() ([]byte, error) {
	h, err := script.CodeHash.Serialize()
	if err != nil {
		return nil, err
	}

	t, err := script.HashType.Serialize()
	if err != nil {
		return nil, err
	}

	a := SerializeBytes(script.Args)

	return SerializeTable([][]byte{h, t, a}), nil
}

// Serialize outpoint
func (o *OutPoint) Serialize() ([]byte, error) {
	h, err := o.TxHash.Serialize()
	if err != nil {
		return nil, err
	}

	i := SerializeUint(o.Index)

	b := new(bytes.Buffer)

	b.Write(h)
	b.Write(i)

	return b.Bytes(), nil
}

// Serialize cell input
func (i *CellInput) Serialize() ([]byte, error) {
	s := SerializeUint64(i.Since)

	o, err := i.PreviousOutput.Serialize()
	if err != nil {
		return nil, err
	}

	return SerializeStruct([][]byte{s, o}), nil
}

// Serialize cell output
func (o *CellOutput) Serialize() ([]byte, error) {
	c := SerializeUint64(o.Capacity)

	l, err := o.Lock.Serialize()
	if err != nil {
		return nil, err
	}

	t, err := SerializeOption(o.Type)
	if err != nil {
		return nil, err
	}

	return SerializeTable([][]byte{c, l, t}), nil
}

// Serialize cell dep
func (d *CellDep) Serialize() ([]byte, error) {
	o, err := d.OutPoint.Serialize()
	if err != nil {
		return nil, err
	}

	dd, err := d.DepType.Serialize()
	if err != nil {
		return nil, err
	}

	return SerializeStruct([][]byte{o, dd}), nil
}

// Serialize transaction
func (t *Transaction) Serialize() ([]byte, error) {
	v := SerializeUint(t.Version)

	// Ok, no way around this
	deps := make([]Serializer, len(t.CellDeps))
	for i, v := range t.CellDeps {
		deps[i] = v
	}
	cds, err := SerializeArray(deps)
	if err != nil {
		return nil, err
	}
	cdsBytes := SerializeFixVec(cds)

	hds := make([][]byte, len(t.HeaderDeps))
	for i := 0; i < len(t.HeaderDeps); i++ {
		hd, err := t.HeaderDeps[i].Serialize()
		if err != nil {
			return nil, err
		}

		hds[i] = hd
	}
	hdsBytes := SerializeFixVec(hds)

	ips := make([][]byte, len(t.Inputs))
	for i := 0; i < len(t.Inputs); i++ {
		ip, err := t.Inputs[i].Serialize()
		if err != nil {
			return nil, err
		}

		ips[i] = ip
	}
	ipsBytes := SerializeFixVec(ips)

	ops := make([][]byte, len(t.Outputs))
	for i := 0; i < len(t.Outputs); i++ {
		op, err := t.Outputs[i].Serialize()
		if err != nil {
			return nil, err
		}

		ops[i] = op
	}
	opsBytes := SerializeDynVec(ops)

	ods := make([][]byte, len(t.OutputsData))
	for i := 0; i < len(t.OutputsData); i++ {
		od := SerializeBytes(t.OutputsData[i])

		ods[i] = od
	}
	odsBytes := SerializeDynVec(ods)

	fields := [][]byte{v, cdsBytes, hdsBytes, ipsBytes, opsBytes, odsBytes}
	return SerializeTable(fields), nil
}

func (w *WitnessArgs) Serialize() ([]byte, error) {
	l, err := SerializeOptionBytes(w.Lock)
	if err != nil {
		return nil, err
	}

	i, err := SerializeOptionBytes(w.InputType)
	if err != nil {
		return nil, err
	}

	o, err := SerializeOptionBytes(w.OutputType)
	if err != nil {
		return nil, err
	}

	return SerializeTable([][]byte{l, i, o}), nil
}

func SerializeHashType(hashType ScriptHashType) (string, error) {
	if HashTypeData == hashType {
		return "00", nil
	} else if HashTypeType == hashType {
		return "01", nil
	} else if HashTypeData1 == hashType {
		return "02", nil
	}

	return "", errors.New("Invalid script hash_type: " + string(hashType))
}

func DeserializeHashType(hashType string) (ScriptHashType, error) {
	if "00" == hashType {
		return HashTypeData, nil
	} else if "01" == hashType {
		return HashTypeType, nil
	} else if "02" == hashType {
		return HashTypeData1, nil
	}

	return "", errors.New("Invalid script hash_type: " + hashType)
}
