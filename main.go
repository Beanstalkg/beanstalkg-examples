package main

import (
	"fmt"
	"github.com/kr/beanstalk"
	"time"
)

func main() {
	producer(20)
	feedback := make(chan bool)
	go worker(feedback)
	go worker(feedback)

	consumeCount := 0
 	for {
		select {
		case <-feedback:
			consumeCount++;
			fmt.Println("Consumed: ", consumeCount)
		}
	}

}

func producer(work int) {
	c, _ := beanstalk.Dial("tcp", "127.0.0.1:11300")
	for i := 1; i <= work; i++ {
		c.Put([]byte(fmt.Sprintf("WORK %d", i)), uint32(i), 0, 1000 *time.Second)
	}
}

func worker(feedback chan <- bool) {
	c, _ := beanstalk.Dial("tcp", "127.0.0.1:11300")
	for {
		id, body, _ := c.Reserve(5 * time.Second)
		fmt.Println("Completed " + string(body))
		c.Delete(id)
		feedback <- true
	}

}
