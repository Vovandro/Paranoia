package kafka

import "github.com/confluentinc/confluent-kafka-go/v2/kafka"

type kafkaHeader struct {
	data map[string][]string
}

func (t *kafkaHeader) Fill(headers []kafka.Header) {
	t.data = make(map[string][]string, len(headers))

	for i := 0; i < len(headers); i++ {
		if _, ok := t.data[headers[i].Key]; !ok {
			t.data[headers[i].Key] = make([]string, 0)
		}

		t.data[headers[i].Key] = append(t.data[headers[i].Key], string(headers[i].Value))
	}
}

func (t *kafkaHeader) Add(key, value string) {
	if _, ok := t.data[key]; !ok {
		t.data[key] = make([]string, 0)
	}

	t.data[key] = append(t.data[key], value)
}

func (t *kafkaHeader) Set(key string, value string) {
	t.data[key] = []string{value}
}

func (t *kafkaHeader) Get(key string) string {
	if v, ok := t.data[key]; ok && len(v) > 0 {
		return v[0]
	}

	return ""
}

func (t *kafkaHeader) Values(key string) []string {
	if v, ok := t.data[key]; ok {
		return v
	}

	return nil
}

func (t *kafkaHeader) Del(key string) {
	delete(t.data, key)
}

func (t *kafkaHeader) GetAsMap() map[string][]string {
	return t.data
}
