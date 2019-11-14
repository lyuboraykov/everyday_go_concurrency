// inspired by https://eng.uber.com/optimizing-m3/
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/lyuboraykov/everyday_go_concurrency/util"
)

const (
	numPeople           = 40
	numParallelRequests = 5
	timeout             = 100 * time.Second
)

type reqRes struct {
	name string
	err  error
}

type runner struct {
	in  chan int
	out chan reqRes
}

func (r *runner) Go(p int) {
	r.in <- p
}

func (r *runner) ResChan() chan reqRes {
	return r.out
}

func (r *runner) Close() {
	close(r.in)
}

func NewRunner(numWorkers int) *runner {
	in := make(chan int)
	out := make(chan reqRes)

	for i := 0; i < numWorkers; i++ {
		go func() {
			for {
				if p, ok := <-in; ok {
					name, err := util.FetchName(p)
					out <- reqRes{name, err}
				} else {
					return
				}
			}
		}()
	}

	return &runner{in, out}
}

func main() {
	util.Timed(dowork)()
}

func dowork() {
	ctx := context.Background()
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, timeout)
	defer cancelFunc()

	people := make([]int, numPeople)
	for i := 0; i < numPeople; i++ {
		people[i] = i + 1
	}

	r := NewRunner(numParallelRequests)

	go func() {
		for _, p := range people {
			r.Go(p)
		}
	}()

peopleloop:
	for range people {
		select {
		case res := <-r.ResChan():
			if res.err != nil {
				fmt.Println("there was an error: ", res.err)
			} else {
				fmt.Println(res.name)
			}
		case <-timeoutCtx.Done():
			fmt.Println("timed out:", timeoutCtx.Err())
			break peopleloop
		}
	}

	r.Close()
}
