package pkg

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Notifier interface
type Notifier interface {
	Notify(p Product, i ReleaseInfo)
}

func NewNilNotifier() *NilNotifier {
	return &NilNotifier{}
}

// NilNotifier implementation of Notifier that does nothing
type NilNotifier struct{}

// Notify is a Nil operation
func (n NilNotifier) Notify(_ Product, _ ReleaseInfo) {}

type SlackNotifier struct {
	webhookURL string
}

// NewSlackNotifier creates a new instance of SlackNotifier with the specified Slack webhook URL.
// This webhook URL is used for sending notifications to Slack.
//
// Parameters:
//
//	url (string): The Slack webhook URL where notifications will be sent.
//
// Returns:
//
//	*SlackNotifier: A pointer to the newly created instance of SlackNotifier.
//
// Example:
//
//	notifier := pkg.NewSlackNotifier("https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX")
//	This creates a new SlackNotifier that will send notifications to the specified webhook URL.
func NewSlackNotifier(url string) *SlackNotifier {
	return &SlackNotifier{webhookURL: url}
}

// Notify sends a notification message to Slack using the SlackNotifier's webhook URL.
// The message includes information about the product and its release.
//
// Parameters:
//
//	p (Product): The product information, including its name and version.
//	i (ReleaseInfo): The release information, including the release date and end of life date.
//
// This method constructs a message from the provided product and release information
// and sends this message as a Slack notification. It does not return any value or error.
// Errors during message construction or sending are logged.
//
// Example usage:
//
//	product := pkg.Product{Name: "MyProduct", Version: "1.2.3"}
//	releaseInfo := pkg.ReleaseInfo{ReleaseDate: "2024-01-01", EndOfLifeDate: "2025-01-01"}
//	notifier.Notify(product, releaseInfo)
//	This sends a notification to Slack with the specified product and release information.
func (s *SlackNotifier) Notify(p Product, i ReleaseInfo) {
	message, err := GenerateSlackMessage(p, i)
	if err != nil {
		return
	}

	req, err := http.NewRequest("POST", s.webhookURL, bytes.NewBuffer(message))
	if err != nil {
		log.Printf("error creating request: %v", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error sending notification to Slack: %v", err)
		return
	}
	defer resp.Body.Close()
}

// SlackMessageBlock represents a single block in the Slack message.
type SlackMessageBlock struct {
	Type   string       `json:"type"`
	Text   *TextObject  `json:"text,omitempty"`
	Fields []TextObject `json:"fields,omitempty"`
}

// TextObject represents text elements in Slack message blocks.
type TextObject struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Emoji *bool  `json:"emoji,omitempty"`
}

// SlackMessage represents the entire Slack message in Block Kit format.
type SlackMessage struct {
	Blocks []SlackMessageBlock `json:"blocks"`
}

// GenerateSlackMessage generates a structured Slack message for product lifecycle alerts.
func GenerateSlackMessage(product Product, release ReleaseInfo) ([]byte, error) {
	message := SlackMessage{
		Blocks: []SlackMessageBlock{
			{
				Type: "section",
				Fields: []TextObject{
					{Type: "mrkdwn", Text: fmt.Sprintf("Product Name:\n*%s*", product.Name)},
					{Type: "mrkdwn", Text: fmt.Sprintf("Version:\n*%s*", product.Version)},
					{Type: "mrkdwn", Text: fmt.Sprintf("Label:\n*%s*", product.Label)},
					{Type: "mrkdwn", Text: fmt.Sprintf("Release Date:\n*%s*", release.ReleaseDate)},
					{Type: "mrkdwn", Text: fmt.Sprintf("End of Life Date:\n*%s*", release.EndOfLifeDate)},
				},
			},
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}

	return jsonMessage, nil
}
