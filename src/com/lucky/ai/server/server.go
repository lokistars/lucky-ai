package server

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"io"
	"lucky-ai/src/com/lucky/ai/api"
	"net"
	"net/http"
)

func Server() {
	http.HandleFunc("/ai", openAiHandle)
	http.HandleFunc("/openai", webSocketHandShake)
	http.HandleFunc("/openai/modelList", gatModelHandle)
	http.HandleFunc("/openai/modelDetails", gatModelHandle)
	http.HandleFunc("/images", imageHandle)

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
	fmt.Print("Request socket URL:", r.RemoteAddr, ": ")
	for {
		_, data, err := conn.ReadMessage()

		if nil != err {
			fmt.Printf("异常关闭:%v", err)
			break
		}
		fmt.Printf("收到消息:%v", data)
		api.RequestCompletions(string(data))
	}
}

// gatModel
func gatModelHandle(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		http.NotFound(w, r)
		return
	}
	fmt.Print("Request model URL:", r.RemoteAddr, ": ")
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
				_, _ = w.Write([]byte("model is null"))
				return
			}
			bytes := api.RetrieveModel(value)
			_, _ = w.Write(bytes)
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
	fmt.Print("Request ai URL:", r.RemoteAddr, ": ")
	switch r.Method {
	case http.MethodGet:
		msg := r.FormValue("msg")
		if "" == msg {
			_, _ = w.Write([]byte("msg is null"))
			fmt.Println()
			return
		}
		api.RequestOpenAiChat(msg, w)
		break
	default:
		http.NotFound(w, r)
	}
}

func imageHandle(w http.ResponseWriter, r *http.Request) {
	if nil == w || nil == r {
		http.NotFound(w, r)
		return
	}
	fmt.Print("Request images URL:", r.RemoteAddr, ": ")
	switch r.Method {
	case http.MethodGet:

		value := r.FormValue("prompt")
		if "" == value {
			_, _ = w.Write([]byte("prompt is null"))
			return
		}
		w.Header().Set("Content-Type", "image/png")

		images := api.RequestOpenAiImages(value)

		// 读取图像文件
		file, err := http.Get(images)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Body.Close()

		// 将图像写入响应中
		_, err = io.Copy(w, file.Body)
		if err != nil {
			fmt.Println(err)
			return
		}
		break
	default:
		http.NotFound(w, r)
	}
}
