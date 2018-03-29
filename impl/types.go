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

	// TypeDouble represents a Float64 type
	TypeDouble = (0x04 << 2)

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

	// TypeDynInt represents an Integer value that may assume several sizes and
	// precisions (byte, uint16, uint32, uint64, float32, float64)
	TypeDynInt = (0x0C << 2)
)

var knownTypes = []ProtocolType{
	TypeUint8, TypeUint32, TypeUint64, TypeDouble, TypeString,
	TypeBlob, TypeBool, TypeArray, TypeUUID, TypeAny, TypeStruct,
	TypeDynInt,
}

type lengthEncoding byte

const (
	lengthEncodingEmpty lengthEncoding = iota
	lengthEncodingUint8
	lengthEncodingUint16
	lengthEncodingUint32
	lengthEncodingUint64
)

// DynIntValueKind represents which kind the current retained value have. It may
// be safe to assume larger values than the one returned, but any value smaller
// than the current one will cause an overflow
type DynIntValueKind byte

const (

	// DynIntValueKindInvalid is used internally to represent an invalid value
	DynIntValueKindInvalid DynIntValueKind = iota

	// DynIntValueKindUint8 represents an 8-bit unsigned integer
	DynIntValueKindUint8

	// DynIntValueKindUint16 represents a 16-bit unsigned integer
	DynIntValueKindUint16

	// DynIntValueKindUint32 represents a 32-bit unsigned integer
	DynIntValueKindUint32

	// DynIntValueKindUint64 represents a 64-bit unsigned integer
	DynIntValueKindUint64

	// DynIntValueKindInt8 represents an 8-bit integer
	DynIntValueKindInt8

	// DynIntValueKindInt16 represents a 16-bit integer
	DynIntValueKindInt16

	// DynIntValueKindInt32 represents a 32-bit integer
	DynIntValueKindInt32

	// DynIntValueKindInt64 represents a 64-bit integer
	DynIntValueKindInt64

	// DynIntValueKindFloat32 represents a 32-bit float value
	DynIntValueKindFloat32

	// DynIntValueKindFloat64 represents an 64-bit float value
	DynIntValueKindFloat64
)
