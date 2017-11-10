package impl

import "fmt"

// LudwiegUnknown is just an utility type used to hold information about an
// unknown type. Currently, it only retains type metadata and a buffer
// containing the object value. Its decoder returns an error on unknown types
// that has no prefixed lenght.
type LudwiegUnknown struct {
	HasValue bool
	Value    []byte
}

func decodeUnknown(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		incr(offset)
		return &LudwiegUnknown{
			HasValue: false,
		}, nil
	}
	if !t.LengthPrefixed {
		return nil, fmt.Errorf("unknown type with no prefixed lenght cannot be decoded")
	}

	typeSize := readSize(offset, b)

	result := &LudwiegUnknown{
		HasValue: true,
		Value:    b[*offset:typeSize],
	}

	incrSize(offset, int(typeSize))

	return result, nil
}
