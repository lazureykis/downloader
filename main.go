package main

import (
	"flag"
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

	work(Url, Concurrency)
}

func work(url string, concurrency int) {
	log.Println("Concurrent threads:", concurrency)
	log.Println("URL:", url)
	ch := FetchUrl(url)
	fr := <-ch
	processResult(fr)
}

func processResult(fr FetchResult) {
	if fr.err != nil {
		log.Println(fr.err.Error())
	} else {
		log.Println("Parsed", fr.url)
		log.Println("Found", len(fr.links), "links:")
		for _, link := range fr.links {
			log.Println(link)
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
