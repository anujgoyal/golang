package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

var wg sync.WaitGroup

func getStock(x int) {
	for i := 1; i <= 12; i++ {
		fmt.Printf("%d x %d = %d\n", i, x, i*x)
		// time.Sleep(100 * time.Millisecond)
	}
	wg.Done()
}

// https://stackoverflow.com/questions/8757389/reading-file-line-by-line-in-go
func readFile() int {
        // open file
	file, err := os.Open("stocklist.txt")
	if err != nil { log.Fatal(err) }
	defer file.Close()

        // setup slice
        // https://blog.golang.org/go-slices-usage-and-internals
        var sl []string
        sl = make([]string, 0)

        // read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
                sl = append(sl, scanner.Text())
		fmt.Println(sl);
	}

        // more error checking
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
        return 0
}

func main() {
	// don't do error checking yet
        var sl int
        sl = readFile()
        fmt.Printf("sl: %d\n", sl)
	/*for n := 2; n <= 12; n++ {
		wg.Add(1)
		go getStock(n)
	}
	wg.Wait()*/
}
