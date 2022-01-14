package netutil

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
	"net"
	"time"
)

var _ net.Conn = (*ConnWithTimeouts)(nil)

// A ConnWithTimeouts is a wrapper to net.Conn that allows to set a read and write timeouts.
type ConnWithTimeouts struct {
	net.Conn
	readTimeout  time.Duration
	writeTimeout time.Duration
}

// NewConnWithTimeouts wraps a net.Conn with read and write deadilnes.
func NewConnWithTimeouts(conn net.Conn, readTimeout time.Duration, writeTimeout time.Duration) ConnWithTimeouts {
	return ConnWithTimeouts{Conn: conn, readTimeout: readTimeout, writeTimeout: writeTimeout}
}

// Implementation of the Conn interface.

// Read sets a read deadilne and delegates to conn.Read.
func (c ConnWithTimeouts) Read(b []byte) (int, error) {
	if c.readTimeout == 0 {
		return c.Conn.Read(b)
	}
	if err := c.Conn.SetReadDeadline(time.Now().Add(c.readTimeout)); err != nil {
		return 0, err
	}
	return c.Conn.Read(b)
}

// Write sets a write deadline and delegates to conn.Write
func (c ConnWithTimeouts) Write(b []byte) (int, error) {
	if c.writeTimeout == 0 {
		return c.Conn.Write(b)
	}
	if err := c.Conn.SetWriteDeadline(time.Now().Add(c.writeTimeout)); err != nil {
		return 0, err
	}
	return c.Conn.Write(b)
}

// SetDeadline implements the Conn SetDeadline method.
func (c ConnWithTimeouts) SetDeadline(t time.Time) error {
	panic("can't call SetDeadline for ConnWithTimeouts")
}

// SetReadDeadline implements the Conn SetReadDeadline method.
func (c ConnWithTimeouts) SetReadDeadline(t time.Time) error {
	panic("can't call SetReadDeadline for ConnWithTimeouts")
}

// SetWriteDeadline implements the Conn SetWriteDeadline method.
func (c ConnWithTimeouts) SetWriteDeadline(t time.Time) error {
	panic("can't call SetWriteDeadline for ConnWithTimeouts")
}
