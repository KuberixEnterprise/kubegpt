package ai

import (
	"context"
	"strings"

	b64 "encoding/base64"

	"github.com/kuberixenterprise/kubegpt/api/v1alpha1"
	openai "github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
)

func GetAnswer(content string, kubegpt v1alpha1.KubegptSpec) string {
	token64, _ := b64.StdEncoding.DecodeString(kubegpt.AI.Secret.Key)
	// 개행문자 제거
	token := strings.Replace(string(token64), "\n", "", -1)
	model := kubegpt.AI.Model
	promptTmpl := PromptMap[kubegpt.AI.Language]

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
					Content: content,
				},
			},
		},
	)

	if err != nil {
		log.Errorf("ChatCompletion error: %v\n", err)
		return err.Error()
	}

	return response.Choices[0].Message.Content

}
