package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"github.com/cilium/ebpf"
	"net"
	"os"
	"os/signal"
	"syscall"
)

const (
	// bpf filesystem
	bpfFSPath = "/sys/fs/bpf"

	programName = "xdp_blacklist"

	xpdObjName = "xdp_blacklist.o"
)

var onlyOneSignalHandler = make(chan struct{})
var shutdownHandler chan os.Signal

// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
// Only one of SetupSignalContext and SetupSignalHandler should be called, and only can
// be called once.
func SetupSignalHandler() <-chan struct{} {
	return SetupSignalContext().Done()
}

// SetupSignalContext is same as SetupSignalHandler, but a context.Context is returned.
// Only one of SetupSignalContext and SetupSignalHandler should be called, and only can
// be called once.
func SetupSignalContext() context.Context {
	close(onlyOneSignalHandler) // panics when called twice

	shutdownHandler = make(chan os.Signal, 2)

	ctx, cancel := context.WithCancel(context.Background())
	signal.Notify(shutdownHandler, []os.Signal{os.Interrupt, syscall.SIGTERM}...)
	go func() {
		<-shutdownHandler
		cancel()
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx
}

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
