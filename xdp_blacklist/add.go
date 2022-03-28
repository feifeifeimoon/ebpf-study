package main

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/spf13/cobra"
	"log"
	"net"
	"os"
	"path"
)

func NewAddCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "add",
		Short:   "Appends a new IP address to the blacklist",
		Example: "xdp_blacklist add 172.17.0.3",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			blackIP := net.ParseIP(args[0])
			if blackIP == nil {
				log.Println("err ip format: ", args[0])
				return
			}
			log.Printf("will add ip [%s] to xdp blacklist\n", blackIP)

			// Allow the current process to lock memory for eBPF resources.
			if err := rlimit.RemoveMemlock(); err != nil {
				log.Fatal(err)
			}

			// Create pin path
			pinPath := path.Join(bpfFSPath, programName)
			if err := os.MkdirAll(pinPath, os.ModePerm); err != nil {
				log.Fatalf("failed to create bpf fs subpath: %+v", err)
			}

			spec, err := ebpf.LoadCollectionSpec(xpdObjName)
			if err != nil {
				log.Panic("load ebpf err, ", err)
			}

			var objs XDPObj
			if err := spec.LoadAndAssign(&objs, &ebpf.CollectionOptions{
				Maps: ebpf.MapOptions{
					PinPath: pinPath,
				},
			}); err != nil {
				panic(err)
			}
			defer objs.Program.Close()
			defer objs.Map.Close()

			err = objs.Map.Put(ConvertIP2Number(blackIP), uint64(0))
			if err != nil {
				log.Println("put map err, ", err)
				return
			}
		},
	}
}
