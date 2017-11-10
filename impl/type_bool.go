package impl

import "bytes"

// LudwiegBool is used to safely represent a nullable bool value
type LudwiegBool struct {
	HasValue bool
	Value    bool
}

// Bool returns a safe nullable Bool value
func Bool(v bool) *LudwiegBool {
	return &LudwiegBool{
		HasValue: true,
		Value:    v,
	}
}

func serializeBool(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}

	rv := c.value.Interface()
	if v, ok := rv.(*LudwiegBool); ok {
		var by byte
		if v.Value {
			by = 0x1
		}
		b.WriteByte(by)
	} else {
		return illegalSetterValueError("bool")
	}
	return nil
}

func decodeBool(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegBool{}, nil
	}
	return &LudwiegBool{
		HasValue: true,
		Value:    b[incrSize(offset, 1)] == 0x1,
	}, nil
}
