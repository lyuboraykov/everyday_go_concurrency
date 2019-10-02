package main

import (
	"fmt"

	"github.com/lyuboraykov/everyday_go_concurrency/util"
)

func main() {
	util.Timed(dowork)()
}

func dowork() {
	name1, err1 := util.FetchName(1)
	name2, err2 := util.FetchName(2)

	if err1 != nil || err2 != nil {
		fmt.Println("Error occurred: ", err1, err2)
	} else {
		fmt.Println(name1, name2)
	}
}
