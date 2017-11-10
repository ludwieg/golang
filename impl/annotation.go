package impl

// LudwiegTypeAnnotation is used to annotate interface types based on indexed
// fields
type LudwiegTypeAnnotation struct {
	Type      ProtocolType
	ArrayType ProtocolType
	ArraySize string
}

func (t LudwiegTypeAnnotation) metaProtocolByte() *metaProtocolByte {
	p := metaTypeFromByte(byte(t.Type))
	return &p
}
