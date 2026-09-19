package pdu

type EntityStatePDU struct {
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
	Marking                 [12]byte // 1-byte character set + 11-byte string/padding
	Capabilities            uint32
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
