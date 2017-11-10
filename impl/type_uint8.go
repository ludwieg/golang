package impl

import (
	"bytes"
)

// LudwiegUint8 is used to safely represent a nullable uint8 value
type LudwiegUint8 struct {
	HasValue bool
	Value    uint8
}

// Uint8 returns a safe nullable uint8 value
func Uint8(v uint8) *LudwiegUint8 {
	return &LudwiegUint8{
		HasValue: true,
		Value:    v,
	}
}

func serializeUint8(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegUint8); ok {
		b.WriteByte(byte(v.Value))
	} else {
		return illegalSetterValueError("uint8")
	}
	return nil
}

func decodeUint8(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegUint8{}, nil
	}
	return &LudwiegUint8{
		HasValue: true,
		Value:    b[incrSize(offset, 1)],
	}, nil
}
