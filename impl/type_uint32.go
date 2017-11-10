package impl

import "bytes"

// LudwiegUint32 is used to safely represent a nullable uint32 value
type LudwiegUint32 struct {
	HasValue bool
	Value    uint32
}

// Uint32 returns a safe nullable Uint32 value
func Uint32(v uint32) *LudwiegUint32 {
	return &LudwiegUint32{
		HasValue: true,
		Value:    v,
	}
}

func serializeUint32(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegUint32); ok {
		writeUint32(v.Value, b)
	} else {
		return illegalSetterValueError("uint32")
	}
	return nil
}

func decodeUint32(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegUint32{}, nil
	}
	tmpBuf := b[*offset : *offset+4]
	incrSize(offset, 4)

	return &LudwiegUint32{
		HasValue: true,
		Value:    readUint32(tmpBuf),
	}, nil
}
