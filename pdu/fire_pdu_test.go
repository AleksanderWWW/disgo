package pdu

import (
	"bytes"
	"encoding/hex"
	"testing"
)

// Standard 96-byte IEEE 1278.1 Fire PDU Raw Byte Stream (Big-Endian)
// Header: Version=6, Exercise=1, Type=2 (Fire), Family=1, Timestamp=0x80000000 (Absolute), Length=96
var sampleFirePDUHex = "" +
	"06010201" + // ProtocolVersion (6), ExerciseID (1), PDUType (2 = Fire), ProtocolFamily (1)
	"80000000" + // Timestamp (Absolute Bit 0 set)
	"00600000" + // Length (96 bytes = 0x0060), PDUStatus (0)
	"00010001000A" + // Firing Entity ID: Site=1, App=1, Entity=10
	"000100010014" + // Target Entity ID: Site=1, App=1, Entity=20
	"00010001001E" + // Munition/Expendable ID: Site=1, App=1, Entity=30
	"000100010064" + // Event ID: Site=1, App=1, Event=100
	"00000001" + // Fire Mission Index = 1
	"408770000000000040877000000000004087700000000000" + // Location: WorldCoordinates (X, Y, Z float64)
	"0102000101010100" + // Burst Descriptor: EntityType (Kind=1, Domain=2, Country=1, Cat=1, Sub=1, Spec=1, Extra=0)
	"000100010001000A" + // Burst Descriptor: Warhead=1, Fuse=1, Quantity=1, Rate=10
	"3f8000003f8000003f800000" + // Velocity: Vector3Float (X=1.0, Y=1.0, Z=1.0)
	"447A0000" // Range: float32 (1000.0)

func TestParseFirePDU(t *testing.T) {
	rawBytes, err := hex.DecodeString(sampleFirePDUHex)
	if err != nil {
		t.Fatalf("Failed to decode test hex string: %v", err)
	}

	if len(rawBytes) != 96 {
		t.Fatalf("Expected raw payload length to be 96 bytes, got %d", len(rawBytes))
	}

	reader := bytes.NewReader(rawBytes)
	firePDU, err := ParseFirePDU(reader)
	if err != nil {
		t.Fatalf("ParseFirePDU failed: %v", err)
	}

	// Validate parsed properties
	if firePDU.Header.PDUType != PDUTypeFire {
		t.Errorf("Expected PDU Type %d, got %d", PDUTypeFire, firePDU.Header.PDUType)
	}

	if firePDU.FiringID.Number != 10 {
		t.Errorf("Expected Firing Entity ID 10, got %d", firePDU.FiringID.Number)
	}

	if firePDU.TargetID.Number != 20 {
		t.Errorf("Expected Target Entity ID 20, got %d", firePDU.TargetID.Number)
	}
}
