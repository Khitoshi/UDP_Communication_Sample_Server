package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
)

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
