package line

// MessageResults type
type MessageResults struct {
	Result []Message `json:"result"`
}

// Message type
type Message struct {
	ID          string         `json:"id,omitempty"`
	From        string         `json:"from,omitempty"`
	FromChannel int64          `json:"fromChannel,omitempty"`
	To          []string       `json:"to"`
	ToChannel   int64          `json:"toChannel"`
	EventType   string         `json:"eventType"`
	Content     MessageContent `json:"content"`
}

// MessageContent type
type MessageContent struct {
	ID              string            `json:"id,omitempty"`
	ContentType     int8              `json:"contentType"`
	From            string            `json:"from,omitempty"`
	CreatedTime     int64             `json:"createdTime,omitempty"`
	To              []string          `json:"to,omitempty"`
	ToType          int8              `json:"toType"`
	ContentMetadata map[string]string `json:"contentMetadata,omitempty"`
	Text            string            `json:"text"`
	Location        map[string]string `json:"location,omitempty"`
}
