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
	"net/http"
	"time"
)

/**
End-point connection for given transport
*/
type Connection interface {
	/**
	Receive one more message, block until received
	*/
	GetMessage() (message string, err error)

	/**
	Send given message, block until sent
	*/
	WriteMessage(message string) error

	/**
	Close current connection
	*/
	Close()

	/**
	Get ping time interval and ping request timeout
	*/
	PingParams() (interval, timeout time.Duration)
}

/**
Connection factory for given transport
*/
type Transport interface {
	/**
	Get client connection
	*/
	Connect(url string) (conn Connection, err error)

	/**
	Handle one server connection
	*/
	HandleConnection(w http.ResponseWriter, r *http.Request) (conn Connection, err error)

	/**
	Serve HTTP request after making connection and events setup
	*/
	Serve(w http.ResponseWriter, r *http.Request)
}
