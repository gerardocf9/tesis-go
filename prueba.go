package main

import (
	"fmt"
	"time"
)

func main() {
	t := time.Now()
	for {
		fmt.Println(time.Since(t))
		u, _ := time.ParseDuration("15s")
		if time.Since(t) > u {
			fmt.Println("pasaron 15 s")
			break
		}
		fmt.Println(time.Since(t))
	}
}
