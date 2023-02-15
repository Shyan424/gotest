package mq

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

const (
	DEFAULT_STREAM  = "jetStream"
	DEFAULT_SUBJECT = "jetStream.test"
)

var nc *nats.Conn

var jsConn nats.JetStreamContext

func JeStreamConnect() {
	js, err := connectNats().JetStream()
	if err != nil {
		logrus.Fatalf("JetStream connect fail: %s", err)
	}
	jsConn = js

	createStream()
	// createConsumer()
}

func JeStreamDisConnect() {
	nc.Drain()
}

func createStream() {
	//Create a Stream
	jsConn.AddStream(&nats.StreamConfig{
		Name:     DEFAULT_STREAM,
		Subjects: []string{"jetStream.*"},
	})

}

func createConsumer() {
	conf := &nats.ConsumerConfig{
		Durable: "me",
		AckWait: 10 * time.Second,
		// ack policy can not be updated 建立之後就不可更改，並且就算修改consumer也要帶入
		AckPolicy: nats.AckExplicitPolicy,
	}

	// jsConn.AddConsumer(DEFAULT_STREAM, conf)

	info, err := jsConn.UpdateConsumer(DEFAULT_STREAM, conf)
	if err != nil {
		logrus.Fatalf("UpdateConsumer err: %s", err)
	}

	logrus.Debugf("Consumer AckWait: %s", info.Config.AckWait)
}

func connectNats() *nats.Conn {
	natsConn, err := nats.Connect(nats.DefaultURL, nats.UserInfo("aacc", "passowrd"))
	if err != nil {
		logrus.Fatalf("Nats connect fail : %s", err)
	}
	nc = natsConn

	return nc
}
