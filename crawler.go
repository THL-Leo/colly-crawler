package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/extensions"
)

const (
	maxDepth      = 16
	maxPages      = 20000
	maxGoroutines = 625
	startURL      = "https://www.latimes.com/"
)

func main() {
	visited := make(map[string]bool)
	var mu sync.RWMutex
	pageCount := 0

	fetchFile, visitFile, urlsFile := createOutputFiles()
	defer closeFiles(fetchFile, visitFile, urlsFile)

	fetchWriter := csv.NewWriter(fetchFile)
	visitWriter := csv.NewWriter(visitFile)
	urlsWriter := csv.NewWriter(urlsFile)
	defer flushWriters(fetchWriter, visitWriter, urlsWriter)

	writeHeaders(fetchWriter, visitWriter, urlsWriter)

	c := colly.NewCollector(
		colly.MaxDepth(maxDepth),
		colly.Async(true),
		colly.AllowedDomains("www.latimes.com", "*.latimes.com"),
	)

	extensions.RandomUserAgent(c)
	// rp, err := proxy.RoundRobinProxySwitcher("socks5://127.0.0.1:8880", "socks5://127.0.0.1:8881")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// c.SetProxyFunc(rp)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: maxGoroutines,
		RandomDelay: 10 * time.Second,
	})

	c.SetRequestTimeout(30 * time.Second)

	c.SetRedirectHandler(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	})

	c.OnRequest(func(r *colly.Request) {
		mu.Lock()
		defer mu.Unlock()
		if pageCount >= maxPages {
			r.Abort()
			return
		}
		if visited[r.URL.String()] {
			r.Abort()
			return
		}
		visited[r.URL.String()] = true
		pageCount++
	})

	c.OnResponse(func(r *colly.Response) {
		mu.Lock()
		defer mu.Unlock()
		fetchWriter.Write([]string{r.Request.URL.String(), fmt.Sprintf("%d", r.StatusCode)})
	})

	var outlinksCount int

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Request.AbsoluteURL(e.Attr("href"))
		if link == "" {
			return
		}

		mu.Lock()
		defer mu.Unlock()

		isSameDomain := isLatimesDomain(e.Request.URL.Host)
		urlsWriter.Write([]string{link, indicator(isSameDomain)})

		outlinksCount++

		if isSameDomain {
			e.Request.Visit(link)
		}
	})

	c.OnScraped(func(r *colly.Response) {
		mu.Lock()
		defer mu.Unlock()

		// Update the visit record with the correct outlinks count
		visitWriter.Write([]string{
			r.Request.URL.String(),
			fmt.Sprintf("%d", len(r.Body)),
			fmt.Sprintf("%d", outlinksCount),
			strings.Split(r.Headers.Get("Content-Type"), ";")[0],
		})

		// Reset outlinks count for the next page
		outlinksCount = 0
	})

	c.OnError(func(r *colly.Response, err error) {
		mu.Lock()
		defer mu.Unlock()
		if r.StatusCode == 0 {
			log.Printf("Error on %s: %s", r.Request.URL, err)
			pageCount--
		} else {
			fetchWriter.Write([]string{r.Request.URL.String(), fmt.Sprintf("%d", r.StatusCode)})
		}
	})

	startTime := time.Now()

	c.Visit(startURL)
	c.Wait()

	duration := time.Since(startTime)
	pagesPerMinute := float64(pageCount) / duration.Minutes()

	fmt.Printf("All tasks completed. Pages crawled: %d\n", pageCount)
	fmt.Printf("Pages per minute: %.2f\n", pagesPerMinute)
}

func createOutputFiles() (*os.File, *os.File, *os.File) {
	fetchFile, _ := os.Create("./fetch_latimes.csv")
	visitFile, _ := os.Create("./visit_latimes.csv")
	urlsFile, _ := os.Create("./urls_latimes.csv")
	return fetchFile, visitFile, urlsFile
}

func closeFiles(files ...*os.File) {
	for _, file := range files {
		file.Close()
	}
}

func flushWriters(writers ...*csv.Writer) {
	for _, writer := range writers {
		writer.Flush()
	}
}

func writeHeaders(fetchWriter, visitWriter, urlsWriter *csv.Writer) {
	fetchWriter.Write([]string{"URL", "Status"})
	visitWriter.Write([]string{"URL", "Size(Bytes)", "Outlinks", "Content-Type"})
	urlsWriter.Write([]string{"URL", "Indicator"})
}

func isLatimesDomain(host string) bool {
	return host == "www.latimes.com" || host == "latimes.com"
}

func indicator(isSameDomain bool) string {
	if isSameDomain {
		return "OK"
	}
	return "N_OK"
}
