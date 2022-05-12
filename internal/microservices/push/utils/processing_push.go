package utils

import (
	"bytes"
	"github.com/mailru/easyjson"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"glide/internal/microservices/push"
	"glide/internal/microservices/push/push/usecase"
	"glide/internal/pkg/rabbit"
)

type SendMessager interface {
	SendMessage(users []string, hsg easyjson.Marshaler)
}

type ProcessingPush struct {
	session *rabbit.Session
	logger  *logrus.Entry
	sendMsg SendMessager
	usecase usecase.Usecase
	stop    chan bool
}

func NewProcessingPush(logger *logrus.Entry, session *rabbit.Session,
	sendMsg SendMessager, usecase usecase.Usecase) *ProcessingPush {
	return &ProcessingPush{
		session: session,
		sendMsg: sendMsg,
		logger:  logger,
		usecase: usecase,
		stop:    make(chan bool),
	}
}

func (pp *ProcessingPush) Stop() {
	pp.stop <- true
}

func (pp *ProcessingPush) initMsg(routerKey string) (<-chan amqp.Delivery, error) {
	ch := pp.session.GetChannel()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		return nil, err
	}

	if err = ch.QueueBind(q.Name, routerKey, pp.session.GetName(), false, nil); err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, err
}

func (pp *ProcessingPush) RunProcessMessage() {
	msg, err := pp.initMsg(push.MessagePush)
	if err != nil {
		pp.logger.Errorf("error init message query from msg with err: %s", err)
		return
	}
	pp.processMessageMsg(msg)
}

func (pp *ProcessingPush) RunProcessGlide() {
	msg, err := pp.initMsg(push.GlidePush)
	if err != nil {
		pp.logger.Errorf("error init glide query from msg with err: %s", err)
		return
	}
	pp.processGlideMsg(msg)
}

func (pp *ProcessingPush) processMessageMsg(msg <-chan amqp.Delivery) {
	for {
		var pushMsg amqp.Delivery
		select {
		case <-pp.stop:
			return
		case pushMsg = <-msg:
			break
		}

		msgInfo := &push.MessageInfo{}
		reader := bytes.NewBuffer(pushMsg.Body)
		if err := easyjson.UnmarshalFromReader(reader, msgInfo); err != nil {
			pp.logger.Errorf("error decode info message from msg with err: %s", err)
			continue
		}

		users, sendPush, err := pp.usecase.PrepareMessagePush(msgInfo)
		if err != nil {
			pp.logger.Errorf("error prepare info message with err: %s", err)
			continue
		}
		pp.logger.Infof("Was send message about new message %v", pushMsg.Body)
		pp.sendMsg.SendMessage(users, PushResponse{Type: push.MessagePush, Push: sendPush})
	}
}

func (pp *ProcessingPush) processGlideMsg(msg <-chan amqp.Delivery) {
	for {
		var pushMsg amqp.Delivery
		select {
		case <-pp.stop:
			return
		case pushMsg = <-msg:
			break
		}
		glide := &push.GlideInfo{}
		reader := bytes.NewBuffer(pushMsg.Body)
		if err := easyjson.UnmarshalFromReader(reader, glide); err != nil {
			pp.logger.Errorf("error decode info glide from msg with err: %s", err)
			continue
		}

		users, sendPush, err := pp.usecase.PrepareGlidePush(glide)
		if err != nil {
			pp.logger.Errorf("error prepare info glide with err: %s", err)
			continue
		}
		pp.logger.Infof("Was send message about new glide %v", pushMsg.Body)
		pp.sendMsg.SendMessage(users, PushResponse{Type: push.GlidePush, Push: sendPush})
	}
}
