package impl

var magicBytes = []byte{0x27, 0x24, 0x50}

const (
	hasPrefixedLengthBit byte = 0x1
	isEmptyBit                = 0x2
)

// ProtocolType represents a known type that may be transferred using Ludwieg
type ProtocolType byte

const (
	// TypeUnknown represents an unknown type, often used as the zero-value of
	// a ProtocolType field
	TypeUnknown ProtocolType = 0x00

	// TypeUint8 represents an Uint8 type
	TypeUint8 = (0x01 << 2)

	// TypeUint32 represents an Uint32 type
	TypeUint32 = (0x02 << 2)

	// TypeUint64 represents an Uint64 type
	TypeUint64 = (0x03 << 2)

	// TypeFloat64 represents a Float64 type
	TypeFloat64 = (0x04 << 2)

	// TypeString represents a String type
	TypeString = (0x05 << 2) | 0x1

	// TypeBlob represents a Blob type
	TypeBlob = (0x06 << 2) | 0x1

	// TypeBool represents a Bool type
	TypeBool = (0x07 << 2)

	// TypeArray represents an Array type
	TypeArray = (0x08 << 2) | 0x1

	// TypeUUID represents an UUID type
	TypeUUID = (0x09 << 2)

	// TypeAny represents any type Ludwieg is capable of handling
	TypeAny = (0x0A << 2) | 0x1

	// TypeStruct is used internally to identify fields containing structs
	TypeStruct = (0x0B << 2) | 0x1
)

var knownTypes = []ProtocolType{
	TypeUint8, TypeUint32, TypeUint64, TypeFloat64, TypeString,
	TypeBlob, TypeBool, TypeArray, TypeUUID, TypeAny, TypeStruct,
}

type lengthEncoding byte

const (
	lengthEncodingEmpty lengthEncoding = iota
	lengthEncodingUint8
	lengthEncodingUint16
	lengthEncodingUint32
	lengthEncodingUint64
)
