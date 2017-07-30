// https://gobyexample.com/worker-pools
// AG: try out on lonestar before going to yahoo
package main

import "fmt"
import "log"
import "net/http"
import "io/ioutil"

// The worker will run several concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding results on `results`.
// Sleep 1 second per job to simulate an expensive task.
func worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		//time.Sleep(time.Second)
		resp, err := http.Get("http://goanuj.freeshell.org/go/" + j + ".txt")
		if err != nil {
			log.Printf(j + ": " + err.Error())
			fmt.Println("worker", id, "finished job", j)
			results <- err.Error() // channel send
		} else {
			body, _ := ioutil.ReadAll(resp.Body)
			resp.Body.Close() // close ASAP to prevent too many open file desriptors
			val := string(body)
			fmt.Println("worker", id, "finished job", j)
			results <- val // channel send
		}
	}
}

func main() {
	// to use `worker pool, need to send work and then collect their results.
	// 2 channels are needed for this
	jobs := make(chan string, 500)
	results := make(chan string, 500) // AG will need to increase this!

	// Startup N workers, initially blocked because there are no jobs yet.
	for w := 1; w <= 40; w++ {
		go worker(w, jobs, results)
	}

	var sl = []string{"AAPL", "AMZN", "GOOG", "FB", "NFLX"}
	// slice append to repeat stocks; 5 for 160, 6 for 320, 7 for 640
	for i := 0; i < 7; i++ {
		sl = append(sl, sl...)
	}
	// Send N stocks to `jobs channel and then close `jobs channel to signify no more work
	// will cause the for loop in `worker to drop out
	for _, s := range sl {
		jobs <- s
	}
	close(jobs)
	fmt.Println("**** JOBS closed ****")

	// use `results channel to gather return values
	var r = make([]string, len(sl))
	for i := 0; i < len(r); i++ {
		r[i] = <-results
	}
	fmt.Println("\nResults:", len(r))
}
