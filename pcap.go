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
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func listenOneSource(handle *pcap.Handle) chan gopacket.Packet {
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	packets := packetSource.Packets()
	return packets
}

func setBPFFilter(handle *pcap.Handle, filterIP string, filterPort uint16) error {
	var bpfFilter = "tcp"
	if filterPort != 0 {
		bpfFilter += " and port " + strconv.Itoa(int(filterPort))
	}
	if filterIP != "" {
		bpfFilter += " and ip host " + filterIP
	}
	return handle.SetBPFFilter(bpfFilter)
}

// CaptureSingleDevice capture the packet by given network device and ip, port filter
func CaptureSingleDevice(device string, filterIP string, filterPort uint16) (<-chan gopacket.Packet, error) {
	handle, err := pcap.OpenLive(device, 65536, false, pcap.BlockForever)
	if err != nil {
		return nil, err
	}

	err = setBPFFilter(handle, filterIP, filterPort)
	if err != nil {
		return nil, err
	}
	return listenOneSource(handle), nil
}
