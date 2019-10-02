package main

import (
	"fmt"

	"github.com/lyuboraykov/everyday_go_concurrency/util"
)

const numPeople = 40
const numParallelRequests = 3

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

	peopleChan := make(chan int)

	go func() {
		for _, p := range people {
			peopleChan <- p
		}
	}()

	for i := 0; i < numParallelRequests; i++ {
		go func() {
			for p := range peopleChan {
				name, err := util.FetchName(p)
				results <- reqRes{name, err}
			}
		}()
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
