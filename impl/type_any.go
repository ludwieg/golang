package impl

import (
	"bytes"
	"fmt"
	"reflect"
)

// LudwiegAny is used to safely represent a nullable "any" value.
type LudwiegAny struct {
	HasValue bool
	Value    interface{}
}

// Any returns a safe nullable Any value
func Any(v interface{}) *LudwiegAny {
	return &LudwiegAny{
		HasValue: true,
		Value:    v,
	}
}

func serializeAny(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}

	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegAny); ok {
		var internalBuffer bytes.Buffer
		reflectValue := reflect.ValueOf(v.Value)

		// Here we will need some helpers. Bear with me:
		serializeArray := func(t ProtocolType) error {
			if reflectValue.Len() == 0 {
				return fmt.Errorf("type Any cannot serialize empty array")
			}
			meta := metaTypeFromByte(byte(TypeArray))
			return serializeArray(&serializationCandidate{
				annotation: &LudwiegTypeAnnotation{ArrayType: t, ArraySize: "*", Type: TypeArray},
				meta:       &meta,
				value:      &reflectValue,
				writeType:  true,
			}, &internalBuffer)
		}

		serializeSimple := func(f serializerFunc, t ProtocolType) error {
			if t == TypeStruct {
				return fmt.Errorf("type Any cannot retain an struct")
			}
			meta := metaTypeFromByte(byte(t))
			return f(&serializationCandidate{
				annotation: &LudwiegTypeAnnotation{Type: t},
				meta:       &meta,
				value:      &reflectValue,
				writeType:  true,
			}, &internalBuffer)
		}

		var err error
		switch v.Value.(type) {
		case *LudwiegUint8:
			err = serializeSimple(serializeUint8, TypeUint8)
		case *LudwiegUint32:
			err = serializeSimple(serializeUint32, TypeUint32)
		case *LudwiegUint64:
			err = serializeSimple(serializeUint64, TypeUint64)
		case *LudwiegDouble:
			err = serializeSimple(serializeDouble, TypeDouble)
		case *LudwiegString:
			err = serializeSimple(serializeString, TypeString)
		case *LudwiegUUID:
			err = serializeSimple(serializeUUID, TypeUUID)
		case []byte:
			err = serializeSimple(serializeBlob, TypeBlob)
		case [][]byte:
			err = serializeArray(TypeBlob)
		case []*LudwiegUint8:
			err = serializeArray(TypeUint8)
		case []*LudwiegUint32:
			err = serializeArray(TypeUint32)
		case []*LudwiegUint64:
			err = serializeArray(TypeUint64)
		case []*LudwiegDouble:
			err = serializeArray(TypeDouble)
		case []*LudwiegString:
			err = serializeArray(TypeString)
		case []*LudwiegUUID:
			err = serializeArray(TypeUUID)
		default:
			// Oh my.
			// Probs a struct. Assume a pointer, panic otherwise.
			if reflectValue.Kind() != reflect.Ptr {
				return fmt.Errorf("cannot serialize non-pointer 'Any' value")
			}

			err = serializeSimple(serializeStruct, TypeStruct)
		}
		if err != nil {
			return err
		}
		writeSize(uint64(internalBuffer.Len()), b)
		b.Write(internalBuffer.Bytes())
	} else {
		return illegalSetterValueError("any")
	}

	return nil
}

func decodeAny(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegAny{}, nil
	}

	size := int(readSize(offset, b))
	tmpBuf := b[*offset : *offset+size]
	incrSize(offset, size)

	val, err := deserialize(tmpBuf)
	if err != nil {
		return nil, err
	}

	return &LudwiegAny{
		HasValue: true,
		Value:    val[0],
	}, nil
}
