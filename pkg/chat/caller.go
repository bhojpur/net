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
	"reflect"
)

type caller struct {
	Func        reflect.Value
	Args        reflect.Type
	ArgsPresent bool
	Out         bool
}

var (
	ErrorCallerNotFunc     = errors.New("f is not function")
	ErrorCallerNot2Args    = errors.New("f should have 1 or 2 args")
	ErrorCallerMaxOneValue = errors.New("f should return not more than one value")
)

/**
Parses function passed by using reflection, and stores its representation
for further call on message or ack
*/
func newCaller(f interface{}) (*caller, error) {
	fVal := reflect.ValueOf(f)
	if fVal.Kind() != reflect.Func {
		return nil, ErrorCallerNotFunc
	}

	fType := fVal.Type()
	if fType.NumOut() > 1 {
		return nil, ErrorCallerMaxOneValue
	}

	curCaller := &caller{
		Func: fVal,
		Out:  fType.NumOut() == 1,
	}
	if fType.NumIn() == 1 {
		curCaller.Args = nil
		curCaller.ArgsPresent = false
	} else if fType.NumIn() == 2 {
		curCaller.Args = fType.In(1)
		curCaller.ArgsPresent = true
	} else {
		return nil, ErrorCallerNot2Args
	}

	return curCaller, nil
}

/**
returns function parameter as it is present in it using reflection
*/
func (c *caller) getArgs() interface{} {
	return reflect.New(c.Args).Interface()
}

/**
calls function with given arguments from its representation using reflection
*/
func (c *caller) callFunc(h *Channel, args interface{}) []reflect.Value {
	//nil is untyped, so use the default empty value of correct type
	if args == nil {
		args = c.getArgs()
	}

	a := []reflect.Value{reflect.ValueOf(h), reflect.ValueOf(args)}
	if !c.ArgsPresent {
		a = a[0:1]
	}

	return c.Func.Call(a)
}
