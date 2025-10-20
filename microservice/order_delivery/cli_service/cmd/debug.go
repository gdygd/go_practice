package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	debugName  string
	debugLevel int
)
var debugArg DebugCommand

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "디버깅 설정",
	Long: `debug command.

Example:  
  cli-service debug --nm="auth-service" --level=1
  `,

	Run: func(cmd *cobra.Command, args []string) {
		if debugArg.serviceName == "" {
			fmt.Println("디버깅할 서비스명을 --nm or -n 옵션으로 지정해야 합니다.")
			return
		}
		fmt.Printf("%s 디버깅 레벨 설정: %d\n", debugArg.serviceName, debugArg.level)
	},
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.Flags().StringVarP(&debugArg.serviceName, "nm", "n", "", "서비스 이름 (예: auth-service)")
	debugCmd.Flags().IntVarP(&debugArg.level, "level", "l", 0, "디버깅 레벨 (예: 1)")
}
