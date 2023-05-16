package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

const (
	// サーバーのIPアドレス
	SERVER_IP = "localhost"
	// サーバーのポート番号
	SERVER_PORT = 1234
	// バッファサイズ
	BUFFER_SIZE = 1024
)

type Data struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func main() {
	// サーバーアドレスの生成
	/*
		serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", SERVER_PORT))
		if err != nil {
			fmt.Println("サーバーアドレスの生成に失敗しました:", err)
			os.Exit(1)
		}*/

	//サーバーアドレスの生成
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP(SERVER_IP),
		Port: SERVER_PORT,
	}

	// ソケットの生成
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("ソケットの生成に失敗しました:", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("サーバーが起動しました。")

	for {
		//受信
		buffer := make([]byte, BUFFER_SIZE)
		_, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("データの受信中にエラーが発生しました:", err)
			continue
		}

		//message := string(buffer)
		//fmt.Printf("クライアントからのメッセージ [%s]: %s\n", clientAddr.String(), message)

		//データが小さかったり大きかったりするとデータの最後にNULLなどが入りエラーが発生するため
		//jsonデータをトリミングする
		trimmedBuffer := bytes.Trim(buffer, "\x00")
		var d Data
		err = json.Unmarshal(trimmedBuffer, &d)
		if err != nil {
			fmt.Println(err)
			return
		}
		//output
		fmt.Println("Name:", d.Name)
		fmt.Println("Message:", d.Message)

		//送信
		_, err = conn.WriteToUDP(buffer, clientAddr)
		if err != nil {
			fmt.Println("read error: ", err)
		}
	}

	fmt.Println("サーバーが終了しました。")
}
