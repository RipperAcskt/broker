package main

import (
	"github.com/RipperAcskt/broker/internal/app"
	"log"
)

func main() {
	err := app.New().Run()
	if err != nil {
		log.Fatal(err)
	}
}
