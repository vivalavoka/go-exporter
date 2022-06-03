package main

import (
	"math/rand"
	"time"

	"github.com/go-resty/resty/v2"
)

func main() {

	rand.Seed(time.Now().Unix())

	client := NewClient(resty.New())
	agent := NewAgent(client)

	agent.Start()
}
