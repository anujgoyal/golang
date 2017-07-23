package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

func getStock(s string) {
	var u1 string = "http://download.finance.yahoo.com/d/quotes.csv?s="
	var u2 string = "&f=srp5p6m3m4m6m8j1j4y"
        var url string = u1 + s + u2
        fmt.Printf("url: %s \n", url)
	//wg.Done()
}

// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
func readFile() []string {
	// open file
	file, err := os.Open("stocklist.txt")
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
	sl = readFile()
	fmt.Printf("sl(%d): %s\n", len(sl), sl)
	//for i, s := range sl { fmt.Printf("sl[%d] %s\n", i, s) }

        getStock("AAPL")
	/*for n := 2; n <= 12; n++ {
		wg.Add(1)
		go getStock(n)
	}
	wg.Wait()*/
}
