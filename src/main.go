package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	services "github.com/Khitoshi/UDP_Communication_Sample_Server/src/pb"

	"google.golang.org/grpc"
)

type myServer struct {
	services.UnimplementedLoginServiceServer
}

// 自作サービス構造体のコンストラクタを定義
func NewMyServer() *myServer {
	return &myServer{}
}

// Unary RPCがレスポンスを返すところ
// ログイン処理を実行
// TODO:インターセプター処理の追加
func (m *myServer) Login(ctx context.Context, req *services.LoginRequest) (*services.LoginResponse, error) {
	/*
		//ctxからメタデータを取得
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			log.Println(md)
		}

		//ヘッダーを作成
		headerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "header"})
		if err := grpc.SetHeader(ctx, headerMD); err != nil {
			return nil, err
		}

		//トレイラーを作成
		trailerMD := metadata.New(map[string]string{"type": "unary", "from": "server", "in": "trailer"})
		if err := grpc.SetTrailer(ctx, trailerMD); err != nil {
			return nil, err
		}
	*/
	// HelloResponse型を1つreturnする
	// (Unaryなので、レスポンスを一つ返せば終わり)
	return &services.LoginResponse{
		Id:   fmt.Sprintf("0: %s", req.GetName()),
		Name: req.GetName(),
	}, nil
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
	services.RegisterLoginServiceServer(s, NewMyServer())

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
