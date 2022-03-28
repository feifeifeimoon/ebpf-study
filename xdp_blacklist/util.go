package main

import (
	"bytes"
	"encoding/binary"
	"github.com/cilium/ebpf"
	"net"
)

const (
	// bpf filesystem
	bpfFSPath = "/sys/fs/bpf"

	programName = "xdp_blacklist"

	xpdObjName = "xdp_blacklist.o"
)

// ConvertIP2Number convert ip to uint32
func ConvertIP2Number(ip net.IP) uint32 {
	var num uint32
	binary.Read(bytes.NewBuffer(ip.To4()), binary.LittleEndian, &num)
	return num
}

// ConvertNum2IP convert uint32 to ip
func ConvertNum2IP(ip uint32) net.IP {
	buffer := bytes.NewBuffer([]byte{})
	err := binary.Write(buffer, binary.LittleEndian, ip)
	if err != nil {
		return nil
	}

	return buffer.Bytes()[:4]
}

type XDPObj struct {
	Program *ebpf.Program `ebpf:"xdp_blacklist_prog"`
	Map     *ebpf.Map     `ebpf:"xdp_blacklist_map"`
}
