package main

import (
	"testing"
)

func TestFetchUrl(t *testing.T) {
	var r FetchResult

	// Bad DNS record
	r = FetchUrl("http://google.comx/")
	if r.err == nil {
		t.Error("http://google.comx/ should not be parsed")
	}

	// Good host.
	r = FetchUrl("http://google.com/")
	if r.err != nil {
		t.Error("http://google.com/ - ", r.err.Error())
	}

	if len(r.links) < 5 {
		t.Error("Found < 5 links on google.com")
	}
}
