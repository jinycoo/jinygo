// Copyright (c) 2017 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package multierr_test

import (
	"errors"
	"fmt"
)

func ExampleCombine() {
	err := Combine(
		errors.New("call 1 failed"),
		nil, // successful request
		errors.New("call 3 failed"),
		nil, // successful request
		errors.New("call 5 failed"),
	)
	fmt.Printf("%+v", err)
	// Output:
	// the following errors occurred:
	//  -  call 1 failed
	//  -  call 3 failed
	//  -  call 5 failed
}

func ExampleAppend() {
	var err error
	err = Append(err, errors.New("call 1 failed"))
	err = Append(err, errors.New("call 2 failed"))
	fmt.Println(err)
	// Output:
	// call 1 failed; call 2 failed
}

func ExampleErrors() {
	err := Combine(
		nil, // successful request
		errors.New("call 2 failed"),
		errors.New("call 3 failed"),
	)
	err = Append(err, nil) // successful request
	err = Append(err, errors.New("call 5 failed"))

	errors := Errors(err)
	for _, err := range errors {
		fmt.Println(err)
	}
	// Output:
	// call 2 failed
	// call 3 failed
	// call 5 failed
}

func ExampleAppendInto() {
	var err error

	if AppendInto(&err, errors.New("foo")) {
		fmt.Println("call 1 failed")
	}

	if AppendInto(&err, nil) {
		fmt.Println("call 2 failed")
	}

	if AppendInto(&err, errors.New("baz")) {
		fmt.Println("call 3 failed")
	}

	fmt.Println(err)
	// Output:
	// call 1 failed
	// call 3 failed
	// foo; baz
}
