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
	"strings"
	"sync"
	"testing"
	"time"
)

func createSocketPair(t *testing.T) (net.Listener, net.Conn, net.Conn) {
	// Create a listener.
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Listen failed: %v", err)
	}
	addr := listener.Addr().String()

	// Dial a client, Accept a server.
	wg := sync.WaitGroup{}

	var clientConn net.Conn
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		clientConn, err = net.Dial("tcp", addr)
		if err != nil {
			t.Errorf("Dial failed: %v", err)
		}
	}()

	var serverConn net.Conn
	wg.Add(1)
	go func() {
		defer wg.Done()
		var err error
		serverConn, err = listener.Accept()
		if err != nil {
			t.Errorf("Accept failed: %v", err)
		}
	}()

	wg.Wait()

	return listener, serverConn, clientConn
}

func TestReadTimeout(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	cConnWithTimeout := NewConnWithTimeouts(cConn, 1*time.Millisecond, 1*time.Millisecond)

	c := make(chan error, 1)
	go func() {
		_, err := cConnWithTimeout.Read(make([]byte, 10))
		c <- err
	}()

	select {
	case err := <-c:
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		if !strings.HasSuffix(err.Error(), "i/o timeout") {
			t.Errorf("Expected error timeout, got %s", err)
		}
	case <-time.After(10 * time.Second):
		t.Errorf("Timeout did not happen")
	}
}

func TestWriteTimeout(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	sConnWithTimeout := NewConnWithTimeouts(sConn, 1*time.Millisecond, 1*time.Millisecond)

	c := make(chan error, 1)
	go func() {
		// The timeout will trigger when the buffer is full, so to test this we need to write multiple times.
		for {
			_, err := sConnWithTimeout.Write([]byte("payload"))
			if err != nil {
				c <- err
				return
			}
		}
	}()

	select {
	case err := <-c:
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}

		if !strings.HasSuffix(err.Error(), "i/o timeout") {
			t.Errorf("Expected error timeout, got %s", err)
		}
	case <-time.After(10 * time.Second):
		t.Errorf("Timeout did not happen")
	}
}

func TestNoTimeouts(t *testing.T) {
	listener, sConn, cConn := createSocketPair(t)
	defer func() {
		listener.Close()
		sConn.Close()
		cConn.Close()
	}()

	cConnWithTimeout := NewConnWithTimeouts(cConn, 0, 24*time.Hour)

	c := make(chan error, 1)
	go func() {
		_, err := cConnWithTimeout.Read(make([]byte, 10))
		c <- err
	}()

	select {
	case <-c:
		t.Fatalf("Connection timeout, without a timeout")
	case <-time.After(100 * time.Millisecond):
		// NOOP
	}

	c2 := make(chan error, 1)
	sConnWithTimeout := NewConnWithTimeouts(sConn, 24*time.Hour, 0)
	go func() {
		// This should not fail as there is not timeout on write.
		for {
			_, err := sConnWithTimeout.Write([]byte("payload"))
			if err != nil {
				c2 <- err
				return
			}
		}
	}()
	select {
	case <-c2:
		t.Fatalf("Connection timeout, without a timeout")
	case <-time.After(100 * time.Millisecond):
		// NOOP
	}
}
