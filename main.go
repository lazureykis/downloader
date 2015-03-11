package main

import (
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
)

type (
	FetchResult struct {
		url   string
		err   error
		links []string
	}
)

var (
	Concurrency int
	Url         string
)

func main() {
	flag.IntVar(&Concurrency, "c", 20, "Concurrency level.")
	flag.StringVar(&Url, "u", "", "URL to fetch")
	flag.Parse()

	if len(Url) == 0 {
		log.Fatalln("Missing argument -u")
	}

	fmt.Println("Concurrent threads:", Concurrency)
	fmt.Println("URL:", Url)
	ch := FetchUrl(Url)
	fr := <-ch
	if fr.err != nil {
		log.Fatal(fr.err.Error())
	} else {
		fmt.Println("Parsed", fr.url)
		fmt.Println("Found", len(fr.links), "links:")
		for _, link := range fr.links {
			fmt.Println(link)
		}
	}
}

func FetchUrl(url string) chan FetchResult {
	ch := make(chan FetchResult)

	go func() {
		res, err := http.Get(url)
		if err != nil {
			ch <- FetchResult{"", err, nil}
		}

		doc, err := goquery.NewDocumentFromResponse(res)
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
