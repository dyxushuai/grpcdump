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

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dyxushuai/grpcdump/grpc_example/helloworld/helloworld"
	"github.com/gogo/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
)

var device = flag.String("d", "eth0", "Interface to get packets from")
var logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")
var filterIP = flag.String("ip", "", "Filter by ip, if either source or target ip is matched, the packet will be processed")
var filterPort = flag.Uint("port", 0, "Filter by port, if either source or target port is matched, the packet will be processed.")
var protoFile = flag.String("proto", "", "Protobuf spec file")

func main() {
	flag.Parse()
	h := &helloworld.HelloRequest{Name: "hello"}
	data, err := proto.Marshal(h)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(data))
	p := &protoparse.Parser{}
	descs, err := p.ParseFiles(*protoFile)
	if err != nil {
		panic(err)
	}
	input, _, err := findMehodSignature(descs, "/helloworld.Greeter/SayHello")
	if err != nil {
		panic(err)
	}
	dmsg := dynamic.NewMessage(input)
	err = dmsg.Unmarshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Println(dmsg.String())
}

func findMehodSignature(descs []*desc.FileDescriptor, path string) (*desc.MessageDescriptor, *desc.MessageDescriptor, error) {
	strs := strings.Split(path, "/")
	if len(strs) != 3 {
		return nil, nil, fmt.Errorf("error path format: %s", path)
	}
	for _, desc := range descs {
		srvDesc := desc.FindService(strs[1])
		if srvDesc == nil {
			return nil, nil, fmt.Errorf("service name not found: %s", strs[1])
		}
		mtdDesc := srvDesc.FindMethodByName(strs[2])
		if mtdDesc == nil {
			return nil, nil, fmt.Errorf("method name not found: %s", strs[2])
		}
		return mtdDesc.GetInputType(), mtdDesc.GetOutputType(), nil
	}
	return nil, nil, fmt.Errorf("message not found: %s", path)
}
