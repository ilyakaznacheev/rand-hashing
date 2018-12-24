package hashgen

import (
	"encoding/binary"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"golang.org/x/crypto/sha3"
)

const (
	// randomEndingLen is length of random key ending in next iteration
	randomEndingLen = 4
)

// keyHashRaw contains key number and binary hash sequence
type keyHashRaw struct {
	key  uint32
	hash [32]byte
}

// keyHash contains string key and hex hash string
type keyHash struct {
	key  string
	hash string
}

// startHashing returnes channel with SHA-3 sums of key and n random strings
//
// key - string that has to be hashed
// n - number of hashed strings based on key with random ending
func startHashing(key uint32, n int) chan keyHash {
	keyList := createUniqueKeyList(key, n)
	return generateHash(keyList)
}

// generateHash returns channel with generated hashes from input string list
func generateHash(keyList []uint32) chan keyHash {
	hashChan := make(chan keyHash, len(keyList))
	sumChans := make([]chan keyHashRaw, 0, len(keyList))

	// calculate hash sums
	for _, nextKey := range keyList {
		// result queue
		taskChan := make(chan keyHashRaw)

		// run hash sum calculation
		go func(key uint32) {
			bKey := make([]byte, 4)
			binary.LittleEndian.PutUint32(bKey, key)
			taskChan <- keyHashRaw{
				key:  key,
				hash: sha3.Sum256(bKey),
			}
		}(nextKey)

		// add result channel into result queue
		sumChans = append(sumChans, taskChan)
	}

	// return calculated sums in input list order
	go func() {
		for _, resChan := range sumChans {
			res := <-resChan
			hashChan <- keyHash{
				key:  strconv.Itoa(int(res.key)),
				hash: fmt.Sprintf("%x", res.hash),
			}
		}
		// close channel after reading from all goroutines
		close(hashChan)
	}()

	return hashChan
}

// createUniqueKeyList creates list with a string key
// and n generated strings based on key with random ending
func createUniqueKeyList(key uint32, n int) []uint32 {
	keyToSum := make([]uint32, 0, n+1)
	if n < 0 {
		return keyToSum
	}

	rand.Seed(time.Now().UnixNano())
	keyToSum = append(keyToSum, key)

	for idx := 0; idx < n; idx++ {
		for {
			// replace last 4 digits with random digits
			nextKey := key - key%10000 + uint32(rand.Intn(10000))
			// nextKey := key[0:len(key)-randomEndingLen] + fmt.Sprintf("%04d", randEnding)
			if !intInSlice(nextKey, keyToSum) {
				keyToSum = append(keyToSum, nextKey)
				break
			}
		}
	}
	return keyToSum
}

// intInSlice chechs if string list contains certain string
func intInSlice(num uint32, list []uint32) bool {
	for _, v := range list {
		if v == num {
			return true
		}
	}
	return false
}
