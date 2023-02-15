package mq

import (
	"encoding/json"
)

type PublishMsg struct {
	Subject string
	Message interface{}
}

func (p *PublishMsg) Publish() error {
	var msg []byte
	switch p.Message.(type) {
	case string:
		msg = []byte(p.Message.(string))
	default:
		jMsg, err := json.Marshal(p.Message)
		if err != nil {
			return err
		}

		msg = jMsg
	}

	_, err := jsConn.Publish(p.Subject, msg)
	if err != nil {
		return err
	}

	return nil
}
