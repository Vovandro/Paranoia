package rabbitmq

type RabbitmqResponse struct {
	Body       []byte
	StatusCode int
	headers    IHeader
	cookie     ICookie
}

func (t *RabbitmqResponse) Clear() {
	t.Body = t.Body[:0]
	t.StatusCode = 200
	t.headers = &rabbitmqHeader{
		make(map[string][]string, 10),
	}
	t.cookie = &RabbitmqCookie{}
}

func (t *RabbitmqResponse) SetBody(data []byte) {
	t.Body = data
}

func (t *RabbitmqResponse) GetBody() []byte {
	return t.Body
}

func (t *RabbitmqResponse) SetStatus(status int) {
	t.StatusCode = status
}

func (t *RabbitmqResponse) GetStatus() int {
	return t.StatusCode
}

func (t *RabbitmqResponse) Header() IHeader {
	return t.headers
}

func (t *RabbitmqResponse) Cookie() ICookie {
	return t.cookie
}
