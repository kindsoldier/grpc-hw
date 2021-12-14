package main

import (
    "context"
    "log"
    "time"

    "google.golang.org/grpc"
    "pmapp/services"
)

const (
    address     = "localhost:8080"
    name        = "world"
)

func main() {
    Hello()
//    Reboot()
    Monitor()
}


type Credential struct {
    Payload map[string]string
}
func NewCredential() *Credential {
    payload := make(map[string]string)
    payload["token"] = "12345678"
    payload["auth"] = "qwerty"
    return &Credential{
        Payload: payload,
    }
}

func (this *Credential) GetRequestMetadata(ctx context.Context, data ...string) (map[string]string, error) {
    var err error
    return this.Payload, err
}

func (this *Credential) RequireTransportSecurity() bool {
    return false
}

func Hello() {
    var err error

    cred := NewCredential()
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithPerRPCCredentials(cred))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := services.NewHelloClient(conn)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    result, err := client.Hello(ctx, &services.HelloRequest{Name: name})
    if err != nil {
        log.Fatalf("call error: %v", err)
    }
    log.Printf("result: %s", result.GetMessage())
}


func Reboot() {
    var err error
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := services.NewSystemClient(conn)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    result, err := client.Reboot(ctx, &services.EmptyRequest{})
    if err != nil {
        log.Fatalf("call error: %v", err)
    }
    log.Printf("result: %s", result)
}

func Monitor() {
    var err error

    cred := NewCredential()
    conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithPerRPCCredentials(cred))
    if err != nil {
        log.Fatalf("did not connect: %v", err)
    }
    defer conn.Close()

    client := services.NewSystemClient(conn)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    monitor, _ := client.Monitor(ctx)
    for {
        req := services.EmptyRequest{}
        monitor.Send(&req)
        time.Sleep(1)
        measure, _ := monitor.Recv()
        log.Println(measure)
    }
}

