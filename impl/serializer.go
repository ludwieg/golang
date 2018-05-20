package impl

import (
	"bytes"
	"fmt"
	"reflect"
)

func illegalSetterValueError(t string) error {
	return fmt.Errorf("illegal attempt to set invalid type into protocol structure of type %s", t)
}

// SerializeMessage serializes a given SerializablePackage and messageID into a
// transferable byte buffer
func SerializeMessage(p SerializablePackage, messageID byte) (bufPtr *bytes.Buffer, err error) {
	var buf bytes.Buffer

	// 1. Write Magic Bytes
	buf.Write(magicBytes)

	// 2. Write Message Meta
	meta := MessageMetadata{
		ProtocolVersion: 0x01, // Version 1
		MessageID:       messageID,
		PackageType:     p.LudwiegID(),
	}
	meta.writeTo(&buf)

	value := reflect.ValueOf(p)

	var tmpBuf bytes.Buffer

	// 3. Prepare Fields
	err = serializeStruct(&serializationCandidate{
		isRoot:    true,
		writeType: false,
		value:     &value,
	}, &tmpBuf)

	// 4. Write Package Size
	writeSize(uint64(tmpBuf.Len()), &buf)

	// 5. Write Payload
	buf.Write(tmpBuf.Bytes())

	return &buf, err
}

// SerializeNonMessage allows a non-message object (preferable a Struct) to be
// serialised without the overhead of a message header. This also does not offer
// any security that messages do, such as protocol versions and such. Use at your
// own risk.
func SerializeNonMessage(p Serializable) (bufPtr *bytes.Buffer, err error) {
	tmpBuf := bytes.Buffer{}
	value := reflect.ValueOf(p)
	err = serializeStruct(&serializationCandidate{
		isRoot:    false,
		writeType: true,
		value:     &value,
		meta: &metaProtocolByte{
			Empty:          false,
			Known:          true,
			LengthPrefixed: true,
			Type:           TypeStruct,
		},
	}, &tmpBuf)

	return &tmpBuf, err
}

type serializerFunc func(c *serializationCandidate, b *bytes.Buffer) error

func serialize(buf *bytes.Buffer, c *serializationCandidate) error {
	value := c.value
	annotation := c.annotation

	if value.IsNil() {
		if c.writeType {
			// Empty values are handled by just writing the protocol type to the
			// stream with the IsEmpty bit set.
			c.meta.Empty = true
			buf.WriteByte(c.meta.byte())
		}
		return nil
	}

	var serializer serializerFunc

	switch annotation.Type {
	case TypeUnknown:
		return fmt.Errorf("cannot serialize unknown type")
	case TypeUint8:
		serializer = serializeUint8
	case TypeUint32:
		serializer = serializeUint32
	case TypeUint64:
		serializer = serializeUint64
	case TypeDouble:
		serializer = serializeDouble
	case TypeString:
		serializer = serializeString
	case TypeBlob:
		serializer = serializeBlob
	case TypeBool:
		serializer = serializeBool
	case TypeArray:
		serializer = serializeArray
	case TypeUUID:
		serializer = serializeUUID
	case TypeStruct:
		serializer = serializeStruct
	case TypeAny:
		serializer = serializeAny
	case TypeDynInt:
		serializer = serializeDynInt
	}

	if serializer == nil {
		return fmt.Errorf("cannot determine serializer for type %#v", annotation.Type)
	}

	return serializer(c, buf)
}
