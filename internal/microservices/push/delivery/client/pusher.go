package push_client

import (
	"github.com/mailru/easyjson"
	"github.com/streadway/amqp"
	models "glide/internal/microservices/push"
	"glide/internal/pkg/rabbit"
	"time"
)

type PushSender struct {
	session *rabbit.Session
}

func NewPushSender(session *rabbit.Session) *PushSender {
	return &PushSender{
		session: session,
	}
}

func (ph *PushSender) NewMessage(messageId int64, companion string) error {
	push := &models.MessageInfo{
		MessageId: messageId,
		Companion: companion,
		Date:      time.Now(),
	}

	publish := amqp.Publishing{
		Type: "text/plain",
		Body: []byte{},
	}

	var err error
	publish.Body, err = easyjson.Marshal(push)
	if err != nil {
		return err
	}
	ch := ph.session.GetChannel()

	err = ch.Publish(
		ph.session.GetName(),
		models.MessagePush,
		false,
		false,
		publish,
	)

	return err
}

func (ph *PushSender) NewGlideMessage(companion string, glideId int64) error {
	push := &models.GlideInfo{
		GlideId:   glideId,
		Companion: companion,
		Date:      time.Now(),
	}

	publish := amqp.Publishing{
		Type: "text/plain",
		Body: []byte{},
	}

	var err error
	publish.Body, err = easyjson.Marshal(push)
	if err != nil {
		return err
	}
	ch := ph.session.GetChannel()

	err = ch.Publish(
		ph.session.GetName(),
		models.GlidePush,
		false,
		false,
		publish,
	)

	return err
}
