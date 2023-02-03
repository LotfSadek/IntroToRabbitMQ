package main

import (
	"testing"
	"time"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
)

func TestPublishToRabbitMQ(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	assert.Nil(t, err) // fail test if err not nil
	// ALWAYS MAKE SURE TO HANDLE THE ERROR
	if err != nil {
		panic(err)
	}
	// You can assume that the connection is valid
	defer conn.Close() // we can defer here because we wil be blocking

	ch, err := conn.Channel()
	assert.Nil(t, err) // fail test if err not nil
	// We can assume ch is valid
	defer ch.Close()
	// exchange name is notifs
	// routing key is message
	err = ch.Publish(
		"Notifications1",
		"message",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("sending this message at " + time.Now().Format(time.RFC3339Nano)),
		})
	assert.Nil(t, err)
}
