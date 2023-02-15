package main

import (
	"bytes"
	"context"
	"fmt"
	"gogo/hello"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type HelloGrpc struct {
	conn   *grpc.ClientConn
	Client hello.HelloToWhoClient
}

func ConnectHelloGrpc() *HelloGrpc {
	// 後面的cc.example.com主要是因為沒domain name 所以給一個假的可以去對應 Server_cert.pem那邊設定的domain name
	// 沒...剛測試過，證書可以用IP去簽署
	cert, err := credentials.NewClientTLSFromFile("../tls/ca_cert.pem", "")
	if err != nil {
		logrus.Fatalf("Credentials Fail %s", err)
	}

	conn, err := grpc.Dial("127.0.0.1:50051", grpc.WithTransportCredentials(cert))
	// conn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Println()
	}
	client := hello.NewHelloToWhoClient(conn)

	return &HelloGrpc{
		conn:   conn,
		Client: client,
	}
}

func (h *HelloGrpc) closeGrpc() {
	if err := h.conn.Close(); err != nil {
		logrus.Error(err)
	}
}

func helloToWho(c *gin.Context) {
	who := c.Query("who")

	returnStr, err := helloGrpc.Client.Hello(context.Background(), &hello.HelloRequest{Name: who})
	if err != nil {
		logrus.Errorln(err)
	}

	c.String(http.StatusOK, returnStr.Str)
}

func helloServerStream(c *gin.Context) {
	who := c.Query("who")

	resStream, err := helloGrpc.Client.HelloServerStream(context.Background(), &hello.HelloRequest{Name: who})
	if err != nil {
		logrus.Errorf("Server Stream Error %s", err)
	}

	buffer := bytes.NewBufferString("")
	for {
		res, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			logrus.Errorf("Server Stream Recv Error %s", err)
			break
		}

		buffer.WriteString(fmt.Sprintf("%s\n", res.Str))
	}

	c.String(http.StatusOK, buffer.String())
}

func helloClientStream(c *gin.Context) {
	who := c.Query("who")

	grpcStream, err := helloGrpc.Client.HelloClientStream(context.Background())
	if err != nil {
		logrus.Errorf("Client Stream Error %s", err)
	}

	names := strings.Split(who, ",")
	for i, name := range names {
		grpcStream.Send(&hello.HelloRequest{Name: fmt.Sprintf("%d %s", i, strings.TrimSpace(name))})
	}

	res, err := grpcStream.CloseAndRecv()
	if err != nil {
		logrus.Errorf("Client Stream Recv Error %s", err)
	}

	c.String(http.StatusOK, res.GetStr())
}

func helloAllStream(c *gin.Context) {
	who := c.Query("who")

	grpcStream, err := helloGrpc.Client.HelloAllStream(context.Background())
	if err != nil {
		logrus.Errorf("All Stream Error %s", err)
	}

	names := strings.Split(who, ",")
	for i, name := range names {
		grpcStream.Send(&hello.HelloRequest{Name: fmt.Sprintf("%d %s", i, name)})
	}

	wait := make(chan struct{})
	buffer := bytes.NewBufferString("")
	go func() {
		for {
			res, err := grpcStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				logrus.Errorf("All Stream Recv Error %s", err)
				break
			}

			buffer.WriteString(fmt.Sprintf("%s\n", res.GetStr()))
		}
		close(wait)
	}()

	// 比較需要注意的地方
	grpcStream.CloseSend()
	<-wait
	c.String(http.StatusOK, buffer.String())
}
