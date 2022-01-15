//go:build server
// +build server

package examples

import (
	"log"
	"net/http"
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

func main() {
	server := chat.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(chat.OnConnection, func(c *chat.Channel) {
		log.Println("Connected")

		c.Emit("/message", Message{10, "main", "using emit"})

		c.Join("test")
		c.BroadcastTo("test", "/message", Message{10, "main", "using broadcast"})
	})
	server.On(chat.OnDisconnection, func(c *chat.Channel) {
		log.Println("Disconnected")
	})

	server.On("/join", func(c *chat.Channel, channel Channel) string {
		time.Sleep(2 * time.Second)
		log.Println("Client joined to ", channel.Channel)
		return "joined to " + channel.Channel
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Println("Starting Bhojpur Chat server...")
	log.Panic(http.ListenAndServe(":3811", serveMux))
}
