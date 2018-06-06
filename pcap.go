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
