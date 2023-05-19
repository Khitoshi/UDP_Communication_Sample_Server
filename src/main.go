package main

import (
	"fmt"
	"log"
)

func main() {
	//プレイヤーとの通信する構造体生成
	pc, err := NewPlayerConnection()
	if err != nil {
		log.Fatal("PlayerConnectionの生成に失敗: ", err)
	}

	endchan := make(chan bool, 1)
	//更新処理
	go func() {
		pc.UpdatePlayerConnection()
		endchan <- true
	}()

	//待機
	<-endchan

	pc.conn.Close()
	//defer conn.Close()
	fmt.Println("サーバーが終了しました。")
}
