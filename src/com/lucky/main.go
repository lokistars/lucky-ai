package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	// 通过代理请求 curl https://www.google.com -x http://127.0.0.1:7890
	//server.Server()
	str := "data: {\"id\":\"chatcmpl-6yOiFpOOqOOsfkj2285j6p1CHZr9Z\",\"object\":\"chat.completion.chunk\",\"created\":1679852695,\"model\":\"gpt-3.5-turbo-0301\",\"choices\":[{\"delta\":{\"content\":\"我\"},\"index\":0,\"finish_reason\":null}]}"
	//fmt.Println(str)
	bytes := []byte(str)
	for {
		var responseJSON map[string]interface{}
		err := json.Unmarshal(bytes, &responseJSON)
		if err != nil {
			break
		}
		choices := responseJSON["choices"].([]interface{})
		if len(choices) > 0 {
			text := choices[0].(map[string]interface{})["delta"].(map[string]interface{})["content"].(string)
			fmt.Print(text)
		}
	}
}
