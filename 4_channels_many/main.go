package main

import (
	"fmt"

	"github.com/lyuboraykov/everyday_go_concurrency/util"
)

const numPeople = 40

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

	for _, p := range people {
		go func(p int) {
			name, err := util.FetchName(p)
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
