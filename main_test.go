package main

import (
	"testing"
)

func TestFetchUrl(t *testing.T) {
	var (
		ch chan FetchResult
		r  FetchResult
	)

	// Bad DNS record
	ch = FetchUrl("http://google.comx/")
	r = <-ch
	if r.err == nil {
		t.Error("http://google.comx/ should not be parsed")
	}

	// Good host.
	ch = FetchUrl("http://google.com/")
	r = <-ch
	if r.err != nil {
		t.Error("http://google.com/ - ", r.err.Error())
	}

	if len(r.links) < 5 {
		t.Error("Found < 5 links on google.com")
	}
}
