package main

import (
	"context"
	"fmt"
	"log"
	"os"

	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

func main() {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY가 설정되어 있지 않습니다.")
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
	)

	resp, err := client.Chat.Completions.New(
		context.Background(),
		openai.ChatCompletionNewParams{
			Model: openai.ChatModelGPT4o, // 상수로 모델 지정
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.UserMessage("Go 언어에서 goroutine이 무엇인가요?"),
			},
		},
	)
	if err != nil {
		log.Fatalf("ChatCompletion 오류: %v", err)
	}

	fmt.Println(resp.Choices[0].Message.Content)
}
