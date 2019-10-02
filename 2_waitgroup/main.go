package main

import (
	"fmt"
	"sync"

	"github.com/lyuboraykov/everyday_go_concurrency/util"
)

func main() {
	util.Timed(dowork)()
}

func dowork() {
	wg := sync.WaitGroup{}

	var name1, name2 string
	var err1, err2 error
	wg.Add(2)
	go func() {
		name1, err1 = util.FetchName(1)
		wg.Done()
	}()
	go func() {
		name2, err2 = util.FetchName(2)
		wg.Done()
	}()
	wg.Wait()

	if err1 != nil || err2 != nil {
		fmt.Println("Error occurred: ", err1, err2)
	} else {
		fmt.Println(name1, name2)
	}
}
