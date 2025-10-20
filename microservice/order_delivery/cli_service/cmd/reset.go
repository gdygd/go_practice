package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var resetArg ResetCommand

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "서비스 리셋",
	Long: `reset command.

Example:  
  cli-service reset --nm="auth-service"
  `,
	Run: func(cmd *cobra.Command, args []string) {
		if resetArg.serviceName == "" {
			fmt.Println("리셋할 서비스명을 --nm 옵션으로 지정해야 합니다.")
			return
		}
		fmt.Printf("%s 리셋\n", resetArg.serviceName)
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
	resetCmd.Flags().StringVarP(&resetArg.serviceName, "nm", "n", "", "서비스 이름 (예: auth-service)")
}
