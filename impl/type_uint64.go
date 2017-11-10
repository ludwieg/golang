package impl

import "bytes"

// LudwiegUint64 is used to safely represent a nullable uint64 value
type LudwiegUint64 struct {
	HasValue bool
	Value    uint64
}

// Uint64 returns a safe nullable Uint64 value
func Uint64(v uint64) *LudwiegUint64 {
	return &LudwiegUint64{
		HasValue: true,
		Value:    v,
	}
}

func serializeUint64(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegUint64); ok {
		writeUint64(v.Value, b)
	} else {
		return illegalSetterValueError("uint64")
	}
	return nil
}

func decodeUint64(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegUint64{}, nil
	}
	tmpBuf := b[*offset : *offset+8]
	incrSize(offset, 8)

	return &LudwiegUint64{
		HasValue: true,
		Value:    readUint64(tmpBuf),
	}, nil
}
