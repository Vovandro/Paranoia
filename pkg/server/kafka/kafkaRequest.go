package kafka

import (
	"bytes"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"io"
)

type kafkaRequest struct {
	msg     *kafka.Message
	cookies ICookie
	headers IHeader
}

func (t *kafkaRequest) Fill(msg *kafka.Message) {
	t.msg = msg

	if t.cookies == nil {
		t.cookies = &KafkaCookie{}
	}

	if t.headers == nil {
		t.headers = &kafkaHeader{}
	}

	t.headers.(*kafkaHeader).Fill(t.msg.Headers)
}

func (t *kafkaRequest) GetBody() io.ReadCloser {
	return io.NopCloser(bytes.NewReader(t.msg.Value))
}

func (t *kafkaRequest) GetBodySize() int64 {
	return int64(len(t.msg.Value))
}

func (t *kafkaRequest) GetCookie() ICookie {
	return t.cookies
}

func (t *kafkaRequest) GetHeader() IHeader {
	return t.headers
}

func (t *kafkaRequest) GetMethod() string {
	return "KAFKA"
}

func (t *kafkaRequest) GetURI() string {
	return *t.msg.TopicPartition.Topic
}

func (t *kafkaRequest) GetQuery() IQuery {
	return nil
}

func (t *kafkaRequest) GetRemoteIP() string {
	return ""
}

func (t *kafkaRequest) GetRemoteHost() string {
	return ""
}

func (t *kafkaRequest) GetUserAgent() string {
	return *t.msg.TopicPartition.Metadata
}
