package impl

import (
	"bytes"
	"fmt"
	"strconv"
)

/* Ludwieg arrays cannot be nil, but they may be empty. Hence the lack of
LudwiegArray */

func serializeArray(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	// Here we may want to perform extra checks, such as if we have correct
	// types and that our size is sane.
	rawArraySize := c.annotation.ArraySize
	arrayType := c.annotation.ArrayType
	if arrayType == TypeUnknown {
		return fmt.Errorf("invalid array type %#v", arrayType)
	}
	if rawArraySize != "*" {
		_, err := strconv.Atoi(rawArraySize)
		if err != nil {
			return fmt.Errorf("invalid array size %s", rawArraySize)
		}
	}

	arrayLogicalSize := c.value.Len()
	array := c.value.Slice(0, arrayLogicalSize)
	// array will always be an slice of pointers.
	var arrBuf bytes.Buffer
	arrTypeAnnotation := &LudwiegTypeAnnotation{Type: arrayType}
	arrAnnotationByte := arrTypeAnnotation.metaProtocolByte()
	for i := 0; i < arrayLogicalSize; i++ {
		itemVal := array.Index(i)
		err := serialize(&arrBuf, &serializationCandidate{&itemVal, arrTypeAnnotation, arrAnnotationByte, false, false})
		if err != nil {
			return err
		}
	}

	writeSize(uint64(arrBuf.Len()), b)
	b.WriteByte(arrAnnotationByte.byte())
	writeSize(uint64(arrayLogicalSize), b)
	b.Write(arrBuf.Bytes())
	return nil
}

func decodeArray(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		// The main decoder, responsible for setting fields may either choose to
		// panic, or to leave the slice field with its initial value.
		return nil, nil
	}

	// Ludwieg arrays are weird monsters. Let's decode it:
	// 1. Size of the array buffer containing array data
	size := int(readSize(offset, b))

	// 2. Type being read
	arrayType := metaTypeFromByte(b[incr(offset)])

	// 3. the virtual size of the array represents how many items
	// the decoder is supposed to yield
	virtualSize := int(readSize(offset, b))

	result := make([]interface{}, 0, virtualSize)

	decoder, ok := registeredTypeDecoder[arrayType.ManagedType]

	tmpBuffer := b[*offset : *offset+size]

	incrSize(offset, size)

	if !ok {
		return nil, fmt.Errorf("unknown decoder for type %#v", arrayType)
	}

	innerOffset := 0

	for innerOffset < len(tmpBuffer) {
		i, err := decoder(arrayType, tmpBuffer, &innerOffset)
		if err != nil {
			return nil, err
		}
		result = append(result, i)
	}

	return result, nil
}
