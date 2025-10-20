package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cli-service",
	Short: "서비스 관리 CLI 도구",
	Long: `서비스 관리 CLI 도구.

Example:
  전체 서비스 상태 확인 : cli-service state
  특정 서비스 상태 확인 : cli-service state --nm="auth-service"
  서비스 리셋          : cli-service reset --nm="auth-service"
  서비스 디버깅 설정    : cli-service debug --nm="auth-service" --level=1`,
	Run: rootCommand,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
}

func rootCommand(cmd *cobra.Command, args []string) {
	_ = cmd.Help()
}
