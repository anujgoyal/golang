// https://gobyexample.com/worker-pools
// AG: try out on lonestar before going to yahoo
package main

import "fmt"
import "log"
import "net/http"
import "io/ioutil"

// The worker will run several concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding results on `results`.
func worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		// fmt.Println("worker", id, "started  job", j)
		// the expensive task
		resp, err := http.Get("http://goanuj.freeshell.org/go/" + j + ".txt")
		if err != nil {
			log.Printf(j + ": " + err.Error())
			//fmt.Println("worker", id, "finished job", j)
			results <- err.Error() // channel send
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close() // close ASAP to prevent too many open file desriptors
			val := string(body)
			//fmt.Println("worker", id, "finished job", j)
			results <- val // channel send
		}
	}
}

func main() {
	// to use `worker pool, need to send work and then collect the results.
	// 2 channels are needed todo this
	// AG: heuristic, channel queue should be approx. the number of stocks
	// NOTE: if jobs queue isn't large enough then a deadlock will occur, this is very subtle
	// basically the jobs channel is blocking because the queue is too small and the code to
	// close the `jobs channel doesn't get executed, and neither does the `results channel for for loop
	jobs := make(chan string, 8000)
	results := make(chan string, 8000) // AG will need to increase this!

	// startup N workers, initially blocked because there are no jobs yet.
	// work network: highest is ~2000 workers,
	// home network: highest is ~250 workers, ... but why?
	for w := 1; w <= 230; w++ {
		go worker(w, jobs, results)
	}

	var sl = []string{"AAPL", "AMZN", "GOOG", "FB", "NFLX"}
	// slice append to repeat stocks
	// 5 for 160, 6 for 320, 7 for 640, 8 for 1280, 9 for 2540, 10 for 5080
	for i := 0; i < 10; i++ {
		sl = append(sl, sl...)
	}
	// send N stocks to `jobs channel and then
	// close `jobs channel to signify no more work,
	// that will cause the for loop in `worker to drop out
	for _, s := range sl {
		jobs <- s
	}
	close(jobs)
	fmt.Println("*** JOBS closed ***")

	// use `results channel to gather return values
	var r = make([]string, len(sl))
	for i := 0; i < len(r); i++ {
		r[i] = <-results // channel recv
	}
	fmt.Println("\nResults:", len(r))
	//fmt.Println("\nResults:", r)
}
