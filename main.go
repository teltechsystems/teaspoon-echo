package main

import (
	"fmt"
	"github.com/teltechsystems/teaspoon"
	"github.com/teltechsystems/teaspoon/binders"
	"math/rand"
	"os"
	"time"
)

var incrementor = 0

type RandomDataBinder struct {
	binders.ConnectionPool
}

func (b *RandomDataBinder) BroadcastRandomData() {
	connections := b.GetConnections()

	for j := range connections {
		requestID := teaspoon.RequestID{}
		for i := 0; i < 16; i++ {
			requestID[i] = byte(rand.Intn(16))
		}

		randomData := []byte{}
		for i := 0; i < 16; i++ {
			randomData = append(randomData, byte(rand.Intn(57)+65))
		}

		r := &teaspoon.Request{
			OpCode:    teaspoon.OPCODE_BINARY,
			Priority:  5,
			Method:    0,
			Resource:  0,
			RequestID: requestID,
			Payload:   randomData,
		}

		r.WriteTo(connections[j])
	}
}

func NewRandomDataBinder() *RandomDataBinder {
	binder := &RandomDataBinder{}

	go func() {
		for {
			time.Sleep(time.Second * 10)
			binder.BroadcastRandomData()
		}
	}()

	return binder
}

func handler(w teaspoon.ResponseWriter, r *teaspoon.Request) {
	switch r.Resource {
	case 1:
		incrementor += 1

		if incrementor%10 != 0 {
			panic("Destroying connection")
		}
	}

	w.Write(r.Payload)
}

func main() {
	server := &teaspoon.Server{Addr: ":" + os.Getenv("PORT"), Handler: teaspoon.HandlerFunc(handler)}
	server.AddBinder(binders.NewPinger(time.Second * 15))
	server.AddBinder(NewRandomDataBinder())
	fmt.Println(server.ListenAndServe())
}
