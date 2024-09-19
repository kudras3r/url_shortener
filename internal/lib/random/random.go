package random

import (
	"math/rand"
)

var (
	symbols    []rune = []rune("qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM")
	symbolsLen        = len(symbols)
)

func NewRandomString(strLen int) string {
	var res []rune = make([]rune, strLen)

	for i := range res {
		res[i] = symbols[rand.Intn(symbolsLen)]
	}

	return string(res)
}
