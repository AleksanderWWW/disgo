package pdu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

func ParseEntityStatePDU(reader io.Reader) (*EntityStatePDU, error) {
	var pdu EntityStatePDU

	err := binary.Read(reader, binary.BigEndian, &pdu.Base)
	if err != nil {
		return nil, err
	}

	if pdu.Base.Header.PDUType != PDUTypeEntityState {
		return nil, fmt.Errorf("invalid PDU type: expected: 1, got: %d", pdu.Base.Header.PDUType)
	}

	numParams := int(pdu.Base.NumVariableParameters)
	if numParams > 0 {
		pdu.VariableParameters.Params = make([]VariableParameter, numParams)

		for i := range numParams {
			var vp VariableParameter
			if err := binary.Read(reader, binary.BigEndian, &vp); err != nil {
				return nil, fmt.Errorf("failed to read VariableParameter[%d]: %w", i, err)
			}
			pdu.VariableParameters.Params[i] = vp
		}
	}

	return &pdu, nil
}

type EntityStatePDU struct {
	Base               EntityStateBase
	VariableParameters VariableParameterList
}

type EntityStateBase struct {
	Header                  EntityHeader
	Id                      EntityId
	ForceId                 uint8
	NumVariableParameters   uint8
	Type                    EntityType
	AlternativeType         EntityType
	LinearVelocity          Vector3Float
	Location                WorldCoordinates
	Orientation             Vector3Float
	Appearance              uint32
	DeadReckoningParameters [40]byte // 40-byte fixed block
	Marking                 EntityMarking
	Capabilities            uint32
}

type VariableParameterList struct {
	Params []VariableParameter
}

type VariableParameter struct {
	RecordType   uint8
	RecordLength uint8
	Data         [14]byte
}

type EntityMarking struct {
	Val [12]byte // 1-byte character set + 11-byte string/padding
}

func (em EntityMarking) String() string {
	trimmed := string(bytes.Trim(em.Val[:], "\x00"))

	if len(trimmed) == 0 {
		return ""
	}
	return trimmed[1:]
}

type EntityHeader struct {
	ProtocolVersion uint8
	ExerciseID      uint8
	PDUType         uint8
	ProtocolFamily  uint8
	Timestamp       uint32
	Length          uint16
	PDUStatus       uint16
}

type WorldCoordinates struct {
	X float64
	Y float64
	Z float64
}

type Vector3Float struct {
	X float32
	Y float32
	Z float32
}

type EntityType struct {
	Kind        uint8
	Domain      uint8
	Country     uint16
	Category    uint8
	Subcategory uint8
	Specific    uint8
	Extra       uint8
}

type EntityId struct {
	Address SimulationAddress
	Number  uint16
}

type SimulationAddress struct {
	Site        uint16
	Application uint16
}

func (pdu EntityStatePDU) Serialize(w io.Writer) error {
	err := binary.Write(w, binary.BigEndian, pdu.Base)
	if err != nil {
		return err
	}

	for i, vp := range pdu.VariableParameters.Params {
		err = binary.Write(w, binary.BigEndian, vp)
		if err != nil {
			return fmt.Errorf("failed to serialize VariableParameter[%d]: %w", i, err)
		}
	}

	return nil
}
