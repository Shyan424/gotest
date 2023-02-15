package mq

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

type NatsRequest struct {
	Subject string
	Message string
	Timeout time.Duration
}

type NatsResponse struct {
	Message string
	NatsMsg *nats.Msg
}

func (r *NatsRequest) SendRequest() *NatsResponse {
	if r.Subject == "" {
		r.Subject = "re"
	}
	if r.Timeout == 0 {
		r.Timeout = 1 * time.Second
	}

	return r.ruequestHandle()
}

func (r *NatsRequest) ruequestHandle() *NatsResponse {
	ms, err := nc.Request(r.Subject, []byte(r.Message), r.Timeout)
	if err != nil {
		logrus.Errorf("Response error: %s", err)
	}

	return &NatsResponse{
		Message: string(ms.Data),
		NatsMsg: ms,
	}
}
