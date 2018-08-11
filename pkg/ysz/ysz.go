package ysz

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	seed_str, ok := os.LookupEnv("YSZ_SEED")
	if ok {
		seed, err := strconv.Atoi(seed_str)
		if err != nil {
			panic(err)
		}
		rand.Seed(int64(seed))
	} else {
		rand.Seed(time.Now().UnixNano())
	}
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
