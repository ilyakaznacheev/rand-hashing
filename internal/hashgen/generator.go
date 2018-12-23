package hashgen

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/sha3"
)

const (
	// randomEndingLen is length of random key ending in next iteration
	randomEndingLen = 4
)

type keyHash struct {
	key  string
	hash [32]byte
}

// startHashing returnes channel with SHA-3 sums of key and n random strings
//
// key - string that has to be hashed
// n - number of hashed strings based on key with random ending
func startHashing(key string, n int) chan keyHash {
	keyList := createUniqueKeyList(key, n)
	return generateHash(keyList)
}

// generateHash returns channel with generated hashes from input string list
func generateHash(keyList []string) chan keyHash {
	hashChan := make(chan keyHash, len(keyList))
	sumChans := make([]chan keyHash, 0, len(keyList))

	// calculate hash sums
	for _, nextKey := range keyList {
		// result queue
		taskChan := make(chan keyHash)

		// run hash sum calculation
		go func(key string) {
			taskChan <- keyHash{
				key:  key,
				hash: sha3.Sum256([]byte(key)),
			}
		}(nextKey)

		// add result channel into result queue
		sumChans = append(sumChans, taskChan)
	}

	// return calculated sums in input list order
	go func() {
		for _, resChan := range sumChans {
			hashChan <- <-resChan
		}
		// close channel after reading from all goroutines
		close(hashChan)
	}()

	return hashChan
}

// createUniqueKeyList creates list with a string key
// and n generated strings based on key with random ending
func createUniqueKeyList(key string, n int) []string {
	keyToSum := make([]string, 0, n+1)
	if n < 0 {
		return keyToSum
	}

	rand.Seed(time.Now().UnixNano())
	keyToSum = append(keyToSum, key)

	for idx := 0; idx < n; idx++ {
		for {
			randEnding := rand.Intn(1000)
			nextKey := key[0:len(key)-randomEndingLen] + fmt.Sprintf("%04d", randEnding)
			if !stringInSlice(nextKey, keyToSum) {
				keyToSum = append(keyToSum, nextKey)
				break
			}
		}
	}
	return keyToSum
}

// stringInSlice chechs if string list contains certain string
func stringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
