package pipeline

import (
	"fmt"
	"time"
)

// Sorry Frank.
type Timestamp int

type Message interface {
}

type Event struct {
	Time Timestamp
	Msg  Message
}

func now() int {
	return int(time.Now().UnixNano())
}

func Logger(in chan Event) {
	for {
		fmt.Println(<-in)
	}
}
