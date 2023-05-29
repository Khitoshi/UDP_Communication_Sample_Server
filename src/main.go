package main

import (
	"fmt"
	"log"
)

func main() {

	server, err := NewServer()
	if err != nil {
		log.Fatal(err)
	}

	err = server.Listen()
	if err != nil {
		log.Fatal(err)
	}

	//server終了
	fmt.Println("server終了")

}
