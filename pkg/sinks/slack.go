package sinks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/kuberixenterprise/kubegpt/pkg/cache"
	"net/http"

	"github.com/kuberixenterprise/kubegpt/api/v1alpha1"
	log "github.com/sirupsen/logrus"
)

type SlackSink struct {
	Endpoint string
	Client   *Client
}

type SlackMessage struct {
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
}

type TextObject struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Field struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Attachment struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Color string `json:"color"`
	Title string `json:"title"`
}

func formatEvent(event v1alpha1.Event) string {
	// Format the event into a string.
	// Modify this according to how you want to display each event.
	return fmt.Sprintf("- Event: %s\n- Count: %v\n- Reason: %s\n- Message: %s", event.Type, event.Count, event.Reason, event.Message)
}
func (s *SlackSink) Configure(config *v1alpha1.Kubegpt) {
	s.Endpoint = config.Spec.Sink.Endpoint
	s.Client = NewClient()
}

func buildSlackMessage(result v1alpha1.ResultSpec, k8sgptCR string) SlackMessage {
	var detailsText string
	for _, event := range result.Event {
		// Use the formatEvent function to get the string representation of each event
		detailsText += formatEvent(event) + "\n" // Add a newline after each event
	}

	labelsText := fmt.Sprintf("%v", result.Labels)
	imagesText := fmt.Sprintf("%v", result.Images)

	return SlackMessage{
		Text: fmt.Sprintf(">*[%s] Error event occurred in the %s : %s \n -Label: %s \n -Image: %s*", k8sgptCR, result.Kind, result.Name, labelsText, imagesText),
		Attachments: []Attachment{
			{
				Type:  "mrkdwn",
				Text:  detailsText,
				Color: "danger",
				Title: "Detailed Report",
			},
		},
	}
}

func StringSlackMessage(text string, result v1alpha1.ResultSpec) SlackMessage {
	var title string
	for _, event := range result.Event {
		title += fmt.Sprintf(">*Event: %s / Count: %v / Reason: %s  \nMessage: %s*", event.Type, event.Count, event.Reason, event.Message)
	}
	return SlackMessage{
		Text: "[Error Message]" + title,
		Attachments: []Attachment{
			{
				Type:  "mrkdwn",
				Color: "good",
				Title: "ChatGPT Answer : ",
				Text:  text,
			},
		},
	}
}

func (c *Client) SendHTTPRequest(method, url string, body []byte) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		log.WithError(err).WithField("component", "HTTPClient").Error("Failed to create HTTP request")
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.hclient.Do(req)
	if err != nil {
		log.WithError(err).WithField("component", "HTTPClient").Error("Failed to send HTTP request")
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.WithField("status", resp.Status).WithField("component", "HTTPClient").Error("HTTP request returned non-success status")
		return resp, fmt.Errorf("HTTP request returned non-success status: %s", resp.Status)
	}

	return resp, nil
}

func (s *SlackSink) Emit(results v1alpha1.ResultSpec, kubegpt v1alpha1.KubegptSpec) (string, error) {
	message := buildSlackMessage(results, "Kubegpt")
	jsonData, err := json.Marshal(message)
	SlackClient(s, jsonData, results.Name)

	if err != nil {
		log.WithError(err).WithField("component", "SlackSink").Error("Failed to marshal message")
		return "", err
	}
	gptMsg := message.Text + message.Attachments[0].Text
	log.Info(message.Attachments[0].Text)
	return gptMsg, nil
}

func RebuildSlackMessage(key string, cachedData cache.CacheItem) SlackMessage {
	return SlackMessage{
		Attachments: []Attachment{
			{
				Type:  "mrkdwn",
				Color: "danger",
				Title: "Unresolved errors",
				Text:  fmt.Sprintf("*%s*\n*Error Message:* \n%s", key, cachedData.Message),
			},
			{
				Type:  "mrkdwn",
				Color: "good",
				Title: "Previous answer",
				Text:  fmt.Sprintf("*%s*", cachedData.Answer),
			},
		},
	}
}

func (s *SlackSink) ReEmit(key string, cachedData cache.CacheItem) error {
	message := RebuildSlackMessage(key, cachedData)
	jsonData, err := json.Marshal(message)
	if err != nil {
		log.WithError(err).WithField("component", "SlackSink").Error("Failed to marshal message")
		return err
	}

	SlackClient(s, jsonData, key) // Send the message to Slack

	log.Printf("Successfully sent report to Slack for %v", key)
	return nil
}

func SlackClient(s *SlackSink, sendData []byte, sendName string) error {
	req, err := http.NewRequest(http.MethodPost, s.Endpoint, bytes.NewBuffer(sendData))
	if err != nil {
		log.WithError(err).WithField("component", "SlackSink").Error("Failed to create HTTP request")
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.Client.hclient.Do(req)
	if err != nil {
		log.WithError(err).WithField("component", "SlackSink").Error("Error sending HTTP request")
		return err
	}

	if resp.StatusCode != http.StatusOK {
		log.WithField("status", resp.Status).Error("Failed to send report")
		return fmt.Errorf("failed to send report: %s", resp.Status)
	}
	log.Printf("Successfully sent report to Slack for %s", sendName)
	return nil
}
