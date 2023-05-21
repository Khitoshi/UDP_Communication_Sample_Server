package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	services "UDPserver/src/protocode"

	"google.golang.org/grpc"
)

type myServer struct {
	services.UnimplementedGreetingServiceServer
}

func main() {

	/*
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
	*/

	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// gRPCサーバーを作成
	s := grpc.NewServer()

	//gRPCサーバーにGreetingServiceを登録
	hellopb.RegisterGreetingServiceServer(s)

	// 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()

}
