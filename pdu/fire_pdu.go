package pdu

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Fire PDU - 96 bytes
type FirePDU struct {
	Header               PDUHeader
	FiringID             EntityId
	TargetID             EntityId
	MunitionExpendableId EntityId
	EventID              EventId
	FireMissionIdx       uint32
	Location             WorldCoordinates
	Description          MunitionDescription
	Velocity             Vector3Float
	Range                float32
}

func ParseFirePDU(reader io.Reader) (*FirePDU, error) {
	var pdu FirePDU

	// Fixed 96-byte single pass read
	if err := binary.Read(reader, binary.BigEndian, &pdu); err != nil {
		return nil, fmt.Errorf("failed to read FirePDU: %w", err)
	}

	if pdu.Header.PDUType != PDUTypeFire {
		return nil, fmt.Errorf("invalid PDU type: expected %d, got %d", PDUTypeFire, pdu.Header.PDUType)
	}

	return &pdu, nil
}

func (pdu FirePDU) Serialize(w io.Writer) error {
	return binary.Write(w, binary.BigEndian, pdu)
}
