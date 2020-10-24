package main

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

var urls = map[string][]string{
	"http://news.yahoo.com/news/topics/": []string{"http://news.yahoo.com",
		"http://news.yahoo.com/news",
		"http://news.google.com",
		"http://news.yahoo.com/us"},
	"http://news.yahoo.com/news": []string{
		"http://news.yahoo.com",
		"http://news.yahoo.com/news/topics/",
		"http://news.yahoo.com/us"},
	"http://news.yahoo.com": []string{
		"http://news.yahoo.com/news",
		"http://news.yahoo.com/us",
		"http://news.yahoo.com/new"},
}

var source = "http://news.yahoo.com/news/topics/"

func main() {
	s := Spider{mux: &sync.RWMutex{}, sites: map[string]bool{source: true}}
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(1)
	go s.crawl(source, &wg)
	wg.Wait()

	fmt.Println("Finished in ", time.Since(start).Milliseconds(), " ms")

	fmt.Println("Sites found:")
	for k := range s.sites {
		fmt.Println(k)
	}
}

// Spider crawls sites for connected sites in domain
type Spider struct {
	mux   *sync.RWMutex
	sites map[string]bool
}

func (s *Spider) crawl(source string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, url := range getURLs(source) {
		if !s.inList(url) && inDomain(url, source) {
			s.addToList(url)
			wg.Add(1)
			go s.crawl(url, wg)
		}
	}
}

func (s *Spider) inList(url string) bool {
	s.mux.RLock()
	defer s.mux.RUnlock()
	if s.sites[url] {
		return true
	}
	return false
}

func inDomain(url, source string) bool {
	parts := strings.Split(source, "/")
	domain := ""
	for _, p := range parts {
		if strings.Contains(p, ".") {
			domain = p
			break
		}
	}

	if strings.Contains(url, domain) {
		return true
	}
	return false
}

func (s *Spider) addToList(url string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	s.sites[url] = true
}

// simulates latency to scrape the page
func getURLs(url string) []string {
	time.Sleep(15 * time.Millisecond)
	return urls[url]
}
