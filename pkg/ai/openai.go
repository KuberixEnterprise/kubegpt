package ai

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
)

func GetAnswer(content string, token string) string {

	client := openai.NewClient(token)
	response, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4TurboPreview,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: "당신의 역할은 Kubernetes 상의 오류를 해결해주는 것입니다. " +
						"Kubernetes 리소스 yaml 를 바탕으로 어디에서 오류가 났는지, 해결방안은 무엇인지 제시해야 합니다." +
						"Kubernetes 리소스 yaml 로는 정보가 부족해도 해당 오류가 생길 수 있는 경우 중 제일 연관성이 높은 순으로 3개 나열하고 그에 대한 해결방안을 제시해주십시오." +
						"명확하고 정확한 설명을 강조하고 부정확하거나 오해의 소지가 있는 정보를 제공하지 않도록 하십시오." +
						"답변은 제일 연관성이 높은 순서대로 어떤 리소스를 어떻게 고쳐야하는지 자세하게 답해야 합니다.",
				},
				// {
				// 	Role: openai.ChatMessageRoleSystem,
				// 	Content: "Your role is to solve errors in Kubernetes. Based on Kubernetes resource yaml, you need to present where the error occurred and what the solution is." +
				// 		"With Kubernetes resource yaml, please list all the numbers of errors that may occur even if there is insufficient information and provide solutions for them." +
				// 		"Highlight clear and accurate descriptions and avoid providing inaccurate or misleading information." +
				// 		"Answers should be answered in detail about which resources should be fixed in the most relevant order.",
				// },
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

	// fmt.Println(response.Choices[0].Message.Content)
	return response.Choices[0].Message.Content

}
