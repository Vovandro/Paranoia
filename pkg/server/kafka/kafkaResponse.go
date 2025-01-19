package kafka

type KafkaResponse struct {
	Body       []byte
	StatusCode int
	headers    IHeader
	cookie     ICookie
}

func (t *KafkaResponse) Clear() {
	t.Body = t.Body[:0]
	t.StatusCode = 200
	t.headers = &kafkaHeader{
		make(map[string][]string, 10),
	}
	t.cookie = &KafkaCookie{}
}

func (t *KafkaResponse) SetBody(data []byte) {
	t.Body = data
}

func (t *KafkaResponse) GetBody() []byte {
	return t.Body
}

func (t *KafkaResponse) SetStatus(status int) {
	t.StatusCode = status
}

func (t *KafkaResponse) GetStatus() int {
	return t.StatusCode
}

func (t *KafkaResponse) Header() IHeader {
	return t.headers
}

func (t *KafkaResponse) Cookie() ICookie {
	return t.cookie
}
