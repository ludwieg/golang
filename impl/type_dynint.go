package impl

import (
	"bytes"
	"math"
)

// This defines the constant margin of error when transforming possible float
// values into another kind.
const dynintEpsilon = 1e-9

// LudwiegDynInt is used to represent an integer field with varying size,
// automatically setting its precision on-the-fly.
type LudwiegDynInt struct {
	hasValue       bool
	value          float64
	UnderlyingType DynIntValueKind
}

// DynInt returns a safe nullable DynInt value
func DynInt(v interface{}) *LudwiegDynInt {
	retVal := &LudwiegDynInt{
		hasValue: true,
	}
	var kind DynIntValueKind
	var value float64
	switch val := v.(type) {
	case int:
		kind, value = dynintInferNumberType(float64(val))
	case int8:
		kind, value = dynintInferNumberType(float64(val))
	case int16:
		kind, value = dynintInferNumberType(float64(val))
	case int32:
		kind, value = dynintInferNumberType(float64(val))
	case int64:
		kind, value = dynintInferNumberType(float64(val))
	case uint:
		kind, value = dynintInferNumberType(float64(val))
	case uint8:
		kind, value = dynintInferNumberType(float64(val))
	case uint16:
		kind, value = dynintInferNumberType(float64(val))
	case uint32:
		kind, value = dynintInferNumberType(float64(val))
	case uint64:
		kind, value = dynintInferNumberType(float64(val))
	case float32:
		kind, value = dynintInferNumberType(float64(val))
	case float64:
		kind, value = dynintInferNumberType(float64(val))
	case nil:
		kind = DynIntValueKindInt8
		value = 0
		retVal.hasValue = false
	default:
		panic(illegalSetterValueError("dynint"))
	}

	retVal.value = value
	retVal.UnderlyingType = kind

	return retVal
}

func serializeDynInt(c *serializationCandidate, b *bytes.Buffer) error {
	if c.writeType {
		b.WriteByte(c.meta.byte())
	}
	vv := c.value.Interface()
	if v, ok := vv.(*LudwiegDynInt); ok {
		b.WriteByte(byte(v.UnderlyingType))

		switch v.UnderlyingType {
		case DynIntValueKindInvalid:
		case DynIntValueKindUint8, DynIntValueKindInt8:
			b.WriteByte(byte(v.value))
		case DynIntValueKindInt16, DynIntValueKindUint16:
			writeUint16(uint16(v.value), b)
		case DynIntValueKindUint32, DynIntValueKindInt32:
			writeUint32(uint32(v.value), b)
		case DynIntValueKindUint64, DynIntValueKindInt64:
			writeUint64(uint64(v.value), b)
		case DynIntValueKindFloat64, DynIntValueKindFloat32:
			writeDouble(v.value, b)
		}
	} else {
		return illegalSetterValueError("dynint")
	}
	return nil
}

func dynintInferNumberType(val float64) (DynIntValueKind, float64) {
	// First, we will want to give special attention to float kinds, since
	// we're about to drop all the fraction part of its value.
	integer, frac := math.Modf(val)
	if frac > dynintEpsilon && frac < 1.0-dynintEpsilon {
		// Float.
		if val >= -math.MaxFloat32 && val <= math.MaxFloat32 {
			return DynIntValueKindFloat32, val
		}
		return DynIntValueKindFloat32, val
	}

	maxMins := [][]float64{
		{0, math.MaxUint8},
		{0, math.MaxUint16},
		{0, math.MaxUint32},
		{0, math.MaxUint64},
		{math.MinInt8, math.MaxInt8},
		{math.MinInt16, math.MaxInt16},
		{math.MinInt32, math.MaxInt32},
		{math.MinInt64, math.MaxInt64},
	}

	types := []DynIntValueKind{
		DynIntValueKindUint8,
		DynIntValueKindUint16,
		DynIntValueKindUint32,
		DynIntValueKindUint64,
		DynIntValueKindInt8,
		DynIntValueKindInt16,
		DynIntValueKindInt32,
		DynIntValueKindInt64,
	}

	for i, v := range maxMins {
		if integer >= v[0] && integer <= v[1] {
			return types[i], integer
		}
	}

	return DynIntValueKindInvalid, 0
}

