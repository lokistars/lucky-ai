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

// RequestOpenAiChat 请求OpenAi 聊天模型
// https://platform.openai.com/docs/api-reference/introduction 接口文档
func RequestOpenAiChat(msg string, w http.ResponseWriter) []byte {

	url := "https://api.openai.com/v1/chat/completions"

	fmt.Println("message:", msg)

	messages := make([]message, 1)
	messages[0] = message{
		Role:    "user",
		Content: msg,
	}
	reqData := requestChat{
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
	scanner := bufio.NewScanner(resp.Body)
	role := ""

	flusher, _ := w.(http.Flusher)
	w.Header().Set("Content-Type", "text/event-stream; charset=utf-8")
	w.Header().Set("Transfer-Encoding", "chunked") // 支持分块传输
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	for scanner.Scan() {
		line := scanner.Text()[6:]
		if line == "[DONE]" {
			fmt.Println()
			_, _ = w.Write([]byte("\n"))
			break
		}
		res := responseChat{}
		_ = json.Unmarshal([]byte(line), &res)
		scanner.Scan()
		choices := res.Choices[0]
		if choices.FinishReason == "stop" {
			role = ""
			continue
		}
		if role == "" {
			role = choices.Delta.Role
		}
		switch role {
		case "assistant":
			//_, _ = w.Write([]byte(choices.Delta.Content))
			_, _ = io.WriteString(w, choices.Delta.Content)
			flusher.Flush()

			// 实现流式响应
			//encoder := json.NewEncoder(w)
			//encoder.Encode(choices.Delta.Content)
			//w.(http.Flusher).Flush()
		case "additions":
		case "batch":
		case "selection":
		default:
			fmt.Println("Unknown delta role:", choices.Delta.Role)
			break
		}

	}

	return nil
}

// RequestCompletions 请求OpenAi 智能补全模型
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
