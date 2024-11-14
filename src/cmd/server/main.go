package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "os"
    "os/signal"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"
    hellopb "mygrpc/pkg/grpc"
)

type myServer struct {
    hellopb.UnimplementedGreetingServiceServer
}

// 自作サービス構造体のコンストラクタを定義
func NewMyServer() *myServer {
    return &myServer{}
}

// Hello メソッドの実装
func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
    // リクエストからnameフィールドを取り出して
    // "Hello, [名前]!"というレスポンスを返す
    return &hellopb.HelloResponse{
        Message: fmt.Sprintf("Hello, %s!", req.GetName()),
    }, nil
}

func main() {
    // 1. 8080番portのLisnterを作成
    port := 8080
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        panic(err)
    }

    // 2. gRPCサーバーを作成
    s := grpc.NewServer()

    // 3. gRPCサーバーにGreetingServiceを登録
    hellopb.RegisterGreetingServiceServer(s, NewMyServer())

    // 4. サーバーリフレクションの設定
    reflection.Register(s)

    // 4. 作成したgRPCサーバーを、8080番ポートで稼働させる
    go func() {
        log.Printf("start gRPC server port: %v", port)
        s.Serve(listener)
    }()

    // 5.Ctrl+Cが入力されたらGraceful shutdownされるようにする
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, os.Interrupt)
    <-quit
    log.Println("stopping gRPC server...")
    s.GracefulStop()
}