package util

import (
	"math/rand"
	"time"
)

func RandomString(n int) string {
	var letters = []byte("kjagfalihgsikufgvbblacgvbsoljhdjguiwghvbcbmcmlahqwphv")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := 0; i < n; i++ {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}
