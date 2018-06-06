//    Copyright 2018 <xu shuai <dyxushuai@gmail.com>>
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package grpcdump

import (
	"fmt"
	"io"
	"os"
)

var DefaultPrinter = NewStdoutPrinter()

const defaultBufferLen = 4096

// Printer buffered print
type Printer struct {
	buffer chan string
	output io.WriteCloser
}

// NewStdoutPrinter print to os.Stdout
func NewStdoutPrinter() *Printer {
	return NewPrinter(os.Stdout)
}

// NewPrinter print ot io.WriteCloser
func NewPrinter(output io.WriteCloser) *Printer {
	p := &Printer{
		buffer: make(chan string, defaultBufferLen),
		output: output,
	}
	go p.process()
	return p
}

// Printlnf print formated messages
func (p *Printer) Printlnf(format string, args ...interface{}) {
	msg := fmt.Sprintf(format+"\r\n", args...)
	p.buffer <- msg
}

func (p *Printer) process() {
	for msg := range p.buffer {
		p.output.Write([]byte(msg))
	}
}

// Close the printer
func (p *Printer) Close() {
	close(p.buffer)
}
