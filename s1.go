package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// https://www.youtube.com/watch?v=f6kdp27TYZs (15m)
// Generator: function that returns a channel
func getStocks(sl []string) <-chan string {
	c := make(chan string)
	for _, s := range sl {
		go getStock(s, c)
	}
	return c
}

func getStock(s string, c chan string) {
	resp, err := http.Get("http://goanuj.freeshell.org/go/" + s + ".txt")
	if err != nil {
		log.Printf(s + ": " + err.Error())
		c <- err.Error() // channel send
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close() // close ASAP to prevent too many open file desriptors
	val := string(body)
	//fmt.Printf("body: %s", val)
	c <- val // channel send
	return
}

func main() {
	start := time.Now()
	var sl = []string{"AAPL", "AMZN", "GOOG", "FB", "NFLX"}
	// creates slice of 1280 elements
	for i := 0; i < 8; i++ {
		sl = append(sl, sl...)
	}
	fmt.Printf("sl(size): %d\n", len(sl))

	// get channel that returns only strings
	c := getStocks(sl)
	for i := 0; i < len(sl); i++ {
		fmt.Printf("%s", <-c) // channel recv
	}

	fmt.Printf("main: %.2fs elapsed.\n", time.Since(start).Seconds())
}
