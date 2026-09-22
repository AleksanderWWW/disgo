package pdu

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEntityStatePduParse(t *testing.T) {
	f, err := os.Open("testfiles/EntityStatePdu-26.raw")
	assert.NoError(t, err)

	defer f.Close()

	data, err := ParseEntityStatePDU(f)
	assert.NoError(t, err)

	assert.Equal(t, uint8(6), data.Base.Header.ProtocolVersion)
	assert.Equal(t, uint8(7), data.Base.Header.ExerciseID)
	assert.Equal(t, uint8(1), data.Base.Header.PDUType)
	assert.Equal(t, uint8(1), data.Base.Header.ProtocolFamily)
	assert.Equal(t, uint32(2003426), data.Base.Header.Timestamp.Value())
	assert.True(t, data.Base.Header.Timestamp.Relative())
	assert.Equal(t, uint16(144), data.Base.Header.Length)

	assert.Equal(t, uint8(1), data.Base.ForceId)
	assert.Equal(t, uint8(0), data.Base.NumVariableParameters)

	assert.Equal(t, "26", data.Base.Marking.String())
}

func TestEntityStatePduSerialize(t *testing.T) {
	f, err := os.Open("testfiles/EntityStatePdu-26.raw")
	assert.NoError(t, err)

	defer f.Close()

	data, err := ParseEntityStatePDU(f)
	assert.NoError(t, err)

	var memStream bytes.Buffer
	err = data.Serialize(&memStream)
	assert.NoError(t, err)

	data1, err := ParseEntityStatePDU(&memStream)
	assert.NoError(t, err)

	assert.Equal(t, data, data1)
}
