// Copyright 2012 Google, Inc. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

// This binary provides sample code for using the gopacket TCP assembler and TCP
// stream reader.  It reads packets off the wire and reconstructs HTTP requests
// it sees, logging them.
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
var protoFile = flag.String("proto", "", "Protobuf spec file")

func main() {
	flag.Parse()
	log.Printf("Starting capture on device %q", *device)
	printer := grpcdump.NewStdoutPrinter()
	sf, err := grpcdump.NewGrpcStreamFactory(*filterIP, uint16(*filterPort), printer, *protoFile)
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
