package pdu

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityStatePduParse(t *testing.T) {
	f, err := os.Open("testfiles/EntityStatePdu-26.raw")
	if err != nil {
		t.Fatalf("error opening file: %s", err)
	}

	defer f.Close()

	data, err := ParseEntityStatePDU(f)
	if err != nil {
		t.Fatalf("error parsing file: %s", err)
	}

	assert.Equal(t, uint8(6), data.Header.ProtocolVersion)
	assert.Equal(t, uint8(7), data.Header.ExerciseID)
	assert.Equal(t, uint8(1), data.Header.PDUType)
	assert.Equal(t, uint8(1), data.Header.ProtocolFamily)
	assert.Equal(t, uint32(2003426), data.Header.Timestamp)
	assert.Equal(t, uint16(144), data.Header.Length)
}

func TestEntityStatePduSerialize(t *testing.T) {
	f, err := os.Open("testfiles/EntityStatePdu-26.raw")
	if err != nil {
		t.Fatalf("error opening file: %s", err)
	}

	defer f.Close()

	data, err := ParseEntityStatePDU(f)
	if err != nil {
		t.Fatalf("error parsing file: %s", err)
	}

	var memStream bytes.Buffer
	err = data.Serialize(&memStream)
	assert.NoError(t, err)

	data1, err := ParseEntityStatePDU(&memStream)
	if err != nil {
		t.Fatalf("error parsing entity state: %s", err)
	}

	assert.Equal(t, data, data1)
}
