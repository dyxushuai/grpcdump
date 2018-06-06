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
	"sync"
	"time"

	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"github.com/jhump/protoreflect/desc"

	"github.com/google/gopacket"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
)

const (
	contentTypeHeader = "content-type"
	pathHeader        = ":path"
)

// gprcMediaTypes the content-type of http2 header for grpc
var grpcMediaTypes = []string{
	"application/grpc",
	"application/grpc+proto",
	"application/grpc+json",
}

// IsGrpc assert grpc protocol by content-type
func isGrpc(contentType string) bool {
	for _, gmt := range grpcMediaTypes {
		if gmt == contentType {
			return true
		}
	}
	return false
}

type GrpcStreamFactory struct {
	printer         *Printer
	hostIP          string
	hostPort        uint16
	protoParser     *protoFileDescs
	rw              sync.RWMutex
	currentRequests map[uint32]*grpcParser
}

func NewGrpcStreamFactory(hostIP string, hostPort uint16, printer *Printer, protoFile string) (*GrpcStreamFactory, error) {
	if hostIP == "" || hostPort == 0 {
		return nil, fmt.Errorf("miss hostIP: %s or hostPort: %d", hostIP, hostPort)
	}
	protoParser, err := protoParse(nil, protoFile)
	if err != nil {
		return nil, err
	}
	return &GrpcStreamFactory{
		printer:         printer,
		hostIP:          hostIP,
		hostPort:        hostPort,
		protoParser:     protoParser,
		currentRequests: map[uint32]*grpcParser{},
	}, nil
}

func (s *GrpcStreamFactory) setCurrentRequest(streamID uint32, request *grpcParser) {
	s.rw.Lock()
	s.currentRequests[streamID] = request
	s.rw.Unlock()
}

func (s *GrpcStreamFactory) getCurrentRequest(streamID uint32) (*grpcParser, bool) {
	s.rw.RLock()
	defer s.rw.RUnlock()
	request, ok := s.currentRequests[streamID]
	return request, ok
}

func (s *GrpcStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	r := tcpreader.NewReaderStream()
	g := &grpcParser{
		factory:     s,
		protoParser: s.protoParser,
		hostIP:      s.hostIP,
		hostPort:    s.hostPort,
		net:         net,
		printer:     s.printer,
		transprot:   transport,
		framer:      http2.NewFramer(nil, newSkipClientPrefaceReader(&r)),
	}
	g.framer.ReadMetaHeaders = hpack.NewDecoder(uint32(4<<10), nil)
	g.framer.ReadMetaHeaders.SetEmitEnabled(false)
	go g.HandleFrames()
	return &r
}

// grpcParser parse the grpc protocol from byte stream
type grpcParser struct {
	factory        *GrpcStreamFactory
	protoParser    *protoFileDescs
	printer        *Printer
	framer         *http2.Framer
	hostIP         string
	hostPort       uint16
	net, transprot gopacket.Flow
	//
	streamID      uint32
	isGrpc        bool
	isRequest     bool
	input, output *desc.MessageDescriptor
}

func (g *grpcParser) setStreamID(streamID uint32) {
	if g.streamID == 0 {
		g.streamID = streamID
	}
}

func (g *grpcParser) myResponse(other *grpcParser) bool {
	if g.streamID != other.streamID {
		return false
	}
	src, dst := g.flow()
	otherSrc, otherDst := other.flow()
	return src == otherDst && dst == otherSrc
}

func (g *grpcParser) flow() (string, string) {
	return fmt.Sprintf("%s:%s", g.net.Src(), g.transprot.Src()), fmt.Sprintf("%s:%s", g.net.Dst(), g.transprot.Dst())
}

func (g *grpcParser) onNewHeaderField(f hpack.HeaderField) {
	if f.Name == pathHeader {
		g.isRequest = true
		var err error
		g.input, g.output, err = g.protoParser.findMehodSignature(f.Value)
		if err != nil {
			g.handleError(err)
		}
	}
	if f.Sensitive {
		g.printlnf("    %s = %q (SENSITIVE)", f.Name, f.Value)
	}
	g.printlnf("    %s = %q", f.Name, f.Value)
}

func (g *grpcParser) handleError(err error) {
	g.printlnf("ERROR: %v", err)
}

func (g *grpcParser) handleHeaders(f *http2.MetaHeadersFrame) bool {
	if !g.isGrpc {
		for _, hf := range f.Fields {
			if hf.Name == contentTypeHeader {
				g.isGrpc = isGrpc(hf.Value)
			}
			if hf.Name == pathHeader {
				g.isRequest = true
				var err error
				g.input, g.output, err = g.protoParser.findMehodSignature(hf.Value)
				if err != nil {
					g.handleError(err)
				}
			}
		}
	}
	if !g.isGrpc {
		return false
	}
	g.setStreamID(f.StreamID)
	now := time.Now().Format(time.RFC3339Nano)
	g.printlnf("")
	src, dst := g.flow()
	if g.isRequest {
		g.printlnf("REQUEST(STREAM=%d) > %s: %s ---> %s", f.StreamID, now, src, dst)
		g.factory.setCurrentRequest(f.StreamID, g)
	} else {
		if requst, ok := g.factory.getCurrentRequest(f.StreamID); ok {
			if requst.myResponse(g) {
				g.printlnf("RESPONSE(STREAM=%d) > %s: %s <--- %s", f.StreamID, now, dst, src)
				g.input = requst.input
				g.output = requst.output
			}
		}
	}
	g.printlnf("  HEADERS:")
	for _, hf := range f.Fields {
		if hf.Sensitive {
			g.printlnf("    %s = %q (SENSITIVE)", hf.Name, hf.Value)
		}
		g.printlnf("    %s = %q", hf.Name, hf.Value)
	}

	return true
}

func (g *grpcParser) handleBody(f *http2.DataFrame) {
	if !g.isGrpc {
		return
	}
	g.printlnf("  BODY:")
	if g.isRequest {
		if g.input != nil {
			pretty, err := g.protoParser.pretty(f.Data()[5:], g.input)
			if err != nil {
				g.handleError(err)
			} else {
				g.printlnf("    %s", pretty)
				return
			}
		}
	} else {
		if g.output != nil {
			pretty, err := g.protoParser.pretty(f.Data()[5:], g.output)
			if err != nil {
				g.handleError(err)
			} else {
				g.printlnf("    %s", pretty)
				return
			}
		}
	}
	g.printlnf("    %q", f.Data()[5:])
}

func (g *grpcParser) printlnf(format string, args ...interface{}) {
	g.printer.Printlnf(format, args...)
}

func (g *grpcParser) HandleFrames() {
	for {
		f, err := g.framer.ReadFrame()
		if err != nil {
			fmt.Printf("HandleFrames error: %v\n", err)
			return
		}
		switch f := f.(type) {
		case *http2.MetaHeadersFrame:
			g.handleHeaders(f)
		case *http2.DataFrame:
			g.handleBody(f)
		default:
			// ignore
		}
	}
}
