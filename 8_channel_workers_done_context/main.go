package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const numPeople = 40
const numParallelRequests = 3

func main() {
	ctx := context.Background()
	timeoutCtx, cancelFunc := context.WithTimeout(ctx, 3*time.Second)
	defer cancelFunc()

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
				name, err := makeRequest(p)
				results <- reqRes{name, err}
			}
		}()
	}

peopleloop:
	for range people {
		select {
		case res := <-results:
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

	close(results)
}

func makeRequest(personId int) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://swapi.co/api/people/%d", personId))
	if err != nil {
		return "", err
	}

	jsonDecoder := json.NewDecoder(resp.Body)
	type person struct {
		Name string `json:"name"`
	}
	p := person{}
	err = jsonDecoder.Decode(&p)
	if err != nil {
		return "", err
	}
	return p.Name, err
}
