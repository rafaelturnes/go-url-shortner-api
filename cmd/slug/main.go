package main

import (
	"fmt"
	"math/rand"
)

const (
	charset    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	mediumSize = 6
	highSize   = 12
)

func main() {
	n := 10

	fmt.Println("--- Short ---")
	for i := 0; i < n; i++ {
		fmt.Println(short(charset))
	}

	fmt.Println("--- Secure ---")
	for i := 0; i < n; i++ {
		fmt.Println(secure(charset))
	}
}

func short(charset string) string {
	b := make([]byte, mediumSize)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func secure(charset string) string {
	b := make([]byte, highSize)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}
