package line

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Bot type
type Bot struct {
	ChannelID     string
	ChannelSecret string
	LineMID       string
}

// ReceiveRequest function
func (bot *Bot) ReceiveRequest(req *http.Request) (*MessageResults, error) {
	defer req.Body.Close()
	// check signature
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	if !bot.checkSignature(req.Header.Get("X-Line-Channelsignature"), body) {
		return nil, errors.New("invalid signature")
	}

	results := &MessageResults{}
	if json.Unmarshal(body, results); err != nil {
		return nil, err
	}
	return results, nil
}

// SendText function
func (bot *Bot) SendText(to string, message string) ([]byte, error) {
	// create request
	payload, err := json.Marshal(&Message{
		To:        []string{to},
		ToChannel: ChannelForSendingMessage,
		EventType: EventTypeForSendingMessage,
		Content: MessageContent{
			ContentType: ContentTypeText,
			ToType:      ToTypeForUser,
			Text:        message,
		},
	})
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", EndpointEventsURL, bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json; charser=UTF-8")
	req.Header.Set("X-Line-ChannelID", bot.ChannelID)
	req.Header.Set("X-Line-ChannelSecret", bot.ChannelSecret)
	req.Header.Set("X-Line-Trusted-User-With-ACL", bot.LineMID)

	// send and read response
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (bot *Bot) checkSignature(signature string, body []byte) bool {
	decoded, _ := base64.StdEncoding.DecodeString(signature)
	hash := hmac.New(sha256.New, []byte(bot.ChannelSecret))
	hash.Write(body)
	return hmac.Equal(hash.Sum(nil), decoded)
}
