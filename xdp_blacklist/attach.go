package main

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	"log"
	"os"
	"path"
)

func NewAttachCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "attach",
		Short:   "Attach xdp to network device",
		Example: "xdp_blacklist attach eth0",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ethName := args[0]
			log.Printf("will attach xdp program on dev[%s]\n", ethName)
			//ctx := SetupSignalContext()

			// Allow the current process to lock memory for eBPF resources.
			if err := rlimit.RemoveMemlock(); err != nil {
				log.Fatal(err)
			}

			// Create pin path
			pinPath := path.Join(bpfFSPath, programName)
			if err := os.MkdirAll(pinPath, os.ModePerm); err != nil {
				log.Fatalf("failed to create bpf fs subpath: %+v", err)
			}

			spec, err := ebpf.LoadCollectionSpec("xdp_blacklist.o")
			if err != nil {
				log.Panic("load ebpf err, ", err)
			}

			var objs XDPObj
			if err := spec.LoadAndAssign(&objs, &ebpf.CollectionOptions{
				Maps: ebpf.MapOptions{
					// Pin the map to the BPF filesystem and configure the
					// library to automatically re-write it in the BPF
					// program so it can be re-used if it already exists or
					// create it if not
					PinPath: pinPath,
				},
			}); err != nil {
				panic(err)
			}
			defer objs.Program.Close()
			defer objs.Map.Close()

			link, err := netlink.LinkByName(ethName)
			if err != nil {
				log.Printf("get link by name err, %+v ", err)
				return
			}

			err = netlink.LinkSetXdpFd(link, objs.Program.FD())
			if err != nil {
				log.Printf("set xdp to dev[%s] err, %+v", ethName, err)
				return
			}
			//defer netlink.LinkSetXdpFd(link, -1)
		},
	}
}
