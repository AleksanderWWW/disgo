package pdu

import (
	"time"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const (
	AbsoluteFlagMask    = 0x80000000 // 1000 0000 ... in binary
	TimestampValueMask = 0x7FFFFFFF // 0111 1111 ... in binary

	// unitsPerHour is 2^31, the total units in one hour.
	unitsPerHour = 2147483648
	// nsPerUnit is 3600 seconds / 2^31 expressed in nanoseconds.
	nsPerUnit = 1676.3806
)

// GetCurrentTimestamp uses bitwise operators and masks to build
// a semantic DIS timestamp from the system clock.
func GetCurrentTimestamp(absolute bool) EntityTimestamp {
	// 1. Get current time
	now := time.Now()

	// 2. Find start of the current hour
	currentHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())

	// 3. Calculate time passed since top of the hour (in nanoseconds)
	nsSinceHour := now.Sub(currentHour).Nanoseconds()

	// 4. Convert nanoseconds to arbitrary DIS units
	// timeUnits = (nanoseconds / nanoseconds_per_unit)
	unitsPassed := uint32(float64(nsSinceHour) / nsPerUnit)

	// --- Final Assembly Using Bitwise Ops ---

	// Start by ensuring only the lower 31 bits are populated
	// 0x7FFFFFFF ensures the top bit is clear (0).
	timestampVal := EntityTimestamp(unitsPassed & TimestampValueMask)

	// If absolute time is requested, use bitwise OR to FORCE the top bit to 1.
	// We don't need a relative case because ORing with 0 has no effect.
	if absolute {
		// (1xxxx...) | (10000...) = (1xxxx...)
		timestampVal = timestampVal | EntityTimestamp(AbsoluteFlagMask)
	}

	return timestampVal
}


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
	Timestamp       EntityTimestamp
	Length          uint16
	PDUStatus       uint16
}

type EntityTimestamp uint32


func (et EntityTimestamp) Absolute() bool {
	// 0x80000000 is 1000 0000... in binary (Bit 0 set).
	// Since we are big-endian in DIS, Bit 0 is the most significant bit.
	const AbsoluteFlagMask = 0x80000000
	return (uint32(et) & AbsoluteFlagMask) != 0
}

func (et EntityTimestamp) Relative() bool {
	return !et.Absolute()
}

func (et EntityTimestamp) Value() uint32 {
	// 0x7FFFFFFF clears Bit 0, leaving Bits 1-31.
	const TimestampValueMask = 0x7FFFFFFF
	return uint32(et) & TimestampValueMask
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
