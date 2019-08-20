package main

import (
	"fmt"
	"github.com/elek/klepif/pkg/run"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "klepif",
		Short: "Application to process new changes from github PRs",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check changes since the last run and execute actions.",
		Run: func(cmd *cobra.Command, args []string) {
			err := run.Check()
			if err != nil {
				panic(err)
			}
		},
	})

	runCmd := &cobra.Command{
		Use:   "run",
		Args:  cobra.ExactArgs(1),
		Short: "Execute actions for the last events of a specific PR.",
		Run: func(cmd *cobra.Command, args []string) {
			prNum, err := strconv.Atoi(args[0])
			if err != nil {
				panic(err)
			}
			err = run.Run(prNum)
			if err != nil {
				panic(err)
			}
		},
	}
	
	rootCmd.AddCommand(runCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
