package impl

import (
	"reflect"
)

// LudwiegTypeAnnotation is used to annotate interface types based on indexed
// fields
type LudwiegTypeAnnotation struct {
	Type          ProtocolType
	ArrayType     ProtocolType
	ArraySize     string
	ArrayUserType reflect.Type
}

func (t LudwiegTypeAnnotation) metaProtocolByte() *metaProtocolByte {
	p := metaTypeFromByte(byte(t.Type))
	return &p
}

// ArrayOf creates a new array annotation with no predefined length
func ArrayOf(r interface{}) LudwiegTypeAnnotation {
	return ArrayOfWithSize(r, "*")
}

// ArrayOfWithSize creates a new array annotation using the provided size
// as the predefined size information.
func ArrayOfWithSize(r interface{}, size string) LudwiegTypeAnnotation {
	return LudwiegTypeAnnotation{
		Type:          TypeArray,
		ArrayType:     TypeStruct,
		ArraySize:     size,
		ArrayUserType: reflect.PtrTo(reflect.TypeOf(r)),
	}
}
