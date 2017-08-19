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

// Generator: function that returns a channel; https://www.youtube.com/watch?v=f6kdp27TYZs (15m)
func getStocks(sl []string) <-chan string {
	c := make(chan string)
	limit := make(chan struct{}, 200) // limit to N parallel operations
	for _, s := range sl {
		limit <- struct{}{}
		go getStock(s, c, limit)
	}
	return c
}

func getStock(s string, c chan string, limit chan struct{}) {
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
	// creates slice of 8(1280), 9(2560), 10(5120) , 11(10240) elements
	for i := 0; i < 11; i++ {
		sl = append(sl, sl...)
	}
	fmt.Printf("main: %.2fs elapsed.\n", time.Since(start).Seconds())
	fmt.Printf("sl(size): %d\n", len(sl))

	// generator pattern, get channel back from function
	var c = getStocks(sl)
	var r = make([]string, len(sl), len(sl)) 
	for i := 0; i < len(r); i++ {
                r[i] = <-c // channel recv
		//fmt.Printf("%s", r[i]) // channel recv
	}

        fmt.Printf("last element: %s\n", r[len(r)-1])
	fmt.Printf("main: %.2fs elapsed.\n", time.Since(start).Seconds())
}
