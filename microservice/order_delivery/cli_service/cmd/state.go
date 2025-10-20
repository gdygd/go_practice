package cmd

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gdygd/goglib"
	"github.com/spf13/cobra"
)

var stateArg StateCommand

var stateCmd = &cobra.Command{
	Use:   "state",
	Short: "서비스 상태 조회",
	Long: `state command.

Example:  
  cli-service state (전체 서비스 확인)
  cli-service state --nm="auth-service (특정 서비스 확인)"
  `,
	Run: serviceState,
}

func init() {
	rootCmd.AddCommand(stateCmd)
	stateCmd.Flags().StringVarP(&stateArg.serviceName, "nm", "n", "", "서비스 이름 (예: auth-service)")
}

func serviceState(cmd *cobra.Command, args []string) {
	if stateArg.serviceName == "" {
		fmt.Println("전체 서비스 상태 조회")
		checkService()
	} else {
		fmt.Printf("%s 상태 조회\n", stateArg.serviceName)
	}
}

func checkSercviceState(url string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// HttpRequest안에서 timeout처리가 됨.
	statuscode, _, err := goglib.HttpRequest(ctx, http.Header{}, nil, "GET", url+"/heartbeat")
	if err != nil || statuscode != http.StatusOK {
		// fmt.Printf("Check Heartbeat fail.. [%s]", url)
		return false
	}

	return true
}

type ServiceState struct {
	Name  string
	State bool
}

func checkService() {
	var wg sync.WaitGroup
	var servState map[string]bool = make(map[string]bool)

	for i, svc := range services {

		wg.Add(1)
		go func(idx int, url string) {
			defer wg.Done()
			rst := checkSercviceState(svc.Url)
			servState[svc.Name] = rst
		}(i, svc.Url)
	}

	wg.Wait()

	printLine1("=", 50)
	fmt.Printf("%-20s | %-5s \n", "SERVICE", "STATE")
	printLine1("-", 50)

	for nm, stt := range servState {
		fmt.Printf("%-20s | %-5v \n", nm, stt)
	}

	printLine1("=", 50)
}
