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

}

func (ph *PushSender) NewGlideMessage(companion int64, glideId int64) error {

}

func (ph *PushSender) NewPost(creatorId int64, postId int64, postTitle string) error {
	push := &models.PostInfo{
		CreatorId: creatorId,
		PostId:    postId,
		PostTitle: postTitle,
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
		models.PostPush,
		false,
		false,
		publish,
	)

	return err
}

func (ph *PushSender) NewComment(commentId int64, authorId int64, postId int64) error {
	push := &models.CommentInfo{
		CommentId: commentId,
		AuthorId:  authorId,
		PostId:    postId,
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
		models.CommentPush,
		false,
		false,
		publish,
	)

	return err
}

func (ph *PushSender) NewSubscriber(subscriberId int64, awardsId int64, creatorId int64) error {
	push := &models.SubInfo{
		UserId:    subscriberId,
		CreatorId: creatorId,
		AwardsId:  awardsId,
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
		models.CommentPush,
		false,
		false,
		publish,
	)

	return err
}
