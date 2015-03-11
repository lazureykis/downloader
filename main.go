package main

import (
	"flag"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"net/url"
	"time"
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
	Timeout     = 500 * time.Millisecond
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

	select {
	case fr := <-ch:
		processResult(fr)
	case <-time.After(Timeout * time.Millisecond):
		log.Println("timed out")
	}
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

func FetchUrl(urlString string) chan FetchResult {
	ch := make(chan FetchResult)
	base_url, err := url.Parse(urlString)

	go func() {
		if err != nil {
			ch <- FetchResult{"", err, nil}
			return
		}

		res, err := http.Get(urlString)
		if err != nil {
			ch <- FetchResult{"", err, nil}
			return
		}

		doc, err := goquery.NewDocumentFromResponse(res)
		if err != nil {
			ch <- FetchResult{"", err, nil}
			return
		}

		links := make([]string, 0)
		doc.Find("a[href]").Each(func(i int, s *goquery.Selection) {
			link, _ := s.Attr("href")
			link_url, err := url.Parse(link)
			if err == nil {
				real_link := base_url.ResolveReference(link_url).String()
				links = append(links, real_link)
			}
		})

		ch <- FetchResult{urlString, nil, links}
		return
	}()

	return ch
}
