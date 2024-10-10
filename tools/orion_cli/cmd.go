package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type cliLogger struct{}

func (cliLogger) Log(msg string) {
	fmt.Println("[orion-cli] " + msg)
}

var cl = cliLogger{}

var cmd = &cobra.Command{
	Use: "",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("* orion framework generator, please use '--help' option for more information.")
	},
	Example: "orion-cli --help",
}

var cmdGenerator = &cobra.Command{
	Use:     "new",
	Example: "orion-cli new folder_name",
	Run: func(cmd *cobra.Command, args []string) {
		handleCmd(func() error {
			gen := &OrionGenerator{}
			return gen.Check(args).Excute()
		})
	},
}

func handleCmd(c func() error) {
	// spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	// spin.Start()
	// defer spin.Stop()

	cl.Log("Preparing ...")
	if err := c(); err != nil {
		cl.Log("Excute error:" + err.Error())
		return
	}

	cl.Log("Excute succ!")
}

func init() {
	cmd.AddCommand(cmdGenerator)
}
