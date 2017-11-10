package impl

import (
	"bytes"
)

// LudwiegFloat64 is used to safely represent a nullable float64 value
type LudwiegFloat64 struct {
	HasValue bool
	Value    float64
}

// Float64 returns a safe nullable Float64 value
func Float64(v float64) *LudwiegFloat64 {
	return &LudwiegFloat64{
		HasValue: true,
		Value:    v,
	}
}

func serializeFloat64(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegFloat64); ok {
		writeFloat64(v.Value, b)
	} else {
		return illegalSetterValueError("float64")
	}
	return nil
}

func decodeFloat64(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegFloat64{}, nil
	}
	tmpBuf := b[*offset : *offset+8]
	incrSize(offset, 8)

	return &LudwiegFloat64{
		HasValue: true,
		Value:    readFloat64(tmpBuf),
	}, nil
}
