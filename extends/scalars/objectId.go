package scalars

import (
	"database/sql/driver"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"strconv"
)

type (
	ObjectID primitive.ObjectID
)

var NilObjectID = ObjectID(primitive.NilObjectID)

func NewObjectID() ObjectID {
	return ObjectID(primitive.NewObjectID())
}

func CreateNilObjectID() ObjectID {
	return ObjectID(primitive.NilObjectID)
}

func CreateObjectIdOrNilByHex(v string) ObjectID {
	id, err := primitive.ObjectIDFromHex(v)
	if err != nil {
		return ObjectID(primitive.NilObjectID)
	} else {
		return ObjectID(id)
	}
}

func CreateObjectIdByHex(v string) (res ObjectID, err error) {
	var id primitive.ObjectID
	if id, err = primitive.ObjectIDFromHex(v); err != nil {
		return NilObjectID, err
	}
	return ObjectID(id), nil
}

func (x *ObjectID) UnmarshalJSON(bytes []byte) (err error) {
	var id primitive.ObjectID
	var str string
	if str, err = strconv.Unquote(string(bytes)); err != nil {
		return
	}

	if id, err = primitive.ObjectIDFromHex(str); err != nil {
		return
	}
	*x = ObjectID(id)
	return nil
}

func (x ObjectID) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(x.String())), nil
}

func (x *ObjectID) Equal(v ObjectID) bool {
	return v.ObjectID().Hex() == x.ObjectID().Hex()
}

func (x *ObjectID) ObjectID() primitive.ObjectID {
	return primitive.ObjectID(*x)
}

func (x ObjectID) Value() (driver.Value, error) {
	return []byte(x.ObjectID().Hex()), nil
}

func (x *ObjectID) Scan(src any) (err error) {
	switch v := src.(type) {
	case string:
		var id primitive.ObjectID
		if id, err = primitive.ObjectIDFromHex(v); err != nil {
			return
		}
		*x = ObjectID(id)
		return
	case []byte:
		var id primitive.ObjectID
		if id, err = primitive.ObjectIDFromHex(string(v)); err != nil {
			return
		}
		*x = ObjectID(id)
		return
	default:
		return fmt.Errorf("InvalidObjectIdValue")
	}
}

func (x *ObjectID) String() string {
	return x.ObjectID().Hex()
}

func (x *ObjectID) IsNil() bool {
	return x.ObjectID().IsZero()
}

func MarshalObjectID(v ObjectID) graphql.Marshaler {
	return graphql.WriterFunc(func(writer io.Writer) {
		_, _ = writer.Write([]byte(strconv.Quote(v.ObjectID().Hex())))
	})
}

func UnmarshalObjectID(v interface{}) (res ObjectID, err error) {
	switch t := v.(type) {
	case string:
		var id primitive.ObjectID
		if id, err = primitive.ObjectIDFromHex(t); err != nil {
			return
		}
		return ObjectID(id), nil
	default:
		return ObjectID(primitive.NilObjectID), fmt.Errorf("InvalidObjectIdValue")
	}
}

type (
	ObjectIDList []ObjectID
)

func (x *ObjectIDList) ToString() (list []string) {
	for _, id := range *x {
		list = append(list, id.String())
	}
	return
}
