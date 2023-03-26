package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type model struct {
	Object string `json:"object"`
	Data   []struct {
		Id string `json:"id"`
	} `json:"data"`
	Ids []string `json:"ids"`
}

// GetOpenAiModel 获取模型列
func GetOpenAiModel() model {
	fmt.Println("开始获取模型列表")

	url := "https://api.openai.com/v1/models"

	resp := httpRequest(http.MethodGet, url, nil)

	defer func() {
		if nil != resp {
			_ = resp.Body.Close()
		}
	}()

	res, _ := io.ReadAll(resp.Body)

	data := model{}

	err := json.Unmarshal(res, &data)

	if nil != err {
		fmt.Println("解析模型列表失败:", err)
	}
	data.Ids = make([]string, len(data.Data))
	for i := range data.Data {
		data.Ids[i] = data.Data[i].Id
	}
	data.Data = nil
	return data
}

// RetrieveModel 检索模型
func RetrieveModel(model string) []byte {
	if model == "" {
		fmt.Println("模型不能为空")
		return nil
	}
	fmt.Println("开始检索模型：", model)

	url := "https://api.openai.com/v1/models/" + model
	resp := httpRequest(http.MethodGet, url, nil)

	defer func() {
		if nil != resp {
			_ = resp.Body.Close()
		}
	}()

	res, _ := io.ReadAll(resp.Body)
	return res
}
