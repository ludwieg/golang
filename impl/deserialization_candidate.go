package impl

import (
	"fmt"
	"reflect"
)

// DeserializationCandidate retains a buffer of raw information that may be
// parsed into fields and set on a registered structure
type DeserializationCandidate struct {
	MessageMeta *MessageMetadata
	buffer      []byte
}

// CanDeserialize may be used to check if a received package can be deserialized
// beforehand. This is useful when supporting several packages across different
// platforms.
func (c *DeserializationCandidate) CanDeserialize() bool {
	_, ok := registeredPackages[c.MessageMeta.PackageType]
	return ok
}

// Deserialize either return a registered interface{}, or an error. You may use
// the language's facilities to determine which type was returned, and act
// accordingly.
func (c *DeserializationCandidate) Deserialize() (interface{}, error) {
	if !c.CanDeserialize() {
		panic(fmt.Errorf("cannot deserialize unknown package %#v", c.MessageMeta.PackageType))
	}

	deserializationResult, err := deserialize(c.buffer)
	if err != nil {
		return nil, err
	}

	packageMeta := registeredPackages[c.MessageMeta.PackageType]
	packageType := packageMeta.nativeType

	return createObjectFromType(packageMeta.meta, packageType, deserializationResult), nil
}

func createObjectFromType(annotations []LudwiegTypeAnnotation, t reflect.Type, values []interface{}) interface{} {

	resultLen := len(values)

	instance := reflect.New(t)
	ptr := reflect.Indirect(instance)

	for i, fieldMeta := range annotations {
		if i > resultLen {
			break
		}
		rawValue := values[i]
		rawPointer := reflect.ValueOf(rawValue)

		field := t.Field(i)
		fieldValue := ptr.Field(i)

		switch fieldMeta.Type {
		case TypeUint8, TypeUint32, TypeUint64, TypeDouble, TypeString, TypeBool, TypeUUID, TypeAny:
			if rawPointer.IsNil() {
				continue
			}
			fieldValue.Set(rawPointer)
		case TypeStruct:
			// This will be slighlty trickier, yet feasible.
			// We will recurse the call to createObjectFromType, but
			// we need annotations and type information.
			metaValue, err := extractAnnotationsFromType(field.Type)
			if err != nil {
				panic(err)
			}
			val := createObjectFromType(metaValue, field.Type.Elem(), rawValue.([]interface{}))
			fieldValue.Set(reflect.ValueOf(val))
		case TypeArray:
			var curArr []interface{}
			// Empty arrays will cause rawValue to be nil. Here we can ensure
			// the array being accessed (and cast) is not nil.
			if rawValue == nil {
				curArr = []interface{}{}
			} else {
				curArr = rawValue.([]interface{})
			}

			// Custom-type arrays need extra attention here.
			if fieldMeta.ArrayType == TypeStruct {
				// At this point, we may have slices of slices, and each group
				// of slices must be used to instantiate a new object that will
				// be placed inside the slice.
				newArr := reflect.MakeSlice(reflect.SliceOf(fieldMeta.ArrayUserType), len(curArr), len(curArr))
				annotations, err := extractAnnotationsFromType(fieldMeta.ArrayUserType)
				if err != nil {
					panic(err)
				}
				for i, v := range curArr {
					// The next line uses Elem() on ArrayUserType, since it will
					// always be cast to a ptr using PtrTo.
					val := createObjectFromType(annotations, fieldMeta.ArrayUserType.Elem(), v.([]interface{}))
					ptr := newArr.Index(i)
					ptr.Set(reflect.ValueOf(val))
				}

				fieldValue.Set(newArr)
			} else {
				newArr := reflect.MakeSlice(field.Type, len(curArr), len(curArr))
				for i, v := range curArr {
					ptr := newArr.Index(i)
					ptr.Set(reflect.ValueOf(v))
				}

				fieldValue.Set(newArr)
			}
		case TypeBlob:
			fieldValue.Set(rawPointer)
		}

	}

	return instance.Interface()
}

func deserialize(buffer []byte) ([]interface{}, error) {
	offset := 0
	items := []interface{}{}

	for offset < len(buffer) {

		metaType := metaTypeFromByte(buffer[offset])
		incr(&offset)

		if decoder, ok := registeredTypeDecoder[metaType.ManagedType]; ok {
			if obj, err := decoder(metaType, buffer, &offset); err == nil {
				items = append(items, obj)
			} else {
				return items, err
			}
		} else {
			panic(fmt.Errorf("unknown decoder for type %#v", metaType))
		}
	}
	return items, nil
}

func incr(offset *int) int {
	return incrSize(offset, 1)
}

func incrSize(offset *int, size int) int {
	old := *offset
	*offset = *offset + size
	return old
}

func readSize(offset *int, buffer []byte) uint64 {

	b := buffer[incr(offset)]

	switch lengthEncoding(b) {
	case lengthEncodingEmpty:
		return uint64(0)
	case lengthEncodingUint8:
		buf := buffer[*offset]
		incr(offset)
		return uint64(buf)
	case lengthEncodingUint16:
		buf := buffer[*offset:2]
		incrSize(offset, 2)
		return uint64(readUint16(buf))
	case lengthEncodingUint32:
		buf := buffer[*offset:4]
		incrSize(offset, 4)
		return uint64(readUint32(buf))
	case lengthEncodingUint64:
		buf := buffer[*offset:8]
		incrSize(offset, 8)
		return readUint64(buf)
	default:
		panic(fmt.Errorf("unknown size prefix %#v", b))
	}
}

func extractAnnotationsFromType(t reflect.Type) ([]LudwiegTypeAnnotation, error) {
	serializable := reflect.TypeOf((*Serializable)(nil)).Elem()
	if !t.Implements(serializable) {
		return nil, fmt.Errorf("type %s does not implement Serializable", t.String())
	}

	if t.Kind() != reflect.Struct {
		t = t.Elem()
	}

	structIndirect := reflect.Indirect(reflect.New(t))
	return structIndirect.
		MethodByName("LudwiegMeta").
		Call([]reflect.Value{})[0].
		Interface().([]LudwiegTypeAnnotation), nil
}
