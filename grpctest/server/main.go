package main

import (
	"bytes"
	"context"
	"fmt"
	"gogo/hello"
	"io"
	"net"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	cred, err := credentials.NewServerTLSFromFile("../tls/server_cert.pem", "../tls/server_key.pem")
	if err != nil {
		logrus.Fatalf("Credentials Fail %s", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		logrus.Fatalln(err)
	}

	server := grpc.NewServer(grpc.Creds(cred))
	hello.RegisterHelloToWhoServer(server, &helloServer{})
	if err := server.Serve(lis); err != nil {
		logrus.Fatalln(err)
	}

}

type helloServer struct {
	hello.UnimplementedHelloToWhoServer
}

// 實做xxx_grpc.pb內的方法 type HelloToWhoServer interface
func (s *helloServer) Hello(c context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	res := &hello.HelloResponse{Str: fmt.Sprintf("Hello %s", req.GetName())}

	return res, nil
}

// server端stream
func (s *helloServer) HelloServerStream(res *hello.HelloRequest, stream hello.HelloToWho_HelloServerStreamServer) error {
	who := res.GetName()

	for i := 0; i < 10; i++ {
		stream.Send(&hello.HelloResponse{Str: fmt.Sprintf("Hi No.%d %s", i, who)})
	}

	return nil
}

// client端stream
func (s *helloServer) HelloClientStream(stream hello.HelloToWho_HelloClientStreamServer) error {
	names := []string{}

	for {
		who, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		names = append(names, who.GetName())
	}

	buffer := bytes.NewBufferString("")
	for i, name := range names {
		buffer.WriteString(fmt.Sprintf("Hi No.%d %s\n", i, name))
	}

	stream.SendAndClose(&hello.HelloResponse{Str: buffer.String()})

	return nil
}

// 雙向stream
func (s *helloServer) HelloAllStream(stream hello.HelloToWho_HelloAllStreamServer) error {
	i := 0
	for {
		whoReq, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		stream.Send(&hello.HelloResponse{Str: fmt.Sprintf("Hi NO.%d %s", i, whoReq.GetName())})
		i++
	}
}
