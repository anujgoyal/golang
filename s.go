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
	//for i := 1; i <= 12; i++ {
        fmt.Printf("%s \n", s)
        // time.Sleep(100 * time.Millisecond)
	wg.Done()
}

// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
func readFile() []string {
        // open file
	file, err := os.Open("stocklist.txt")
	if err != nil { log.Fatal(err) }
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
        //fmt.Printf("sl(%d): %s\n", len(sl), sl)
        for i, s := range sl {
            fmt.Printf("sl[%d] %s\n", i, s)
        }

	/*for n := 2; n <= 12; n++ {
		wg.Add(1)
		go getStock(n)
	}
	wg.Wait()*/
}

