package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
)

func MakeRequest(client *http.Client, url string) {
	data := ""
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBufferString(data))
	if err != nil {
		fmt.Println(err)
	}
	request.Header.Add("Content-Type", "text/plain")
	request.Header.Add("Content-Length", strconv.Itoa(len(data)))

	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}

	defer response.Body.Close()
}
