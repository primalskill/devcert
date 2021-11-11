package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "devcert [space separated domain names]",
		Short: "short descr",
		Long:  `long description`,
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			err := devcertExec(args)
			if err != nil {
				return err
			}
			return nil
		},
	}

	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

}
