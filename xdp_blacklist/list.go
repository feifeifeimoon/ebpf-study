package main

import (
	"fmt"
	"github.com/cilium/ebpf"
	"github.com/cilium/ebpf/rlimit"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Show ip in blacklist",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// Allow the current process to lock memory for eBPF resources.
			if err := rlimit.RemoveMemlock(); err != nil {
				log.Fatal(err)
			}

			// Create pin path
			pinPath := path.Join(bpfFSPath, programName)
			_, err := os.Stat(pinPath)
			if err != nil {
				if os.IsExist(err) {
					log.Printf("can't not find xdp blacklist path in %s, use xdp_blacklist attach first\n", pinPath)
					return
				}
				log.Println("stat err ", err)
				return
			}

			spec, err := ebpf.LoadCollectionSpec(xpdObjName)
			if err != nil {
				log.Panic("load ebpf err, ", err)
			}

			var objs XDPObj
			if err := spec.LoadAndAssign(&objs, &ebpf.CollectionOptions{Maps: ebpf.MapOptions{PinPath: pinPath}}); err != nil {
				panic(err)
			}

			defer objs.Program.Close()
			defer objs.Map.Close()

			key := uint32(0)
			value := uint64(0)

			fmt.Println("IP\t\tHit")
			iter := objs.Map.Iterate()
			for iter.Next(&key, &value) {
				fmt.Println(ConvertNum2IP(key), value)
			}
		},
	}
}
