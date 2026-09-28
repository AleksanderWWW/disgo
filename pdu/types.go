package pdu


// PDU Header - 12 bytes
type PDUHeader struct {
	ProtocolVersion uint8
	ExerciseID      uint8
	PDUType         uint8
	ProtocolFamily  uint8
	Timestamp       EntityTimestamp
	Length          uint16
	PDUStatus       uint16
}

// Entity type - 8 bytes
type EntityType struct {
	Kind        uint8
	Domain      uint8
	Country     uint16
	Category    uint8
	Subcategory uint8
	Specific    uint8
	Extra       uint8
}

// Munition Description type - 16 bytes
type MunitionDescription struct {
	Type EntityType
	Warhead uint16
	Fuse uint16
	Quantity uint16
	Rate uint16
}

// Entity ID type - 6 bytes
type EntityId struct {
	Address SimulationAddress
	Number  uint16
}

// Event ID type - 6 bytes
type EventId struct {
	Address SimulationAddress
	Number  uint16
}

// Simulation Address type - 4 bytes
type SimulationAddress struct {
	Site        uint16
	Application uint16
}

// World Coordinates type - 24 bytes
type WorldCoordinates struct {
	X float64
	Y float64
	Z float64
}

// 3-element vector (float elements) - 12 bytes
type Vector3Float struct {
	X float32
	Y float32
	Z float32
}
