package cores

import (
	"fmt"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/bsoncodec"
	"go.mongodb.org/mongo-driver/bson/bsonrw"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type DecimalCodec struct{}

var _ bsoncodec.ValueCodec = &DecimalCodec{}

func DecimalCodecRegister() (reflect.Type, bsoncodec.ValueCodec) {
	return reflect.TypeOf(decimal.Decimal{}), &DecimalCodec{}
}

func (dc *DecimalCodec) EncodeValue(ectx bsoncodec.EncodeContext, vw bsonrw.ValueWriter, val reflect.Value) error {
	// Use reflection to convert the reflect.Value to decimal.Decimal.
	dec, ok := val.Interface().(decimal.Decimal)
	if !ok {
		return fmt.Errorf("InvalidDecimalValue")
	}

	// Convert decimal.Decimal to primitive.Decimal128.
	primDec, err := primitive.ParseDecimal128(dec.String())
	if err != nil {
		return fmt.Errorf("error converting decimal.Decimal %v to primitive.Decimal128: %v", dec, err)
	}
	return vw.WriteDecimal128(primDec)
}

func (dc *DecimalCodec) DecodeValue(ectx bsoncodec.DecodeContext, vr bsonrw.ValueReader, val reflect.Value) error {
	// Read primitive.Decimal128 from the ValueReader.
	primDec, err := vr.ReadDecimal128()
	if err != nil {
		return fmt.Errorf("error reading primitive.Decimal128 from ValueReader: %v", err)
	}

	// Convert primitive.Decimal128 to decimal.Decimal.
	dec, err := decimal.NewFromString(primDec.String())
	if err != nil {
		return fmt.Errorf("InvalidDecimalValue")
	}

	// Set val to the decimal.Decimal value contained in dec.
	val.Set(reflect.ValueOf(dec))
	return nil
}
