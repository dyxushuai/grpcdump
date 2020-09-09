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
	"log"
	"time"

	"github.com/dyxushuai/grpcdump"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/tcpassembly"
)

var device = flag.String("d", "eth0", "Interface to get packets from")
var logAllPackets = flag.Bool("v", false, "Logs every packet in great detail")
var filterIP = flag.String("ip", "", "Filter by ip, if either source or target ip is matched, the packet will be processed")
var filterPort = flag.Uint("port", 0, "Filter by port, if either source or target port is matched, the packet will be processed.")

func main() {
	var protoFile grpcdump.ArrayFlags
	var protoPath grpcdump.ArrayFlags
	flag.Var(&protoPath, "proto_path", "Specify the directory in which to search for imports")
	flag.Var(&protoFile, "proto", "Protobuf spec file")
	flag.Parse()
	log.Printf("Starting capture on device %q", *device)
	printer := grpcdump.NewStdoutPrinter()
	sf, err := grpcdump.NewGrpcStreamFactory(*filterIP, uint16(*filterPort), printer, protoFile.ParseDir(".proto"), protoPath)
	//sf, err := grpcdump.NewGrpcStreamFactory(*filterIP, uint16(*filterPort), printer, protoFile, protoPath)
	if err != nil {
		log.Fatal(err)
	}
	streamPool := tcpassembly.NewStreamPool(sf)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Println("reading in packets")
	packets, err := grpcdump.CaptureSingleDevice(*device, *filterIP, uint16(*filterPort))
	if err != nil {
		log.Fatal(err)
	}
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			// A nil packet indicates the end of a pcap file.
			if packet == nil {
				return
			}
			if *logAllPackets {
				log.Println(packet)
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				log.Println("Unusable packet")
				continue
			}
			tcp := packet.TransportLayer().(*layers.TCP)
			assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
		case <-ticker:
			// Every minute, flush connections that haven't seen activity in the past 2 minutes.
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		}
	}
}
