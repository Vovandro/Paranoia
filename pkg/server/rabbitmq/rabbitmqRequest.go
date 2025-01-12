package rabbitmq

import (
	"bytes"
	amqp "github.com/rabbitmq/amqp091-go"
	"io"
)

type rabbitmqRequest struct {
	msg     *amqp.Delivery
	cookies ICookie
	headers IHeader
}

func (t *rabbitmqRequest) Fill(msg *amqp.Delivery) {
	t.msg = msg

	if t.cookies == nil {
		t.cookies = &RabbitmqCookie{}
	}

	if t.headers == nil {
		t.headers = &rabbitmqHeader{}
	}

	t.headers.(*rabbitmqHeader).Fill(t.msg.Headers)
}

func (t *rabbitmqRequest) GetBody() io.ReadCloser {
	return io.NopCloser(bytes.NewReader(t.msg.Body))
}

func (t *rabbitmqRequest) GetBodySize() int64 {
	return int64(len(t.msg.Body))
}

func (t *rabbitmqRequest) GetCookie() ICookie {
	return t.cookies
}

func (t *rabbitmqRequest) GetHeader() IHeader {
	return t.headers
}

func (t *rabbitmqRequest) GetMethod() string {
	return "RABBITMQ"
}

func (t *rabbitmqRequest) GetURI() string {
	return t.msg.ConsumerTag
}

func (t *rabbitmqRequest) GetQuery() IQuery {
	return nil
}

func (t *rabbitmqRequest) GetRemoteIP() string {
	return ""
}

func (t *rabbitmqRequest) GetRemoteHost() string {
	return ""
}

func (t *rabbitmqRequest) GetUserAgent() string {
	return t.msg.RoutingKey
}
