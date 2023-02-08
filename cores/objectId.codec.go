package cores

import (
	"fmt"
	"github.com/bamtols/mong-go/extends/scalars"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type (
	ObjectIDCodec struct{}
)

// interface 적용된걸 알려주는 힌트
var _ bsoncodec.ValueCodec = &ObjectIDCodec{}

func ObjectIDCodecRegister() (reflect.Type, bsoncodec.ValueCodec) {
	return reflect.TypeOf(scalars.ObjectID{}), &ObjectIDCodec{}
}

func (x *ObjectIDCodec) EncodeValue(ctx bsoncodec.EncodeContext, w bsonrw.ValueWriter, v reflect.Value) (err error) {
	dec, ok := v.Interface().(scalars.ObjectID)
	if !ok {
		return fmt.Errorf("InvalidObjectIDValue")
	}
	return w.WriteObjectID(primitive.ObjectID(dec))
}

func (x *ObjectIDCodec) DecodeValue(ctx bsoncodec.DecodeContext, r bsonrw.ValueReader, v reflect.Value) (err error) {
	obj, err := r.ReadObjectID()
	if err != nil {
		return fmt.Errorf("InvalidObjectIDValue")
	}
	v.Set(reflect.ValueOf(scalars.ObjectID(obj)))
	return
}