func decodeDynInt(t metaProtocolByte, b []byte, offset *int) (interface{}, error) {
	if t.Empty {
		return &LudwiegDynInt{}, nil
	}
	dynKind := b[incrSize(offset, 1)]
	var value float64
	var tmpBuf []byte

	switch DynIntValueKind(dynKind) {
	case DynIntValueKindUint8, DynIntValueKindInt8:
		value = float64(b[incrSize(offset, 1)])
	case DynIntValueKindInt16, DynIntValueKindUint16:
		tmpBuf = b[*offset : *offset+2]
		incrSize(offset, 2)
		value = float64(readUint16(tmpBuf))
	case DynIntValueKindUint32, DynIntValueKindInt32:
		tmpBuf = b[*offset : *offset+4]
		incrSize(offset, 4)
		value = float64(readUint32(tmpBuf))
	case DynIntValueKindUint64, DynIntValueKindInt64:
		tmpBuf = b[*offset : *offset+8]
		incrSize(offset, 8)
		value = float64(readUint64(tmpBuf))
	case DynIntValueKindFloat64, DynIntValueKindFloat32:
		tmpBuf = b[*offset : *offset+8]
		incrSize(offset, 8)
		value = readDouble(tmpBuf)
	default:
		return &LudwiegDynInt{}, nil
	}

	result := LudwiegDynInt{
		hasValue:       true,
		value:          value,
		UnderlyingType: DynIntValueKind(dynKind),
	}

	return &result, nil
}

// Next we need to coerce back to known types

// Uint8 returns the internal representation of this type as an uint8
func (d *LudwiegDynInt) Uint8() uint8 {
	if !d.hasValue {
		return 0
	}
	return uint8(d.value)
}

// Uint16 returns the internal representation of this type as an uint16
func (d *LudwiegDynInt) Uint16() uint16 {
	if !d.hasValue {
		return 0
	}
	return uint16(d.value)
}

// Uint32 returns the internal representation of this type as an uint32
func (d *LudwiegDynInt) Uint32() uint32 {
	if !d.hasValue {
		return 0
	}
	return uint32(d.value)
}

// Uint64 returns the internal representation of this type as an uint64
func (d *LudwiegDynInt) Uint64() uint64 {
	if !d.hasValue {
		return 0
	}
	return uint64(d.value)
}

// Int8 returns the internal representation of this type as an int8
func (d *LudwiegDynInt) Int8() int8 {
	if !d.hasValue {
		return 0
	}
	return int8(d.value)
}

// Int16 returns the internal representation of this type as an int16
func (d *LudwiegDynInt) Int16() int16 {
	if !d.hasValue {
		return 0
	}
	return int16(d.value)
}

// Int32 returns the internal representation of this type as an int32
func (d *LudwiegDynInt) Int32() int32 {
	if !d.hasValue {
		return 0
	}
	return int32(d.value)
}

// Int64 returns the internal representation of this type as an int64
func (d *LudwiegDynInt) Int64() int64 {
	if !d.hasValue {
		return 0
	}
	return int64(d.value)
}

// Float32 returns the internal representation of this type as a float32
func (d *LudwiegDynInt) Float32() float32 {
	if !d.hasValue {
		return 0
	}
	return float32(d.value)
}

// Float64 returns the internal representation of this type as a float64
func (d *LudwiegDynInt) Float64() float64 {
	if !d.hasValue {
		return 0
	}
	return float64(d.value)
}

// Int returns the internal representation of this type as a int
func (d *LudwiegDynInt) Int() int {
	if !d.hasValue {
		return 0
	}
	return int(d.value)
}
