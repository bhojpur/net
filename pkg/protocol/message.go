package protocol

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

const (
	/**
	Message with connection options
	*/
	MessageTypeOpen = iota
	/**
	Close connection and destroy all handle routines
	*/
	MessageTypeClose = iota
	/**
	Ping request message
	*/
	MessageTypePing = iota
	/**
	Pong response message
	*/
	MessageTypePong = iota
	/**
	Empty message
	*/
	MessageTypeEmpty = iota
	/**
	Emit request, no response
	*/
	MessageTypeEmit = iota
	/**
	Emit request, wait for response (ack)
	*/
	MessageTypeAckRequest = iota
	/**
	ack response
	*/
	MessageTypeAckResponse = iota
)

type Message struct {
	Type   int
	AckId  int
	Method string
	Args   string
	Source string
}
