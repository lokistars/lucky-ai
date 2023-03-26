package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"lucky-ai/src/com/lucky/ai/api"
	"net"
	"net/http"
)

func Server() {
	http.HandleFunc("/ai", openAiHandle)
	http.HandleFunc("/openai", webSocketHandShake)
	http.HandleFunc("/openai/modelList", gatModelHandle)
	http.HandleFunc("/openai/modelDetails", gatModelHandle)

	var str string
	addrs, _ := net.InterfaceAddrs()
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				str = ipnet.IP.String()
				break
			}
		}
	}
	fmt.Println("启动HTTP服务成功：地址：http://" + str + ":8083")

	err := http.ListenAndServe(":8083", nil)

	if nil != err {
		fmt.Println("启动HTTP服务失败:", err)
	}
}

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func webSocketHandShake(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		http.NotFound(w, r)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)

	if nil != err {
		fmt.Printf("WebSocket upgrade error, %v+", err)
		http.NotFound(w, r)
		return
	}

	defer func() {
		_ = conn.Close()
	}()

	for {
		_, data, err := conn.ReadMessage()

		if nil != err {
			fmt.Printf("异常关闭:%v", err)
			break
		}
		fmt.Printf("收到消息:%v", data)
	}
}

// gatModel
func gatModelHandle(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/openai/modelList" {
			model := api.GetOpenAiModel()

			marshal, _ := json.Marshal(model)
			_, err := w.Write(marshal)
			if err != nil {
				fmt.Printf("异常关闭:%v", err)
			}
		} else if r.URL.Path == "/openai/modelDetails" {
			value := r.FormValue("model")
			if "" == value {
				w.Write([]byte("model is null"))
				return
			}
			bytes := api.RetrieveModel(value)
			w.Write(bytes)
		}
		break
	default:
		http.NotFound(w, r)
	}
}

func openAiHandle(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		msg := r.FormValue("msg")
		if "" == msg {
			w.Write([]byte("msg is null"))
			return
		}
		bytes := api.RequestOpenAiChat(msg)
		w.Write(bytes)
		break
	default:
		http.NotFound(w, r)
	}
}
