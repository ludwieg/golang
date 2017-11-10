package impl

import "bytes"

/* In Go, blob is represented by the native byte slice, which is also nullable.
Hence the lack of LudwiegBlob and such */

func serializeBlob(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}

	rv := c.value.Slice(0, c.value.Len()).Interface()
	if v, ok := rv.([]byte); ok {
		writeSize(uint64(len(v)), b)
		b.Write(v)
	} else {
		return illegalSetterValueError("string")
	}
	return nil
}

func decodeBlob(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return []byte{}, nil
	}

	size := int(readSize(offset, b))
	tmpBuf := b[*offset : *offset+size]
	incrSize(offset, size)

	return tmpBuf, nil
}
