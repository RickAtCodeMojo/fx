//
//  Pubsub envelope subscriber
//

package main

import (
	"fmt"

	zmq "github.com/pebbe/zmq4"
)

func main() {
	context, _ := zmq.NewContext()
	defer context.Term()

	subscriber, _ := context.NewSocket(zmq.SUB)
	defer subscriber.Close()
	subscriber.Connect("tcp://localhost:5563")
	subscriber.SetSubscribe("USD")

	for {
		// address, _ := subscriber.Recv(0)
		// content, _ := subscriber.Recv(0)
		// print("[" + string(address) + "] " + string(content) + "\n")
		content, _ := subscriber.Recv(0)
		fmt.Println(content)
		fmt.Println("====================================================================")
	}
}
