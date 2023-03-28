package api

type requestChat struct {
	Model       string    `json:"model"`
	Message     []message `json:"messages"`
	MaxTokens   int       `json:"-"`
	Temperature float32   `json:"temperature"`
	TopP        float32   `json:"top_p"`
	N           int       `json:"n"`
	Stream      bool      `json:"stream"`
	Stop        string    `json:"-"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestImage struct {
	Prompt         string `json:"prompt"` // 描述信息
	N              int    `json:"n" default:"1"`
	Size           string `json:"size" default:"1024x1024"`
	ResponseFormat string `json:"response_format" default:"url"` // b64_json
}

type ImageEdit struct {
	Image          string `json:"image"` // 修改的头像
	Mask           string `json:"mask"`  // 附加头像,需要编辑的位置
	Prompt         string `json:"prompt"`
	Size           string `json:"size"`
	N              int    `json:"n"`
	ResponseFormat string `json:"response_format"` // 响应格式 url 或 b64_json

}
