package mq

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type SubscribeMsg struct {
	Subject  string
	Durable  string
	MaxFetch int
}

type msgHandle func(who string, msg string)

func SubMsg() {
	sub := SubscribeMsg{
		Subject:  DEFAULT_SUBJECT,
		Durable:  "sub",
		MaxFetch: 5,
	}

	sub.PullSubscribe("first", hnadleMsg)
	// 刊刊開兩ㄍ會怎樣
	sub.PullSubscribe("second", hnadleMsg)

	ResponseHandle("Your friend")
	ResponseHandle("No one")

	noAckSubscript()
}

func hnadleMsg(who string, m string) {
	logrus.Info("%s handle %s ohhhyaaaa", who, m)
}

// 用pullSubscript
func (s *SubscribeMsg) PullSubscribe(who string, msgHandle msgHandle) error {
	if s.Subject == "" {
		return errors.New("subject empty")
	}
	if s.Durable == "" {
		return errors.New("durable empty")
	}
	if s.MaxFetch < 1 {
		s.MaxFetch = 1
	}

	sub, err := jsConn.PullSubscribe(s.Subject, s.Durable)
	if err != nil {
		return err
	}

	go func() {
		for {
			msgs, err := sub.Fetch(s.MaxFetch)
			if err != nil {
				if err == nats.ErrTimeout {
					continue
				}

				logrus.Fatal(err)
			}

			for _, m := range msgs {
				msgHandle(who, string(m.Data))
				m.Ack()
			}
		}
	}()

	return nil
}

func (s *SubscribeMsg) ChanSubscribe(who string, handle msgHandle) {
	ch := make(chan *nats.Msg)

	_, err := jsConn.ChanQueueSubscribe("jetStream.chan", "ch", ch)
	if err != nil {
		logrus.Errorf("ChanQueue error: %s", err)
	}

	for {
		msg := <-ch

		handle(who, string(msg.Data))
	}
}

// 用chansubscript
func ResponseHandle(who string) {
	go func() {
		subChan := make(chan *nats.Msg)
		sub, err := nc.ChanSubscribe("re", subChan)
		if err != nil {
			logrus.Fatalf("Response is fail: %s", err)
		}
		defer sub.Drain()

		for {
			msg := <-subChan

			nc.Publish(msg.Reply, []byte(fmt.Sprintf("%s hnadle %s", who, msg.Data)))
		}
	}()
}

// 測試 noack 是否會重複發送
// 測完差不多是下面那樣
// 預設ack設定為 nats.AckExplicit() 必須手動ack
// 預設redelivery為30秒
// 使用會自動建立consumer所以如果要修改consumer設定就必須使用js.UpdateConsumer
func noAckSubscript() {
	sub, err := jsConn.PullSubscribe("jetStream.noack", "me", nats.AckWait(10*time.Second))
	if err != nil {
		logrus.Fatalf("noack err: %s", err)
	}

	go func() {
		for {
			msgs, err := sub.Fetch(1)
			if err != nil {
				if err == nats.ErrTimeout {
					continue
				}
				logrus.Errorf("fetch err: %s", err)
			}

			for _, m := range msgs {
				logrus.Infof("noack handle msg: %s", m.Data)
			}
		}
	}()

}
