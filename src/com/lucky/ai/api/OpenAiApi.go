package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var client = &http.Client{}

type request struct {
	Model       string    `json:"model"`
	Message     []message `json:"messages"`
	MaxTokens   int       `json:"-"`
	Temperature float32   `json:"temperature"`
	TopP        float32   `json:"top_p"`
	N           int       `json:"n"`
	Stream      bool      `json:"stream"`
	Stop        string    `json:"stop"`
}

type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type response struct {
	Model    string
	messages []string
	Choices  []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

// RequestOpenAiChat 请求OpenAi 聊天模型
// https://platform.openai.com/docs/api-reference/introduction 接口文档
func RequestOpenAiChat(msg string) []byte {

	url := "https://api.openai.com/v1/chat/completions"
	messages := make([]message, 1)
	messages[0] = message{
		Role:    "user",
		Content: msg,
	}
	reqData := request{
		Model:       "gpt-3.5-turbo",
		Message:     messages,
		Temperature: 0.5,
		N:           1,
		Stream:      true,
		Stop:        "\n",
	}

	reqBody, _ := json.Marshal(reqData)

	resp := httpRequest(http.MethodPost, url, bytes.NewReader(reqBody))

	defer func() {
		if nil != resp {
			_ = resp.Body.Close()
		}
	}()

	reader := bufio.NewReader(resp.Body)

	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	for {
		var delta struct {
			Role string `json:"role"`
		}
		err := decoder.Decode(&delta)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error:", err)
			continue
		}
		switch delta.Role {
		case "batch":
			break
		case "additions":
			fmt.Printf("Completion: %s\n", decoder)
		case "selection":
			break
		default:
			fmt.Println("Unknown delta role:", delta.Role)
			break
		}
	}

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("成功,返回结果：", string(body))

	return body
}

func RequestCompletions(msg string) []byte {
	url := "https://api.openai.com/v1/completions"
	type request struct {
		Model  string  `json:"model"`
		Prompt string  `json:"prompt"`
		TopP   float32 `json:"top_p"`
	}
	reqData := request{Model: "text-davinci-003", Prompt: msg, TopP: 0.5}
	reqBody, _ := json.Marshal(reqData)

	resp := httpRequest(http.MethodPost, url, bytes.NewReader(reqBody))

	defer func() {
		if nil != resp {
			_ = resp.Body.Close()
		}
	}()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println("成功,返回结果：", string(body))

	return body
}

// httpRequest 发送http请求
func httpRequest(method, url string, body io.Reader) *http.Response {
	apiKey := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ik1UaEVOVUpHTkVNMVFURTRNMEZCTWpkQ05UZzVNRFUxUlRVd1FVSkRNRU13UmtGRVFrRXpSZyJ9.eyJodHRwczovL2FwaS5vcGVuYWkuY29tL3Byb2ZpbGUiOnsiZW1haWwiOiIweHRvb2xraXRAZ21haWwuY29tIiwiZW1haWxfdmVyaWZpZWQiOnRydWV9LCJodHRwczovL2FwaS5vcGVuYWkuY29tL2F1dGgiOnsidXNlcl9pZCI6InVzZXItR3kyREo4SlduaVNlV2xleWxzN0p3TUlkIn0sImlzcyI6Imh0dHBzOi8vYXV0aDAub3BlbmFpLmNvbS8iLCJzdWIiOiJhdXRoMHw2MzkwNDlkMTI1MDViMTBkOGQxNzBmN2IiLCJhdWQiOlsiaHR0cHM6Ly9hcGkub3BlbmFpLmNvbS92MSIsImh0dHBzOi8vb3BlbmFpLm9wZW5haS5hdXRoMGFwcC5jb20vdXNlcmluZm8iXSwiaWF0IjoxNjc5ODA1NTY3LCJleHAiOjE2ODEwMTUxNjcsImF6cCI6IlRkSkljYmUxNldvVEh0Tjk1bnl5d2g1RTR5T282SXRHIiwic2NvcGUiOiJvcGVuaWQgcHJvZmlsZSBlbWFpbCBtb2RlbC5yZWFkIG1vZGVsLnJlcXVlc3Qgb3JnYW5pemF0aW9uLnJlYWQgb2ZmbGluZV9hY2Nlc3MifQ.F_yhY9CGUirDPY-DrmUKyyRtjH1g5gnLkgW6gAT37uAJNEcyEOASwUFiF1jNGkprSlaaRnebv-ZwAMMpmxUbWB_imVViPb-lp7yTd2LLHFhuH19ZBR963PPZM7JC5GIDfBMTyI6W_k8YohzOoXhNHokYTtdTid2-tbk21f0W7RtUVKxTxrEKc0fcFp0ljDZBqG8ldizDQxo_7W_4_WalS5ijc_rcp56wg9oJRjyjm0_9iraeFy4u6KQnoeB7b9CBOalcofotM_7-jZ9G8sdmp6eMBEQFwmgVWVosWKJEX6-XrUtK_7ZaprNrraeuHE_PR_zIvI3qlyD39l73v259jw"
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		fmt.Println("无法连接:", err)
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if nil != err {
		fmt.Println("请求失败:", err)
	}
	return resp
}
