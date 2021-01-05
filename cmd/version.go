package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "版本信息",
	Long:  "当前版本信息",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("当前版本 V1.0.0")
	},
}
