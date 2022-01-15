package chat

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
	"sync"
)

var (
	ErrorWaiterNotFound = errors.New("waiter not found")
)

/**
Processes functions that require answers, also known as acknowledge or ack
*/
type ackProcessor struct {
	counter     int
	counterLock sync.Mutex

	resultWaiters     map[int](chan string)
	resultWaitersLock sync.RWMutex
}

/**
get next id of ack call
*/
func (a *ackProcessor) getNextId() int {
	a.counterLock.Lock()
	defer a.counterLock.Unlock()

	a.counter++
	return a.counter
}

/**
Just before the ack function called, the waiter should be added
to wait and receive response to ack call
*/
func (a *ackProcessor) addWaiter(id int, w chan string) {
	a.resultWaitersLock.Lock()
	a.resultWaiters[id] = w
	a.resultWaitersLock.Unlock()
}

/**
removes waiter that is unnecessary anymore
*/
func (a *ackProcessor) removeWaiter(id int) {
	a.resultWaitersLock.Lock()
	delete(a.resultWaiters, id)
	a.resultWaitersLock.Unlock()
}

/**
check if waiter with given ack id is exists, and returns it
*/
func (a *ackProcessor) getWaiter(id int) (chan string, error) {
	a.resultWaitersLock.RLock()
	defer a.resultWaitersLock.RUnlock()

	if waiter, ok := a.resultWaiters[id]; ok {
		return waiter, nil
	}
	return nil, ErrorWaiterNotFound
}
