package main

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jessevdk/go-flags"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var opts struct {
	// Folow Redirects For Requests
	Redirect bool `short:"r" long:"folow-redirect" description:"folow-redirect or not" `
	// Timeout For the http get Requets
	Timeout int `short:"t" default:"5" long:"timeout" description:"Timeout For Requests" `
	// Concurrency For the http requests (something like threads)
	Concurrency int `short:"c" long:"concurrency" default:"40" description:"Concurrency For Requests"`
}
var redirectFunc func(req *http.Request, via []*http.Request) error

func main() {
	// Flage the args to start use it
	_, err := flags.Parse(&opts)
	// Check if the parse occurred any error
	if err != nil {
		fmt.Println("try to use -h to show the options and usage :)")
		return
	}

	// setting timeout
	timeout := time.Duration(opts.Timeout) * time.Millisecond
	concurrency := opts.Concurrency
	redirect := opts.Redirect
	// Get Urls from input
	urls := getlines()

	var wg sync.WaitGroup
	// Start the concurrency and goroutine
	for i := 0; i < concurrency; i++ {
		// Add 1 every concurrency
		wg.Add(1)
		go func() {
			// Done the sync
			defer wg.Done()
			// Start Reading urls from channel
			for domain := range urls {
				// Pass url chan and timeout to return domain and cname
				get(domain, timeout, redirect)

			}
		}()
	}
	// Wait untill workers finish
	wg.Wait()

}

func get(url string, timeout time.Duration, redirect bool) {
	// Check for follow redirect or not
	if redirect == false {
		redirectFunc = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		redirectFunc = nil
	}
	trans := &http.Transport{
		MaxIdleConns:      30,
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
		// Skip Certificate Error
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		TLSHandshakeTimeout: 5 * time.Second,
		Dial: (&net.Dialer{
			Timeout:   timeout * time.Second,
			KeepAlive: time.Second,
		}).Dial,
	}
	client := &http.Client{
		// Passing transport var
		Transport: trans,
		// Follow Redirect
		CheckRedirect: redirectFunc,
		Timeout:       timeout * time.Second,
	}
	// Add https if string didn't contain it
	if !strings.Contains(url, "http") {
		url = "https://" + url
	}

	// Create New Request with url and method
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return
	}
	// Close Connection
	req.Header.Set("Connection", "close")
	// Do the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(url + " | ")
		return
	}
	defer resp.Body.Close()
	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}
	title := string(document.Find("title").Text())
	fmt.Println(url + " | " + title)
}

// Get Urls From File Function
func getlines() <-chan string {
	// Create a channel to store the lines into it
	out := make(chan string)
	// Create Scanner to read lines
	sc := bufio.NewScanner(os.Stdin)
	// looping the scanner
	go func() {
		defer close(out)
		for sc.Scan() {
			// Convert urls to lower and pass to channel
			domain := strings.ToLower(sc.Text())
			// Store lines into channels
			out <- domain
		}
	}()
	return out
}
