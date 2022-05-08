package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	maxNum := 100
	rand.Seed(time.Now().UnixNano())
	ans := rand.Intn(maxNum)
	for {
		var x int
		_, err := fmt.Scanf("%d", &x)
		if err != nil {
			fmt.Println("input invalid, please try again1", err)
			continue
		}

		if err != nil {
			fmt.Println("input invalid, please try again", err)
			continue
		}
		if x > ans {
			fmt.Println("too big!")
		} else if x < ans {
			fmt.Println("too small!")
		} else {
			fmt.Println("Right!")
			break
		}
	}
}
