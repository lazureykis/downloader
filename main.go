package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	concurrency int
	urls        []string
)

func main() {
	flag.IntVar(&concurrency, "c", 20, "Concurrency level.")
	flag.Parse()
	urls = flag.Args()

	if len(urls) == 0 {
		log.Fatalln("Pass some urls as arguments.")
	}

	fmt.Println("urls:", flag.Args())
	fmt.Println(" con:", concurrency)
	fmt.Println("TODO: implement fetcher")
}
