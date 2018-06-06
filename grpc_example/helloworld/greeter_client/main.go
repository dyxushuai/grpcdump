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

package main

import (
	"flag"
	"log"
	"time"

	pb "github.com/dyxushuai/grpcdump/grpc_example/helloworld/helloworld"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var addr = flag.String("addr", "localhost:8085", "Grpc server address")

func main() {
	flag.Parse()
	for {
		// Set up a connection to the server.
		conn, err := grpc.Dial(*addr, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		c := pb.NewGreeterClient(conn)

		r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: "world", I18: "zh"})
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		log.Printf("Greeting: %s", r.Message)
		time.Sleep(time.Second)
	}
}
