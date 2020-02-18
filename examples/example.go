package main

import (
	"fmt"
	"time"

	"github.com/kachan1208/ratelimiter"
)

var (
	rateLimiter *ratelimiter.SemiLazyLimiter
)

func init() {
	rateLimiter = ratelimiter.NewSemiLazyLimiter()
}

func main() {
	fmt.Println("Error:", limitMe("1"))
	fmt.Println("Error:", limitMe("1"))
	fmt.Println("Error:", limitMe("2"))
}

func limitMe(key string) error {
	//allow to execute one time per second
	err := rateLimiter.Limit(key, time.Second, 1)
	if err != nil {
		return err
	}

	/*
		your code
		....
	*/

	return nil
}
