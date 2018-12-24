package hashgen

import (
	"encoding/binary"
	"fmt"
	"testing"

	"golang.org/x/crypto/sha3"
)

func TestGenerateHash(t *testing.T) {
	tests := []struct {
		name string
		keys []uint64
	}{
		{
			name: "1 key",
			keys: []uint64{123456},
		},
		{
			name: "3 keys",
			keys: []uint64{
				123456,
				126543,
				120000,
			},
		},
		{
			name: "5 keys",
			keys: []uint64{
				123456,
				126543,
				123333,
				125555,
				120000,
			},
		},
		{
			name: "same keys",
			keys: []uint64{
				123456,
				123456,
				123456,
				123456,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sumChan := generateHash(tt.keys)
			for idx := 0; idx < len(tt.keys); idx++ {
				act := <-sumChan
				expKey := make([]byte, 8)
				binary.LittleEndian.PutUint64(expKey, tt.keys[idx])
				exp := sha3.Sum256(expKey)
				if act.hash != fmt.Sprintf("%x", exp) {
					t.Errorf("wrong sha-3 sum generated %v, expected %v", act.hash, exp)
				}
			}

		})
	}
}

func TestCreateUniqueKeyList(t *testing.T) {
	tests := []struct {
		name string
		key  uint64
		n    int
	}{
		{
			name: "single key",
			key:  123456,
			n:    1,
		},
		{
			name: "several keys",
			key:  123456,
			n:    10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys := createUniqueKeyList(tt.key, tt.n)
			// chech key number
			if len(keys) != tt.n+1 {
				t.Errorf("wrong key number %d, expected %d", len(keys), tt.n+1)
			}

			// check if all keys are unique

			set := make(map[uint64]struct{})
			for _, key := range keys {
				if _, ok := set[key]; ok {
					t.Errorf("key %d is not unique", key)
					continue
				}
				set[key] = struct{}{}

			}
		})
	}
}
