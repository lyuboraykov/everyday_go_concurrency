package main

import (
	"context"
	"fmt"

	"github.com/lyuboraykov/everyday_go_concurrency/util"
	"golang.org/x/sync/semaphore"
)

const numPeople = 40
const numParallelRequests = 5

func main() {
	util.Timed(dowork)()
}

func dowork() {
	type reqRes struct {
		name string
		err  error
	}
	results := make(chan reqRes)

	people := make([]int, numPeople)
	for i := 0; i < numPeople; i++ {
		people[i] = i + 1
	}

	sem := semaphore.NewWeighted(numParallelRequests)
	ctx := context.Background()
	for _, p := range people {
		go func(p int) {
			_ = sem.Acquire(ctx, 1)
			name, err := util.FetchName(p)
			sem.Release(1)
			results <- reqRes{name, err}
		}(p)
	}

	for range people {
		res := <-results
		if res.err != nil {
			fmt.Println("there was an error: ", res.err)
		} else {
			fmt.Println(res.name)
		}
	}

	close(results)
}
