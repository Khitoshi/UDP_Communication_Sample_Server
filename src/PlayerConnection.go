package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

/*
// UDP通信で受け取るデータ
type recvData struct {
	Name string `json:"name"`
	//RoomNumber int    `json:"room_number"`
	Message string `json:"message"`
}

// TODO:構造体の名前が変
// UDP通信で受け取るデータ
type sendData struct {
	data             recvData
	playerConnection PlayerConnection
	addr             *net.UDPAddr
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
// pings受信用
// pong送信用
func (pc *PlayerConnection) RecvPlayerConnection(userdata chan<- sendData) {
	//ユーザーのパラメーター
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

		//チャンネルに情報を渡す
		userdata <- sendData{
			playerConnection: PlayerConnection{conn: pc.conn},
			data:             recvdata,
			addr:             clientAddr,
		}
	}
}

func (pc *PlayerConnection) SendPlayerConnection(senddata sendData) {
	//ユーザーのパラメーター
	for {

		// JSONエンコード
		jsonData, err := json.Marshal(senddata.data)
		if err != nil {
			log.Println("json encode error: ", err)
			continue
		}

		//送信
		_, err = pc.conn.WriteToUDP(jsonData, senddata.addr)
		if err != nil {
			log.Println("read error: ", err)
			continue
		}
	}
}
*/

// プレイヤーの情報
type Player struct {
	// A unique identifier for the player.
	ID      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
	// Additional fields to hold player state can be added here.
	// The specific fields will depend on the requirements of your game.
}

// Server represents a game server.
type Server struct {
	conn *net.UDPConn
	//gorutine管理用map
	players map[string]*Player
	// プレイヤーごとのチャネルを管理するマップ
	playerChans map[string]chan *Player
	lock        sync.Mutex
}

// serverの初期処理をしたインスタンス作成
func NewServer() (*Server, error) {
	//サーバーアドレスの生成
	serverAddr := &net.UDPAddr{
		IP:   net.ParseIP(SERVER_IP),
		Port: SERVER_PORT,
	}

	//ソケット作成
	conn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		return nil, err
	}

	return &Server{
		conn:        conn,
		players:     make(map[string]*Player),
		playerChans: make(map[string]chan *Player),
	}, nil

}

// serverメインloop
func (s *Server) Listen() error {
	defer s.conn.Close()

	buf := make([]byte, 1024)
	endchan := make(chan error)
	fmt.Println("server起動")
	//受信
	go func() {
		for {
			n, addr, err := s.conn.ReadFromUDP(buf)
			if err != nil {
				//return err
				endchan <- err
			}

			// Create a new goroutine for each message received.
			go s.handleMessage(addr, buf[:n])
		}
	}()

	return <-endchan
}

// プレイヤーの情報をプレイヤー情報処理用gorutineに渡す
func (s *Server) handleMessage(addr *net.UDPAddr, message []byte) {
	// Here you would parse the message to extract the player ID and the
	// action the player is performing. The specifics of this will depend
	// on how you have chosen to structure your messages.

	// Example:
	// playerID, action := parseMessage(message)

	// Lock the server to safely read from the players map.

	var p Player
	err := json.Unmarshal(message, &p)
	if err != nil {
		//log.Println("jsonデータの変換に失敗: ", err)
		log.Fatal("jsonデータの変換に失敗: ", err)
	}
	fmt.Println("data:", p)
	s.lock.Lock()
	//player, ok := s.players[p.ID]
	_, ok := s.players[p.ID]
	fmt.Println("id: ", p.ID)
	if !ok {
		// If the player does not yet exist, create a new Player object and
		// a new goroutine to handle the player's actions.
		//player = p
		newPlayer := &p
		s.players[p.ID] = newPlayer
		s.playerChans[p.ID] = make(chan *Player, 1)
		s.playerChans[p.ID] <- newPlayer
		fmt.Println("new gorutine: ", newPlayer)
		go s.handlePlayer(&p)
	} else {
		fmt.Println("inc chan: ", p)
		s.playerChans[p.ID] <- &p
	}
	s.lock.Unlock()
	fmt.Printf("\n\n")
	// Send the action to the player's goroutine.
	// player.actions <- action
}

// プレイヤーの数だけ存在するgorutine
// プレイヤーの情報の処理
func (s *Server) handlePlayer(player *Player) {
	// Loop, reading actions from the player's channel and processing them.
	// The specifics of this will depend on your game's requirements.

	// Example:
	// for action := range player.actions {
	//     // Process the action.
	// }

	for {
		action := <-s.playerChans[player.ID]

		action.Message = "server" + action.Message

		s.players[player.ID] = action
	}
}
