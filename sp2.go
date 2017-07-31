// https://stackoverflow.com/questions/45332170/golang-strategies-to-prevent-connection-reset-by-peer-errors/45338264#45338264
// https://play.golang.org/p/WJdrYa-A4C
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
	limit := make(chan struct{}, 20) // limit to 20 parallel operations
	for _, s := range sl {
		limit <- struct{}{}
		go getStock(s, c, limit)
	}
	return c
}

func getStock(s string, c chan string, limit chan struct{}) {
	// time.Sleep(500 * time.Millisecond)
	resp, err := http.Get("http://goanuj.freeshell.org/go/" + s + ".txt")
	if err != nil {
		log.Printf(s + ": " + err.Error())
		<-limit          // release limiting resource
		c <- err.Error() // channel send
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close() // close ASAP to prevent too many open file desriptors
	val := string(body)
	<-limit // release limiting resource
	c <- val
	return
}

func main() {
	var start = time.Now()
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
