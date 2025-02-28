package types

type AiMessage struct {
	ID          string                 `json:"id"`
	Object      string                 `json:"object"`
	Role        string                 `json:"role"`
	Content     []AiMessageContent     `json:"content"`
	Attachments []interface{}          `json:"attachments"`
	Metadata    map[string]interface{} `json:"metadata"`
}

type AiMessageRequest struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AiMessageContent struct {
	Type string               `json:"type"`
	Text AiMessageTextContent `json:"text"`
}

type AiMessageTextContent struct {
	Value       string        `json:"value"`
	Annotations []interface{} `json:"annotations"`
}
