package main

import (
	"testing"

	"github.com/stretchr/testify/assert"

	impl "github.com/ludwieg/golang/impl"
)

func init() {
	impl.RegisterPackage(Test{})
}

type Test struct {
	FieldA *impl.LudwiegUint8
	FieldB *impl.LudwiegUint32
	FieldC *impl.LudwiegUint64
	FieldD *impl.LudwiegFloat64
	FieldE *impl.LudwiegString
	FieldF []byte
	FieldG *impl.LudwiegBool
	FieldH *impl.LudwiegUUID
	FieldY *impl.LudwiegAny
	FieldZ [](*impl.LudwiegString)
	FieldI *TestSub
}

func (t Test) LudwiegID() byte {
	return 0x01
}

func (t Test) LudwiegMeta() []impl.LudwiegTypeAnnotation {
	return []impl.LudwiegTypeAnnotation{
		{Type: impl.TypeUint8},
		{Type: impl.TypeUint32},
		{Type: impl.TypeUint64},
		{Type: impl.TypeFloat64},
		{Type: impl.TypeString},
		{Type: impl.TypeBlob},
		{Type: impl.TypeBool},
		{Type: impl.TypeUUID},
		{Type: impl.TypeAny},
		{Type: impl.TypeArray, ArraySize: "*", ArrayType: impl.TypeString},
		{Type: impl.TypeStruct},
	}
}

type TestSub struct {
	FieldJ *impl.LudwiegString
	FieldK *TestSubOther
}

func (t TestSub) LudwiegMeta() []impl.LudwiegTypeAnnotation {
	return []impl.LudwiegTypeAnnotation{
		{Type: impl.TypeString},
		{Type: impl.TypeStruct},
	}
}

type TestSubOther struct {
	FieldL *impl.LudwiegString
}

func (t TestSubOther) LudwiegMeta() []impl.LudwiegTypeAnnotation {
	return []impl.LudwiegTypeAnnotation{
		{Type: impl.TypeString},
	}
}

func TestEncoderDecoder(t *testing.T) {
	obj := Test{
		FieldA: impl.Uint8(27),
		FieldB: impl.Uint32(28),
		FieldC: impl.Uint64(29),
		FieldD: impl.Float64(30.2),
		FieldE: impl.String("String"),
		FieldF: []byte{0x27, 0x24, 0x50},
		FieldG: impl.Bool(true),
		FieldH: impl.UUID("3232ee42c2f24baf841318335b4d5640"),
		FieldY: impl.Any(impl.String("Any field retaining a string")),
		FieldZ: []*impl.LudwiegString{impl.String("Robin"), impl.String("Tom")},
		FieldI: &TestSub{
			FieldJ: impl.String("Structure"),
			FieldK: &TestSubOther{
				FieldL: impl.String("Other Structure"),
			},
		},
	}

	buf, err := impl.Serialize(obj, 0x01)
	if err != nil {
		t.Error(err)
	}

	d := impl.Deserializer{}
	for _, b := range buf.Bytes() {
		if r := d.Feed(b); r != nil {
			v, err := r.Deserialize()
			if err != nil {
				t.Error(err)
			}
			r := v.(*Test)
			assert.True(t, r.FieldA.HasValue)
			assert.Equal(t, uint8(27), r.FieldA.Value)
			assert.Equal(t, uint32(28), r.FieldB.Value)
			assert.Equal(t, uint64(29), r.FieldC.Value)
			assert.Equal(t, 30.2, r.FieldD.Value)
			assert.Equal(t, "String", r.FieldE.Value)
			assert.Equal(t, uint8(0x27), r.FieldF[0])
			assert.Equal(t, uint8(0x24), r.FieldF[1])
			assert.Equal(t, uint8(0x50), r.FieldF[2])
			assert.True(t, r.FieldG.Value)
			assert.Equal(t, "3232ee42c2f24baf841318335b4d5640", r.FieldH.Value)
			av, ok := r.FieldY.Value.(*impl.LudwiegString)
			assert.True(t, ok)
			assert.Equal(t, "Any field retaining a string", av.Value)
			assert.Equal(t, "Robin", r.FieldZ[0].Value)
			assert.Equal(t, "Tom", r.FieldZ[1].Value)
			assert.NotNil(t, r.FieldI)
			assert.Equal(t, "Structure", r.FieldI.FieldJ.Value)
			assert.NotNil(t, r.FieldI.FieldK)
			assert.Equal(t, "Other Structure", r.FieldI.FieldK.FieldL.Value)
			return
		}
	}

	t.Errorf("Decoding failed.")
}
