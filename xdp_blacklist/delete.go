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

func NewDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "delete",
		Short:   "Delete IP from blacklist",
		Example: "xdp_blacklist delete 172.17.0.3",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			blackIP := net.ParseIP(args[0])
			if blackIP == nil {
				log.Println("err ip format: ", args[0])
				return
			}
			log.Printf("will delete ip [%s] from xdp blacklist\n", blackIP)

			// Allow the current process to lock memory for eBPF resources.
			if err := rlimit.RemoveMemlock(); err != nil {
				log.Fatal(err)
			}

			pinPath := path.Join(bpfFSPath, programName)
			_, err := os.Stat(pinPath)
			if err != nil {
				if os.IsNotExist(err) {
					log.Printf("can't find xdp blacklist path in [%s], use xdp_blacklist attach first\n", pinPath)
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
			if err := spec.LoadAndAssign(&objs, &ebpf.CollectionOptions{
				Maps: ebpf.MapOptions{
					PinPath: pinPath,
				},
			}); err != nil {
				panic(err)
			}
			defer objs.Program.Close()
			defer objs.Map.Close()

			err = objs.Map.Delete(ConvertIP2Number(blackIP))
			if err != nil {
				log.Println("delete map err, ", err)
				return
			}
		},
	}
}
