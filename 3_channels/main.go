package main

import (
	"fmt"

	"github.com/lyuboraykov/everyday_go_concurrency/util"
)

func main() {
	util.Timed(dowork)()
}

func dowork() {
	type reqRes struct {
		name string
		err  error
	}
	results := make(chan reqRes)

	people := []int{1, 2, 3}

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
