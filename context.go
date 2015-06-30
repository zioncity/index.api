package main

import (
	"math/rand"
	"time"
)

var letters = []rune("1234567890")

func random_no(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func panic_error(err error) {
	if err != nil {
		panic(err)
	}
}

const (
	cmd = 1000 //degree / 1000
	cm  = 100
)

func random_date() int64 {
	return time.Now().Add(-time.Hour * time.Duration(rand.Intn(24))).Unix()
}
