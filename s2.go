package main

// read a file of stock symbols
// gather info from website using goroutines

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
        "time"
)

// using init doesn't seem to be helping
var netClient = &http.Client{}

func init() {
        fmt.Println("init in s2.go")
        // setup Client
	tr := &http.Transport{
		//IdleConnTimeout:  time.Second*5,
		MaxIdleConns:       20,
                MaxIdleConnsPerHost:  20,
		//dDisableCompression: true,
	}
	netClient = &http.Client{Transport: tr}

        // set connection timeout
        timeout := time.Duration(1 * time.Second)
        client := http.Client{
            Timeout: timeout,
        }
        _ = client
}

func main() {
	// read file, and output slice of ticker symbols
	var sl []string
	sl = readFile("stocklist.txt")
	sl = readFile("wilshire.txt")

	// create channel, create slice
	var c chan string
	c = make(chan string)
	// fire off N routines
	for _, s := range sl {
		go getStock(s, c)
	}

	info := make([]string, 0) // create slice
	// get info from channel
	for i := 0; i < len(sl); i++ {
		info = append(info, <-c)
	}

        // print stats
	fmt.Printf("\n\ninfo(%d):\n%s\nmain done.\n", len(info), info)
}

func getStock(s string, c chan string) {
        // HTTP get stock symbol
	fmt.Printf("%s ", s)
	resp, err := http.Get("http://download.finance.yahoo.com/d/quotes.csv?s=" + s + "&f=srp5p6m3m4m6m8j1j4y")
        // resp, err := http.Get("http://download.finance.yahoo.com/d/quotes.csv?s=" + s + "&f=srp5p6m3m4m6m8j1j4y")
	if err != nil {
		//log.Fatalf("http.Get => %v", err.Error())
		log.Printf("http.Get => %v", err.Error())
		c <- err.Error()
                return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body) // read body
	if err != nil {
		log.Fatalf("ioutil.ReadAll => %v", err.Error())
	}
	c <- string(body) // https://stackoverflow.com/questions/38673673/access-http-response-as-string-in-go
        return
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
