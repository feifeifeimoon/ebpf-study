package main

import (
	"github.com/spf13/cobra"
	"github.com/vishvananda/netlink"
	"log"
)

func NewDetachCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "detach",
		Short:   "Detach xdp from network device",
		Example: "xdp_blacklist detach eth0",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			ethName := args[0]
			log.Printf("will detach xdp program from dev[%s]\n", ethName)

			link, err := netlink.LinkByName(ethName)
			if err != nil {
				log.Printf("get link by name err, %+v ", err)
				return
			}

			_ = netlink.LinkSetXdpFd(link, -1)
		},
	}
}
