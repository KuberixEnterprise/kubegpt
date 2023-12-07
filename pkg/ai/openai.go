package ai

import (
	"context"
	"fmt"
	"strings"

	b64 "encoding/base64"

	"github.com/kuberixenterprise/kubegpt/api/v1alpha1"
	openai "github.com/sashabaranov/go-openai"
)

func GetAnswer(content string, kubegpt v1alpha1.KubegptSpec) string {
	token64, _ := b64.StdEncoding.DecodeString(kubegpt.AI.Secret.Key)
	// 개행문자 제거
	token := strings.Replace(string(token64), "\n", "", -1)
	model := kubegpt.AI.Model
	promptTmpl := PromptMap["default"]

	client := openai.NewClient(token)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: promptTmpl,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "문제해결해줘 " + content,
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return err.Error()
	}

	// tmp.Println(response.Choices[0].Message.Content)
	return response.Choices[0].Message.Content

}
