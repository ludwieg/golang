package impl

import (
	"bytes"
	"fmt"
	"reflect"
)

func serializeStruct(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}

	reflectValue := c.value

	// Here we need to forcefully coerce a ptr into its direct value,
	// since this is used by common serializers, and the main serializer
	// entrypoint.
	if reflectValue.Kind() == reflect.Ptr {
		nv := reflect.Indirect(*reflectValue)
		reflectValue = &nv
	}

	// At this point, our value must be serializable. Otherwise, panic.
	serializableType := reflect.TypeOf((*Serializable)(nil)).Elem()
	if !reflectValue.Type().Implements(serializableType) {
		return fmt.Errorf("illegal attempt to serialize a non-serializable type %#v", reflectValue)
	}

	// Here we can safely convert it to a Serializable and extract metadata
	coercedSerializable := reflectValue.Convert(serializableType).Interface().(Serializable)
	typeMeta := coercedSerializable.LudwiegMeta()

	var internalBuffer bytes.Buffer

	baseObjType := reflectValue.Type()
	for i := 0; i < baseObjType.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldAnnotation := typeMeta[i]

		err := serialize(&internalBuffer, &serializationCandidate{&fieldValue, &fieldAnnotation, fieldAnnotation.metaProtocolByte(), true, false})
		if err != nil {
			return err
		}
	}

	if !c.isRoot {
		writeSize(uint64(internalBuffer.Len()), b)
	}
	b.Write(internalBuffer.Bytes())
	return nil
}

func decodeStruct(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		// Empty structs are handled as nil pointers
		return nil, nil
	}

	strSize := int(readSize(offset, b))
	tmpBuffer := b[*offset : *offset+strSize]
	incrSize(offset, strSize)

	// Also, struct values are returned as an array of possible fields and other
	// structs. The main decoder is responsible for converting to known types
	// and such.
	return deserialize(tmpBuffer)
}
