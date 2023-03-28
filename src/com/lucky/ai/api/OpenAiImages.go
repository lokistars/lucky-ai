package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// RequestOpenAiImages 请求OpenAi 图片模型
func RequestOpenAiImages(prompt string) string {
	url := "https://api.openai.com/v1/images/generations"
	fmt.Println("prompt:", prompt)
	reqData := requestImage{
		Prompt:         prompt,
		Size:           "1024x1024",
		N:              1,
		ResponseFormat: "url",
	}
	reqBody, _ := json.Marshal(reqData)
	resp := httpRequest(http.MethodPost, url, bytes.NewReader(reqBody))

	defer func() {
		if nil != resp {
			_ = resp.Body.Close()
		}
	}()
	decoder := json.NewDecoder(resp.Body)

	resBody := &responseImage{}

	_ = decoder.Decode(resBody)
	return resBody.Data[0].Url
}

// RequestOpenAiImageEdit 请求OpenAi 图片编辑模型
func RequestOpenAiImageEdit(prompt string) string {
	url := "https://api.openai.com/v1/images/edits"
	reqData := ImageEdit{
		Image:          "",
		Mask:           "",
		Prompt:         prompt,
		Size:           "1024x1024",
		N:              1,
		ResponseFormat: "url",
	}
	reqBody, _ := json.Marshal(reqData)
	resp := httpRequest(http.MethodPost, url, bytes.NewReader(reqBody))
	defer func() {
		if nil != resp {
			_ = resp.Body.Close()
		}
	}()
	return ""
}
