package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// MakeRequest makes a request for a person from star wars
func MakeRequest(personId int) (string, error) {
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

// Timed prints the duration of the execution of the function
func Timed(f func()) func() {
	return func() {
		start := time.Now()
		f()
		fmt.Println(time.Since(start))
	}
}
