// https://gobyexample.com/worker-pools
package main

import "fmt"
import "time"

// The worker will run several concurrent instances. These workers will receive
// work on the `jobs` channel and send the corresponding results on `results`. 
// Sleep 1 second per job to simulate an expensive task.
func worker(id int, jobs <-chan string, results chan<- string) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
                results <- ("MSG: " + j)
	}
}

func main() {
	// In order to use our pool of workers we need to send them work and 
        // collect their results. We make 2 channels for this.
	jobs := make(chan string, 100)
	results := make(chan string, 100)

	// Startup 3 workers, initially blocked because there are no jobs yet.
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// Send 5 `jobs` and then `close` that channel to indicate that's all the work we have.
        var sl = []string{"AAPL", "AMZN", "GOOG", "FB", "NFLX"}
        for _, s := range sl {
		jobs <- s
	}
	close(jobs)

	// Finally collect all the results of the work.
        var r = make([]string, len(sl))
	for i := 0; i < len(r); i++ {
		r[i] = <-results
	}
        fmt.Println("Results:", r)
}

