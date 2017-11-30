package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	cobra.OnInitialize()
}

var rootCmd = &cobra.Command{
	Use:   "garupalife",
	Short: "motto garupa life command line tool",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
