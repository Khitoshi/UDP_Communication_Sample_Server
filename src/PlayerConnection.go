package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
)

// UDP通信で受け取るデータ
type recvData struct {
	Name string `json:"name"`
	//RoomNumber int    `json:"room_number"`
	Message string `json:"message"`
}

// ユーザーのパラメータ
type UserParam struct {
	//位置データ
	//向いている方向
	Message string
}

// プレイヤーと通信処理する構造体
type PlayerConnection struct {
	conn *net.UDPConn
}

// PlayerConnection コンストラクタ
func NewPlayerConnection() (*PlayerConnection, error) {
	//サーバーアドレスの生成
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP(SERVER_IP),
		Port: SERVER_PORT,
	}

	// ソケットの生成
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		//fmt.Println("ソケットの生成に失敗しました:", err)
		//log.Fatal("ソケットの生成に失敗しました:", err)
		return nil, err
	}

	//生成したソケットを含めて返す
	return &PlayerConnection{conn: conn}, err
}

// PlayerConnection 更新処理
func (pc *PlayerConnection) UpdatePlayerConnection() {
	fmt.Println("サーバーが起動しました。")

	//ユーザーのパラメーター
	userParam := map[string][]UserParam{}
	for {
		//受信
		buffer := make([]byte, BUFFER_SIZE)
		recvbyte, clientAddr, err := pc.conn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("データの受信中にエラーが発生しました: ", err)
			continue
		}

		//受信したデータを構造体形式に変換する
		var recvdata recvData
		err = json.Unmarshal(buffer[:recvbyte], &recvdata)
		if err != nil {
			log.Println("jsonデータの変換に失敗: ", err)
			continue
		}

		//受信したデータを登録
		userParam[recvdata.Name] = append(userParam[recvdata.Name], UserParam{Message: recvdata.Message})

		//output
		fmt.Println(userParam)

		//TODO:デバック時のみ有効にする
		//終了宣言のチェック
		if recvdata.Message == "exit" {
			break
		}

		//送信
		_, err = pc.conn.WriteToUDP(buffer, clientAddr)
		if err != nil {
			log.Println("read error: ", err)
			continue
		}
	}
}
