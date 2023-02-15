package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var helloGrpc *HelloGrpc

func main() {
	helloGrpc = ConnectHelloGrpc()
	defer helloGrpc.closeGrpc()

	logrus.SetFormatter(&logrus.TextFormatter{})

	router := gin.Default()

	router.GET("/hello", helloToWho)
	router.GET("/helloServerStream", helloServerStream)
	router.GET("/helloClientStream", helloClientStream)
	router.GET("/helloAllStream", helloAllStream)

	server := &http.Server{
		Addr:    ":8088",
		Handler: router,
	}

	wait := make(chan struct{})
	go showDownServer(wait, server)

	server.ListenAndServe()
	<-wait
	logrus.Infoln("ok 88")
}

func showDownServer(s chan struct{}, server *http.Server) {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit

	context, cancel := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancel()

	err := server.Shutdown(context)
	if err != nil {
		logrus.Errorln(err)
	} else {
		logrus.Infoln("88")
	}
	close(s)
}
