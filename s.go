package main

// read a file of stock symbols
// gather info from website using goroutines

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
        "net/http"
        "io/ioutil"
)

var wg sync.WaitGroup

func getStock(s string) string {
	// https://stackoverflow.com/questions/1760757/how-to-efficiently-concatenate-strings-in-go?rq=1
	var u1 string = "http://download.finance.yahoo.com/d/quotes.csv?s="
	var u2 string = "&f=srp5p6m3m4m6m8j1j4y"
	var url string = u1 + s + u2
	//fmt.Printf("url: %s \n", url) // debugging
	fmt.Printf("s: %s", s)

        // get url, https://golang.org/pkg/net/http/
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("http.Get => %v", err.Error())
	}
        defer resp.Body.Close()

        //https://stackoverflow.com/questions/16784419/in-golang-how-to-determine-the-final-url-after-a-series-of-redirects
        //finalURL := resp.Request.URL.String()
        //fmt.Printf("The URL you ended up at is: %v\n", finalURL)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll => %v",err.Error())
	}
        // https://stackoverflow.com/questions/38673673/access-http-response-as-string-in-go
        val := string(body)
        fmt.Printf("body: %s", val)
	wg.Done()
        return val
}

// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
func readFile(filename string) []string {
	// open file
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// setup slice, https://blog.golang.org/go-slices-usage-and-internals
	var sl []string
	sl = make([]string, 0)

	// read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sl = append(sl, scanner.Text())
		// fmt.Println(sl);
	}

	// more error checking
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return sl
}

func main() {
	// slice of ticker symbols
	var sl []string
	sl = readFile("stocklist.txt")
	sl = readFile("wilshire.txt")
        _ = sl
	//fmt.Printf("sl(%d): %s\n", len(sl), sl)
	//for i, s := range sl { fmt.Printf("sl[%d] %s\n", i, s) }

	//getStock("AAPL")
	for i, s := range sl {
                _ = i
		wg.Add(1)
		go getStock(s)
	}
	wg.Wait()
}
