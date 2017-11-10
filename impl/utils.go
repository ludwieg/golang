package impl

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
)

func writeUint16(v uint16, b *bytes.Buffer) {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, v)
	b.Write(buf)
}

func readUint16(b []byte) uint16 {
	return binary.LittleEndian.Uint16(b)
}

func writeUint32(v uint32, b *bytes.Buffer) {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, v)
	b.Write(buf)
}

func readUint32(b []byte) uint32 {
	return binary.LittleEndian.Uint32(b)
}

func writeUint64(v uint64, b *bytes.Buffer) {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	b.Write(buf)
}

func readUint64(b []byte) uint64 {
	return binary.LittleEndian.Uint64(b)
}

func writeFloat64(v float64, b *bytes.Buffer) {
	bits := math.Float64bits(v)
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, bits)
	b.Write(buf)
}

func readFloat64(b []byte) float64 {
	return math.Float64frombits(readUint64(b))
}

func writeSize(size uint64, b *bytes.Buffer) {
	switch {
	case size <= math.MaxUint8:
		b.WriteByte(byte(lengthEncodingUint8))
		b.WriteByte(byte(size))
	case size <= math.MaxUint16:
		b.WriteByte(byte(lengthEncodingUint16))
		writeUint16(uint16(size), b)
	case size <= math.MaxUint32:
		b.WriteByte(byte(lengthEncodingUint32))
		writeUint32(uint32(size), b)
	case size <= math.MaxUint64:
		b.WriteByte(byte(lengthEncodingUint64))
		writeUint64(uint64(size), b)
	default: // Sanity check
		panic(fmt.Errorf("invalid size value %d", size))
	}
}
