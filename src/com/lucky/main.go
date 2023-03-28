package main

import "lucky-ai/src/com/lucky/ai/server"

func main() {
	// 通过代理请求 curl https://www.google.com -x http://127.0.0.1:7890
	server.Server()
}
