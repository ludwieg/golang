package impl

import "bytes"

// LudwiegString is used to safely represent a nullable string value
type LudwiegString struct {
	HasValue bool
	Value    string
}

// String returns a safe nullable String value
func String(v string) *LudwiegString {
	return &LudwiegString{
		HasValue: true,
		Value:    v,
	}
}

func serializeString(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegString); ok {
		writeSize(uint64(len(v.Value)), b)
		b.Write([]byte(v.Value))
	} else {
		return illegalSetterValueError("string")
	}
	return nil
}

func decodeString(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegString{}, nil
	}
	size := int(readSize(offset, b))
	tmpBuf := b[*offset : *offset+size]
	incrSize(offset, size)

	return &LudwiegString{
		HasValue: true,
		Value:    string(tmpBuf),
	}, nil
}
