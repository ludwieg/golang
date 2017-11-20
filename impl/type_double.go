package impl

import (
	"bytes"
)

// LudwiegDouble is used to safely represent a nullable float64 value
type LudwiegDouble struct {
	HasValue bool
	Value    float64
}

// Double returns a safe nullable Double value
func Double(v float64) *LudwiegDouble {
	return &LudwiegDouble{
		HasValue: true,
		Value:    v,
	}
}

func serializeDouble(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegDouble); ok {
		writeDouble(v.Value, b)
	} else {
		return illegalSetterValueError("float64")
	}
	return nil
}

func decodeDouble(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegDouble{}, nil
	}
	tmpBuf := b[*offset : *offset+8]
	incrSize(offset, 8)

	return &LudwiegDouble{
		HasValue: true,
		Value:    readDouble(tmpBuf),
	}, nil
}
