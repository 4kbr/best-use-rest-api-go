package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// create a new http client
	client := &http.Client{}

	// url := "https://jsonplaceholder.typicode.com/posts/1"
	url := "https://swapi.dev/api/people/1"
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("error making GET request", err)
	}
	defer resp.Body.Close()

	// read and print the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("error reading response body:", err)
		return
	}

	// ini hasilnya  [123 10 32 32 34 ....]
	fmt.Println(body)
	// ini hasilnya  {"userId": 1,"id": 1,"title": "sunt aut facere repellat provident occaecati excepturi o",...}
	fmt.Println(string(body))

}
