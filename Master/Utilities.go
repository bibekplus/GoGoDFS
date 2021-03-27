package main

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"time"
)

func isUnique(list []int, a int) bool {
	for _, b := range list {
		if b == a {
			return false
		}
	}
	return true
}


// newUUID generates a random UUID according to RFC 4122
func newUUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(crand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits;
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random);
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}

func sample(data []string, n int) []string{
	var i int
	size := len(data)

	value := make([]string, n)

	for i=0; i < n; i++{
		value[i] = data[i]
	}


	x1 := rand.NewSource(time.Now().UnixNano())
	y1 := rand.New(x1)

	jUnique :=[]int{-1}
	for ; i < size; i++{

		j:= y1.Intn(i+1)
		for !isUnique(jUnique, j) {
			j = y1.Intn(i+1)
		}
		_ = append(jUnique, j)

		if j < n{
			value[j] = data[i]
		}
	}
	return value


}


