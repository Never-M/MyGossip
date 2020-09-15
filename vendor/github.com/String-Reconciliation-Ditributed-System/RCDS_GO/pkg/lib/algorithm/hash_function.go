package algorithm

import (
	"crypto"
	"hash/fnv"
)

type hashData struct {
	bytes []byte
	err   error
}

func HashString(s string) *hashData {
	return &hashData{
		bytes: []byte(s),
		err:   nil,
	}
}

func HashBytesWithCryptoFunc(b []byte, hash crypto.Hash) *hashData {
	h := hash.New()
	_, herr := h.Write(b)
	return &hashData{
		bytes: h.Sum(nil),
		err:   herr,
	}

}

// TODO: Provide methods to select different hash functions or provide custom hash functions.
func (d *hashData) ToUint64() (uint64, error) {
	if d.err != nil {
		return 0, d.err
	}
	h := fnv.New64()
	_, err := h.Write(d.bytes)
	return h.Sum64(), err
}

func (d *hashData) ToUint32() (uint32, error) {
	if d.err != nil {
		return 0, d.err
	}
	h := fnv.New32()
	_, err := h.Write(d.bytes)
	return h.Sum32(), err
}

func (d *hashData) ToBytes() ([]byte, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d.bytes, nil
}
