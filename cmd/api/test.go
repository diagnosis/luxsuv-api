package main

import (
	"log"

	"github.com/diagnosis/luxsuv-api-v2/internal/secure"
)

func main() {
	hash, err := secure.HashPassword("password1")
	if err != nil {
		log.Fatal(err)
	}
	g := secure.VerifyPassword("password1", hash)
	log.Printf("hash:%v result:%v", hash, g)

}
