package transport

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	upgradeFailed = "Upgrade failed: "

	WsDefaultPingInterval   = 30 * time.Second
	WsDefaultPingTimeout    = 60 * time.Second
	WsDefaultReceiveTimeout = 60 * time.Second
	WsDefaultSendTimeout    = 60 * time.Second
	WsDefaultBufferSize     = 1024 * 32
)

var (
	ErrorBinaryMessage     = errors.New("binary messages are not supported")
	ErrorBadBuffer         = errors.New("buffer error")
	ErrorPacketWrong       = errors.New("wrong packet type error")
	ErrorMethodNotAllowed  = errors.New("method not allowed")
	ErrorHttpUpgradeFailed = errors.New("the HTTP upgrade failed")
)

type WebsocketConnection struct {
	socket    *websocket.Conn
	transport *WebsocketTransport
}

func (wsc *WebsocketConnection) GetMessage() (message string, err error) {
	wsc.socket.SetReadDeadline(time.Now().Add(wsc.transport.ReceiveTimeout))
	msgType, reader, err := wsc.socket.NextReader()
	if err != nil {
		return "", err
	}

	//support only text messages exchange
	if msgType != websocket.TextMessage {
		return "", ErrorBinaryMessage
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", ErrorBadBuffer
	}
	text := string(data)

	//empty messages are not allowed
	if len(text) == 0 {
		return "", ErrorPacketWrong
	}

	return text, nil
}

func (wsc *WebsocketConnection) WriteMessage(message string) error {
	wsc.socket.SetWriteDeadline(time.Now().Add(wsc.transport.SendTimeout))
	writer, err := wsc.socket.NextWriter(websocket.TextMessage)
	if err != nil {
		return err
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}
	return nil
}

func (wsc *WebsocketConnection) Close() {
	wsc.socket.Close()
}

func (wsc *WebsocketConnection) PingParams() (interval, timeout time.Duration) {
	return wsc.transport.PingInterval, wsc.transport.PingTimeout
}

type WebsocketTransport struct {
	PingInterval   time.Duration
	PingTimeout    time.Duration
	ReceiveTimeout time.Duration
	SendTimeout    time.Duration

	BufferSize int

	RequestHeader http.Header
}

func (wst *WebsocketTransport) Connect(url string) (conn Connection, err error) {
	dialer := websocket.Dialer{}
	socket, _, err := dialer.Dial(url, wst.RequestHeader)
	if err != nil {
		return nil, err
	}

	return &WebsocketConnection{socket, wst}, nil
}

func (wst *WebsocketTransport) HandleConnection(
	w http.ResponseWriter, r *http.Request) (conn Connection, err error) {

	if r.Method != "GET" {
		http.Error(w, upgradeFailed+ErrorMethodNotAllowed.Error(), 503)
		return nil, ErrorMethodNotAllowed
	}

	socket, err := websocket.Upgrade(w, r, nil, wst.BufferSize, wst.BufferSize)
	if err != nil {
		http.Error(w, upgradeFailed+err.Error(), 503)
		return nil, ErrorHttpUpgradeFailed
	}

	return &WebsocketConnection{socket, wst}, nil
}

/**
Websocket connection do not require any additional processing
*/
func (wst *WebsocketTransport) Serve(w http.ResponseWriter, r *http.Request) {}

/**
Returns websocket connection with default params
*/
func GetDefaultWebsocketTransport() *WebsocketTransport {
	return &WebsocketTransport{
		PingInterval:   WsDefaultPingInterval,
		PingTimeout:    WsDefaultPingTimeout,
		ReceiveTimeout: WsDefaultReceiveTimeout,
		SendTimeout:    WsDefaultSendTimeout,
		BufferSize:     WsDefaultBufferSize,
	}
}
