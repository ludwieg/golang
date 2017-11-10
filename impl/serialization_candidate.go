package impl

import (
	"reflect"
)

// serializableCandidate is used internally by the serialiser to retain metadata
// about a value being serialised.
type serializationCandidate struct {
	// value retains the reflected value of the entity being serialised
	value *reflect.Value

	// annotation retains information recovered through the LudwiegMeta() method
	// defined by the structure holding the value being serialised
	annotation *LudwiegTypeAnnotation

	// meta retains information about the type being serialised
	meta *metaProtocolByte

	// writeType indicates whether the serialised must write information about
	// the type when writing this entity to the buffer. Operations such as
	// array serialisation set this field as false, in order to avoid overhead
	writeType bool

	// isRoot is only used by the struct serialiser, since it is shared among
	// other generic serializers.
	isRoot bool
}
