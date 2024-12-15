package srvUtils

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type rabbitmqHeader struct {
	data map[string][]string
}

func (t *rabbitmqHeader) Fill(headers amqp.Table) {
	t.data = make(map[string][]string, len(headers))

	for k, v := range headers {
		t.data[k] = []string{v.(string)}
	}
}

func (t *rabbitmqHeader) Add(key, value string) {
	if _, ok := t.data[key]; !ok {
		t.data[key] = make([]string, 0)
	}

	t.data[key] = append(t.data[key], value)
}

func (t *rabbitmqHeader) Set(key string, value string) {
	t.data[key] = []string{value}
}

func (t *rabbitmqHeader) Get(key string) string {
	if v, ok := t.data[key]; ok && len(v) > 0 {
		return v[0]
	}

	return ""
}

func (t *rabbitmqHeader) Values(key string) []string {
	if v, ok := t.data[key]; ok {
		return v
	}

	return nil
}

func (t *rabbitmqHeader) Del(key string) {
	delete(t.data, key)
}

func (t *rabbitmqHeader) GetAsMap() map[string][]string {
	return t.data
}
