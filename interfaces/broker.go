package interfaces

type IBroker interface {
	Init(app IService) error
	Stop() error
	String() string

	Publish(topic string, m *Message) error
	Subscribe(topic string, s *Subscriber) error
}

type Event interface {
	Topic() string
	Message() *Message
	Ack() error
	Error() error
}

type Message struct {
	Header map[string]string
	Body   []byte
}

type Subscriber interface {
	Handler(event *Event) error
	Unsubscribe() error
}
