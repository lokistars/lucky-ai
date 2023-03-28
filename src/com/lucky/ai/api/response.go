package api

type responseImage struct {
	Created int `json:"created"`
	Data    []struct {
		B64Json string `json:"b64_json"`
		Url     string `json:"url"`
	} `json:"data"`
}

type responseChat struct {
	Id      string `json:"id"`
	Model   string `json:"model"`
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
		Index        string `json:"index"`
	} `json:"choices"`
}
