package storage

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

func calculeIndex(key string, arrLength int) int {
	if len(key) > 100 {
		key = key[0:100]
	}

	var total int
	const WEIRD_PRIME = 31
	for len(key) > 0 {
		r, size := utf8.DecodeRuneInString(key)
		rInt, _ := strconv.Atoi(fmt.Sprintf("%d", r))
		total = (total*WEIRD_PRIME + rInt) % arrLength
		key = key[size:]
	}

	fmt.Println("the total: ", total)

	return total
}
