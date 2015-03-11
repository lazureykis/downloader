package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
)

type (
	FetchResult struct {
		url   string
		err   error
		links []string
	}
)

var (
	concurrency int
	url         string
)

func main() {
	flag.IntVar(&concurrency, "c", 20, "Concurrency level.")
	flag.StringVar(&url, "u", "", "URL to fetch")
	flag.Parse()

	if len(url) == 0 {
		log.Fatalln("Missing argument -u")
	}

	fmt.Println("Concurrent threads:", concurrency)
	fmt.Println("URL:", url)
	ch := FetchUrl(url)
	fr := <-ch
	if fr.err != nil {
		log.Fatal(fr.err.Error())
	} else {
		fmt.Println("Parsed ", fr.url)
		fmt.Println("Found", len(fr.links), "links:")
		for _, link := range fr.links {
			fmt.Println(link)
		}
	}
}

func FetchUrl(url string) chan FetchResult {
	ch := make(chan FetchResult)

	go func() {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			ch <- FetchResult{"", err, nil}
		}

		links := make([]string, 0)
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			links = append(links, link)
		})

		ch <- FetchResult{url, nil, links}
	}()

	return ch
}
