package main

import (
	"context"
	"log"
	"net"
    "io"

	"pmapp/services"
	"google.golang.org/grpc"
    "google.golang.org/grpc/metadata"
)

const (
	port = ":8080"
)

type SystemServer struct {
    services.UnimplementedSystemServer
}

func (this *SystemServer) Reboot(ctx context.Context, request *services.EmptyRequest) (*services.EmptyReply, error) {
    var err error
    log.Println("call reboot")
    return &services.EmptyReply{}, err
}

func (this *SystemServer) Monitor(srv services.System_MonitorServer) error {
    var err error
    log.Println("call monitor")
    for {
        _, err := srv.Recv()
        meta, _ := metadata.FromIncomingContext(srv.Context())            
        log.Println("request meta:", meta)
        log.Println("monitor request")
        if err == io.EOF {
            return nil
        }
        if err != nil {
            return err
        }
        res := services.Measure{ Power: 10 }
        srv.Send(&res)
    }
    return err
}

func NewSystemServer() *SystemServer {
	return &SystemServer{}
}


func (this *HelloServer) Hello(ctx context.Context, request *services.HelloRequest) (*services.HelloReply, error) {
	var err error
	log.Printf("requst name: %v", request.GetName())
	return &services.HelloReply{ Message: "hello " + request.GetName() }, err
}
type HelloServer struct {
	services.UnimplementedHelloServer
}
func NewHelloServer() *HelloServer {
	return &HelloServer{}
}

func unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    log.Println("unary interceptor: ", info.FullMethod)
    meta, _ := metadata.FromIncomingContext(ctx)            
    log.Println("reqest token:", meta["token"])
    return handler(ctx, req)
}

func streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
    log.Println("stream interceptor: ", info.FullMethod)
    meta, _ := metadata.FromIncomingContext(stream.Context())            
    log.Println("reqest token:", meta["token"])
    return handler(srv, stream)
}

func main() {
	var err error
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	serv := grpc.NewServer(
                    grpc.UnaryInterceptor(unaryInterceptor),
                    grpc.StreamInterceptor(streamInterceptor),
    )

	services.RegisterHelloServer(serv, NewHelloServer())
	services.RegisterSystemServer(serv, NewSystemServer())

	log.Printf("server listening at %v", listener.Addr())

    err = serv.Serve(listener)
	if err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
