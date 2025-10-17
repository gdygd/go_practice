package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "service",
	Short: "MSA 서비스 관련 기능",
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}
