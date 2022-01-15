//go:build client
// +build client

package examples

import (
	"log"
	"runtime"
	"time"

	"github.com/bhojpur/net/pkg/chat"
	"github.com/bhojpur/net/pkg/transport"
)

type Channel struct {
	Channel string `json:"channel"`
}

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func sendJoin(c *chat.Client) {
	log.Println("acking /join")
	result, err := c.Ack("/join", Channel{"main"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("ack result to /join: ", result)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := chat.Dial(
		chat.GetUrl("localhost", 3811, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/message", func(h *chat.Channel, args Message) {
		log.Println("--- got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(chat.OnDisconnection, func(h *chat.Channel) {
		log.Fatal("disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(chat.OnConnection, func(h *chat.Channel) {
		log.Println("connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)
	go sendJoin(c)

	time.Sleep(60 * time.Second)
	c.Close()

	log.Println(" [x] complete")
}
