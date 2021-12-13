package main

import (
	"flag"
	"log"
	"sync"

	"github.com/flier/gohs/hyperscan"
)

func main() {
	patfile := flag.String("patterns", "patterns", "Path to pattern file")
	infile := flag.String("input", "input", "Path to input file (list of URLs)")
	outfile := flag.String("output", "output", "Path to output file")
	workers := flag.Int("workers", 5, "Amount of workers to spawn")
	flag.Parse()

	pat, err := ReadFile(*patfile)
	if err != nil {
		log.Fatalln("Error reading pattern file:", err)
	}
	urls, err := ReadFile(*infile)
	if err != nil {
		log.Fatalln("Error reading input file:", err)
	}

	p := make([]*hyperscan.Pattern, 0)
	for _, i := range pat {
		p = append(p, hyperscan.NewPattern(i, hyperscan.SomLeftMost))
	}

	db, err := hyperscan.NewStreamDatabase(p...)
	if err != nil {
		log.Fatalln("Error creating Hyperscan DB:", err)
	}
	defer db.Close()

	log.Printf("%d patterns loaded", len(p))
	log.Printf("%d URLs queued", len(urls))
	log.Printf("Running with %d workers", *workers)

	var wg sync.WaitGroup
	recv := make(chan string)
	pool := make(chan struct{}, *workers)

	go func() {
		err := WriteResults(*outfile, recv)
		if err != nil {
			log.Fatalln("Error creating output file:", err)
		}
	}()

	for _, i := range urls {
		wg.Add(1)
		go func(i string) {
			defer wg.Done()
			pool <- struct{}{}

			log.Println("Scanning", i)
			err := Scan(i, db, recv)
			if err != nil {
				log.Println("Error scanning:", err)
			}

			<-pool
		}(i)
	}

	wg.Wait()
	defer close(recv)
}
