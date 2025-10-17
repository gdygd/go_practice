package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "msa-admin",
	Short: "MSA 서비스 관리용 CLI",
	Long:  `각 MSA 서비스의 상태조회, 배포, 토큰발급 등을 관리합니다.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MSA CLI - Use subcommands like 'service', 'auth', 'saga'")
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
