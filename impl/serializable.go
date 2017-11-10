package impl

// Serializable is used internally by Ludwieg to safely coerce and validate
// serializable types.
type Serializable interface {
	LudwiegMeta() []LudwiegTypeAnnotation
}

// SerializablePackage is used internally by Ludwieg to coerce and validate
// top-level serializable structures also known as "Packages"
type SerializablePackage interface {
	LudwiegID() byte
	LudwiegMeta() []LudwiegTypeAnnotation
}
