/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

package main

import (
	"flag"
	"log"
	"net"

	pb "github.com/dyxushuai/grpcdump/grpc_example/helloworld/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", ":8085", "Grpc listen address")

type i18 struct {
	value     string
	languages map[string]string
}

var defaultI18s = []i18{
	{
		value: "world",
		languages: map[string]string{
			"zh": "你好",
			"en": "world",
		},
	},
}

func lookupI18(value, language string) string {
	for _, i := range defaultI18s {
		if i.value == value {
			for l, result := range i.languages {
				if l == language {
					return result
				}
			}
		}
	}
	return value
}

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello, " + lookupI18(in.Name, in.I18)}, nil
}

func main() {
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
