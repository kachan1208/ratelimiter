# Semi Lazy Rate limiter

It uses idea of Redis that uses 2 different techniques to clean data:
 - Lazy expiring
 - Active random expiring

Tests:
```
    make test
```

Benchmark:
```
    make bench
```

Example of usage:
```Go
package main

import (
	"fmt"
	"time"

	ratelimiter "github.com/kachan1208/rate-limiter"
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
```

Output
```Bash
Error: <nil>
Error: rate limit reached
Error: <nil>
```