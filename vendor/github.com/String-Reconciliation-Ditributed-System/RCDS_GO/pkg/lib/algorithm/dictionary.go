package algorithm

import (
	"fmt"
)

// Dictionary records the mapping between hash value and string.
type Dictionary map[uint64]string

// This is a local Dictionary to store string and hash transition.
var localDictionary = make(Dictionary)

// AddToDict converts a string in a hash value and add this pair of string and hash to the local Dictionary.
// It returns the hash value of the string and errors out if there exist hash collision or hash convection error.
func (d *Dictionary) AddToDict(entry string) (uint64, error) {
	if entry == "" {
		return 0, fmt.Errorf("no empty string should be added to the Dictionary")
	}
	hash, err := HashString(entry).ToUint64()
	if err != nil {
		return 0, fmt.Errorf("failed to convert string '%s' to hash value, %v", entry, err)
	}

	if val, isExist := (*d)[hash]; isExist && val != entry {
		return 0, fmt.Errorf("hash collision for string '%s' and '%s', %v", val, entry, err)
	} else if !isExist {
		(*d)[hash] = entry
	}
	return hash, nil
}

// LookupDict returns string that maps to the hash value.
// The function returns error if hash value does not exist in the Dictionary or the mapped string is empty.
func (d *Dictionary) LookupDict(hash uint64) (string, error) {
	val, isExist := (*d)[hash]
	if !isExist {
		return "", fmt.Errorf("hash value %d does not exist in the local Dictionary", hash)
	}
	if val == "" {
		return "", fmt.Errorf("error looking up hash value %d, the string value is empty", hash)
	}
	return val, nil
}
