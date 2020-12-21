package main

import (
	"fmt"
	"time"
)

func main() {
	q := make(chan struct{})
	for i := 0; i < 10; i++ {
		go te(q)
	}
	close(q)
	time.Sleep(time.Second * 1)
}

func te(q chan struct{}) {
	for {
		select {
		case <-q:
			fmt.Println("Quit")
			return
		}
	}
}
