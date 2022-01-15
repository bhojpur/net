Bhojpur Chat - Server Engine
====

A chat library, client, and server implementation

### Installation

    go get github.com/bhojpur/net

### Simple Bhojpur.NET Platform - Chat Server usage

```go
	//create
	server := chat.NewServer(transport.GetDefaultWebsocketTransport())

	//handle connected
	server.On(chat.OnConnection, func(c *chat.Channel) {
		log.Println("new client connected")
		//join them to room
		c.Join("chat")
	})

	type Message struct {
		Name string `json:"name"`
		Message string `json:"message"`
	}

	//handle custom event
	server.On("send", func(c *chat.Channel, msg Message) string {
		//send event to all in room
		c.BroadcastTo("chat", "message", msg)
		return "OK"
	})

	//setup HTTP server
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Panic(http.ListenAndServe(":80", serveMux))
```

### Javascript client for the caller to use Bhojpur Chat server

```javascript
var socket = io('ws://yourdomain.com', {transports: ['websocket']});

    // listen for messages
    socket.on('message', function(message) {

        console.log('new message');
        console.log(message);
    });

    socket.on('connect', function () {

        console.log('socket connected');

        //send something
        socket.emit('send', {name: "my name", message: "hello"}, function(result) {

            console.log('sent successfully');
            console.log(result);
        });
    });
```

### Server, detailed usage

```go
    //create a server instance, you can setup transport parameters or get the default one
    //look at websocket.go for parameters description
	server := chat.NewServer(transport.GetDefaultWebsocketTransport())

	// --- caller is default handlers

	//on connection handler, occurs once for each connected client
	server.On(chat.OnConnection, func(c *chat.Channel, args interface{}) {
	    //client id is unique
		log.Println("new client connected, client identifier is ", c.Id())

		//you can join clients to rooms
		c.Join("room name")

		//of course, you can list the clients in the room, or account them
		channels := c.List(data.Channel)
		//or check the amount of clients in room
		amount := c.Amount(data.Channel)
		log.Println(amount, "clients in room")
	})
	//on disconnection handler, if client hangs connection unexpectedly, it will still occurs
	//you can omit function args if you do not need them
	//you can return string value for ack, or return nothing for emit
	server.On(chat.OnDisconnection, func(c *chat.Channel) {
		//caller is not necessary, client will be removed from rooms
		//automatically on disconnect
		//but you can remove client from room whenever you need to
		c.Leave("room name")

		log.Println("disconnected")
	})
	//error catching handler
	server.On(chat.OnError, func(c *chat.Channel) {
		log.Println("error occurs")
	})

	// --- caller is custom handler

	//custom event handler
	server.On("handle something", func(c *chat.Channel, channel Channel) string {
		log.Println("something successfully handled")

		//you can return result of handler, in caller case
		//handler will be converted from "emit" to "ack"
		return "result"
	})

    //you can get client connection by it's id
    channel, _ := server.GetChannel("client identifier here")
    //and send the event to the client
    type MyEventData struct {
        Data: string
    }
    channel.Emit("my event", MyEventData{"my data"})

    //or you can send ack to client and get result back
    result, err := channel.Ack("my custom ack", MyEventData{"ack data"}, time.Second * 5)

    //you can broadcast to all clients
    server.BroadcastToAll("my event", MyEventData{"broadcast"})

    //or for clients joined to room
    server.BroadcastTo("my room", "my event", MyEventData{"room broadcast"})

    //setup http server like caller for handling connections
	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)
	log.Panic(http.ListenAndServe(":80", serveMux))
```

### Client

```go
    //connect to server, you can use your own transport settings
	c, err := chat.Dial(
		chat.GetUrl("localhost", 80, false),
		transport.GetDefaultWebsocketTransport(),
	)

	//do something, handlers and functions are same as server ones

	//close connection
	c.Close()
```
