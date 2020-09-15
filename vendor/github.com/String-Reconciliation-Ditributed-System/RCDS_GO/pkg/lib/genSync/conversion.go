package genSync

import (
	"log"
	"math/big"
	"reflect"
)

type bigint big.Int

func ToBigInt(input interface{}) *bigint {
	zz := new(big.Int)
	switch input.(type) {
	case string:
		b := []byte(input.(string))
		zz.SetBytes(b)
	case uint64:
		zz.SetUint64(input.(uint64))
	case []byte:
		zz.SetBytes(input.([]byte))
	default:
		log.Panicf("input %v is not supported for converting to 'Big/Int' as type %s", input, reflect.TypeOf(input).Name())
	}
	return (*bigint)(zz)
}

func (b *bigint) ToString() string {
	i := (big.Int)(*b)
	return string(i.Bytes())
}

func (b *bigint) ToUint64() uint64 {
	i := (big.Int)(*b)
	return i.Uint64()
}

func (b *bigint) ToBytes() []byte {
	i := (big.Int)(*b)
	return i.Bytes()
}
