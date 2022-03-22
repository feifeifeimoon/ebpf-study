package main

import (
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/spf13/cobra"
	"log"
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
			blackIP := args[0]
			log.Printf("will add ip[%s] to xdp blacklist\n", blackIP)
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

			spec, err := ebpf.LoadCollectionSpec(xpdObjName)
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

			log.Println("add tp map", ConvertIP2Number(blackIP))
			err = objs.Map.Put(ConvertIP2Number(blackIP), uint64(0))
			if err != nil {
				log.Panic("update map err, ", err)
				return
			}
		},
	}
}
