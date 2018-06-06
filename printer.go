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
