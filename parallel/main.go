package main

import (
	"fmt"
	"time"
)

type MyObj struct {
	ID uint64
}

func main() {
	objs := []MyObj{{1}, {2}, {3}, {4}, {5}}

	expectations := make(chan chan *uint64, len(objs))

	for _, obj := range objs {
		expectation := make(chan *uint64, 1)
		expectations <- expectation
		go func(obj MyObj) {
			time.Sleep(150 * time.Millisecond)
			fmt.Println(obj)
			if obj.ID%2 != 0 {
				expectation <- &obj.ID
			} else {
				expectation <- nil
			}
		}(obj)
	}
	close(expectations)

	results := []uint64{}
	for expectation := range expectations {
		if err := <-expectation; err != nil {
			results = append(results, *err)
		}
	}

	fmt.Println(results)
}
