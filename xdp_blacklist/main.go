package main

import (
	"fmt"
	"github.com/spf13/cobra"
)

func main() {

	rootCmd := &cobra.Command{
		Use:  "xdp_blacklist ",
		Long: "xdp_blacklist is a blacklist program implemented through ebpf",
	}
	rootCmd.AddCommand(NewAttachCommand())
	rootCmd.AddCommand(NewAddCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println("execute err", err)
	}
}
