package main

import (
	"fmt"
	"log"
)

func check(e error, flag string) {
	if e != nil {
		fmt.Printf("Flag: %s\n", flag)
		log.Fatal(e)
	}
}
