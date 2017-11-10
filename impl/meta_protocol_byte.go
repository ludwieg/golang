package impl

type metaProtocolByte struct {
	Type           byte
	ManagedType    ProtocolType
	LengthPrefixed bool
	Empty          bool
	Known          bool
}

func (p *metaProtocolByte) read(b byte) {
	p.LengthPrefixed = (b & hasPrefixedLengthBit) == hasPrefixedLengthBit
	p.Empty = (b & isEmptyBit) == isEmptyBit
	p.Type = b &^ isEmptyBit
	if p.isKnown() {
		p.ManagedType = ProtocolType(p.Type)
	} else {
		p.ManagedType = TypeUnknown
	}
	p.Known = p.ManagedType != TypeUnknown
}

func (p *metaProtocolByte) isKnown() bool {
	for _, t := range knownTypes {
		if byte(t) == p.Type {
			return true
		}
	}
	return false
}

func (p *metaProtocolByte) byte() byte {
	result := p.Type
	if p.Empty {
		result |= isEmptyBit
	} else {
		result &^= isEmptyBit
	}

	if p.LengthPrefixed {
		result |= hasPrefixedLengthBit
	} else {
		result &^= hasPrefixedLengthBit
	}

	return result
}

func metaTypeFromByte(b byte) metaProtocolByte {
	var t metaProtocolByte
	t.read(b)
	return t
}
