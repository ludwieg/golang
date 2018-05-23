package impl

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// LudwiegUUID is used to safely represent a nullable UUID value
type LudwiegUUID LudwiegString

// UUID returns a safe nullable UUID value
func UUID(v string) *LudwiegUUID {
	uuid := LudwiegUUID(*String(v))
	return &uuid
}

func serializeUUID(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}

	rv := c.value.Interface()
	var val string
	switch v := rv.(type) {
	case *LudwiegString:
		val = v.Value
	case *LudwiegUUID:
		val = v.Value
	default:
		return illegalSetterValueError("uuid")
	}

	regexp := regexp.MustCompile(`[0-9a-f]+`)
	val = strings.ToLower(strings.Replace(val, "-", "", 0))
	if len(val) != 32 || !regexp.Match([]byte(val)) {
		return fmt.Errorf("invalid value %s for UUID field", val)
	}

	data := make([]byte, 16)
	dataCursor := 0
	for i := 0; i < 32; i += 2 {
		end := i + 2
		i, err := strconv.ParseUint(val[i:end], 16, 8)
		if err != nil {
			// Should not happen, since we're validating conformity and
			// length one step behind, but let's leave this here as a sanity
			// check.
			return fmt.Errorf("invalid value %s for UUID field", val)
		}
		data[dataCursor] = byte(i)
		dataCursor++
	}

	b.Write(data)

	return nil
}

func decodeUUID(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegUUID{}, nil
	}

	tmpBuf := b[*offset : *offset+16]
	incrSize(offset, 16)

	value := make([]string, 16)
	for i, b := range tmpBuf {
		value[i] = fmt.Sprintf("%x", b)
	}

	return &LudwiegUUID{
		HasValue: true,
		Value:    strings.Join(value, ""),
	}, nil
}
