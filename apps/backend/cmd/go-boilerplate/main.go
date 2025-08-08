// Package main is the entry point to the program
package main

import (
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go sample("one", &wg)
	wg.Wait()
}

func sample(payload any, wg *sync.WaitGroup) {
	log.Println("interface", payload)
	wg.Done()
}
