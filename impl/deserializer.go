package impl

import (
	"fmt"
	"reflect"
)

type registeredPackage struct {
	nativeType reflect.Type
	id         byte
	meta       []LudwiegTypeAnnotation
}

// RegisterPackages is responsible for registering a known package under a
// message ID. It is then used afterwards to decode incoming messages. This
// method is called automatically, when generating sources using the "ludco"
// command line utility.
func RegisterPackages(pkgs ...SerializablePackage) {
	for _, pkg := range pkgs {
		id := pkg.LudwiegID()
		if _, ok := registeredPackages[id]; ok {
			panic(fmt.Errorf("illegal attempt to register two packages with same id: %#v", id))
		}

		registeredPackages[id] = registeredPackage{
			id:         id,
			meta:       pkg.LudwiegMeta(),
			nativeType: reflect.TypeOf(pkg),
		}
	}
}

type typeDecoderFunc func(t metaProtocolByte, b []byte, offset *int) (interface{}, error)

var registeredPackages = map[byte]registeredPackage{}
var registeredTypeDecoder map[ProtocolType]typeDecoderFunc

func init() {
	registeredTypeDecoder = map[ProtocolType]typeDecoderFunc{
		TypeUnknown: decodeUnknown,
		TypeUint8:   decodeUint8,
		TypeUint32:  decodeUint32,
		TypeUint64:  decodeUint64,
		TypeDouble:  decodeDouble,
		TypeString:  decodeString,
		TypeBlob:    decodeBlob,
		TypeBool:    decodeBool,
		TypeUUID:    decodeUUID,
		TypeAny:     decodeAny,
		TypeArray:   decodeArray,
		TypeStruct:  decodeStruct,
		TypeDynInt:  decodeDynInt,
	}
}

// Deserializer

type deserializerStatus byte

const (
	statusPrelude deserializerStatus = iota
	statusProtocolVersion
	statusMessageID
	statusPackageType
	statusPackageSizePrelude
	statusPackageSizeValue
	statusPayload
)

// Deserializer is responsible for keeping track of a basic FSM used to parse
// metadata about a received message, returning a DeserializationCandidate when
// all bytes are received.
type Deserializer struct {
	state     deserializerStatus
	msgMeta   *MessageMetadata
	tmpBuffer []byte
	lastError error
	readBytes uint64
}

func (d *Deserializer) reset() {
	d.state = statusPrelude
	d.msgMeta = nil
	d.lastError = nil
	d.readBytes = 0
}

// Feed provides a single byte to the FSM. If an invalid or unexpected value is
// received, the FSM is automatically reseted.
func (d *Deserializer) Feed(b byte) *DeserializationCandidate {
	switch d.state {
	case statusPrelude:
		expectedByte := magicBytes[d.readBytes]
		if b == expectedByte {
			d.readBytes++
			if d.readBytes == uint64(len(magicBytes)) {
				d.state = statusProtocolVersion
				d.msgMeta = &MessageMetadata{}
			}
		} else {
			d.reset()
		}
	case statusProtocolVersion:
		d.msgMeta.ProtocolVersion = b
		d.readBytes++
		d.state = statusMessageID
	case statusMessageID:
		d.msgMeta.MessageID = b
		d.readBytes++
		d.state = statusPackageType
	case statusPackageType:
		d.msgMeta.PackageType = b
		d.readBytes++
		d.state = statusPackageSizePrelude
	case statusPackageSizePrelude:
		if b > byte(lengthEncodingUint64) {
			d.reset()
		}
		switch lengthEncoding(b) {
		case lengthEncodingEmpty:
			defer d.reset()
			return d.candidate()
		case lengthEncodingUint8:
			d.tmpBuffer = make([]byte, 0, 1)
		case lengthEncodingUint16:
			d.tmpBuffer = make([]byte, 0, 2)
		case lengthEncodingUint32:
			d.tmpBuffer = make([]byte, 0, 4)
		case lengthEncodingUint64:
			d.tmpBuffer = make([]byte, 0, 8)
		}
		d.readBytes++
		d.state = statusPackageSizeValue
	case statusPackageSizeValue:
		d.tmpBuffer = append(d.tmpBuffer, b)
		d.readBytes++
		if len(d.tmpBuffer) == cap(d.tmpBuffer) {
			var newBuf []byte
			switch cap(d.tmpBuffer) {
			case 1:
				newBuf = make([]byte, 0, d.tmpBuffer[0])
			case 2:
				newBuf = make([]byte, 0, readUint16(d.tmpBuffer))
			case 4:
				newBuf = make([]byte, 0, readUint32(d.tmpBuffer))
			case 8:
				newBuf = make([]byte, 0, readUint64(d.tmpBuffer))
			}
			d.tmpBuffer = newBuf
			d.state = statusPayload
		}
	case statusPayload:
		d.tmpBuffer = append(d.tmpBuffer, b)
		d.readBytes++
		if len(d.tmpBuffer) == cap(d.tmpBuffer) {
			defer d.reset()
			return d.candidate()
		}
	}
	return nil
}

func (d *Deserializer) candidate() *DeserializationCandidate {
	return &DeserializationCandidate{
		buffer:      d.tmpBuffer,
		MessageMeta: d.msgMeta,
	}
}

// DeserializeNonMessage converts a previously serialised object (through
// SerializeNonMessage) to be converted back to a known type.
// Returns a pointer to the object with the provided type, or an error.
func DeserializeNonMessage(data []byte, into Serializable) (interface{}, error) {
	fields, err := deserialize(data)
	if err != nil {
		return nil, err
	}
	return createObjectFromType(into.LudwiegMeta(), reflect.TypeOf(into), fields[0].([]interface{})), nil
}
