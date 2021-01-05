package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show current version",
	Long:  "当前版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("当前版本 V2.0.0\n")
	},
}
