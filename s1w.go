package main

import (
    "fmt"
    "log"
    "sync"
    "net/http"
    "io/ioutil"
    "time"
)

var wg sync.WaitGroup

func getStock(s string) {
        resp, err := http.Get("http://goanuj.freeshell.org/go/" + s + ".txt")
	if err != nil {
            log.Printf(s + ": "+ err.Error())
            wg.Done()
            return
        }
	body, _:= ioutil.ReadAll(resp.Body)
        resp.Body.Close() // close ASAP to prevent too many open file desriptors
        val := string(body)
        fmt.Printf("body: %s", val)
	wg.Done()
}

func main() {
        start := time.Now()
	var sl = []string{"AAPL","AMZN","GOOG","FB","NFLX"}
        for i := 0; i<10; i++ {
            sl = append(sl, sl)
        }
        fmt.Printf("sl(size): %d\n", len(sl))

	for _, s := range sl {
		wg.Add(1)
		go getStock(s)
	}
	wg.Wait()
        fmt.Printf("main: %.2fs elapsed.\n", time.Since(start).Seconds())
}

