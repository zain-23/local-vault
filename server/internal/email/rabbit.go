package email

import (
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// One place for every RabbitMQ object name so publisher and consumer always agree
const (
	exchangeName 	= "email"
	queueSend		= "email.send"
	queueRetry		= "email.retry"
	queueDead		= "email.dead"
	routingSend		= "send"
	routingDead		= "dead"
)

// Connect dial the broker and opens one channel (channels carry all AMQP operations)
func Connect(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url) // opens the TCP connection to RabbitMQ
	if err != nil {
		return  nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close() // don't leak the connection if opening the channel failsre
		return  nil, nil, err
	}

	return conn, ch, nil
}

// DeclareTopology creates the exchange + queues. It's idempotent, so every binary can call it.
func DeclareTopology(ch *amqp.Channel, retryDelay time.Duration) error {
	// Direct exchange = route by exact routing-key match, durable=true survives broker restart
	if err := ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil); err != nil {
		return err
	}

	// Main work queue - durable so queued jobs servive a broker restart
	if _, err := ch.QueueDeclare(queueSend, true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind(queueSend, routingSend, exchangeName, false, nil); err != nil {
		return err
	}

	// Retry queue - NO consumer, when a message TTL expires it dead-letters back to email.send
	retryArgs := amqp.Table{
		"x-message-ttl": 				int32(retryDelay.Microseconds()), // how long a job waits here
		"x-dead-letter-exchange":		exchangeName,  // where it goes after ttl
		"x-dead-letter-routing-key":	routingSend, // enters the work queue
	}

	if _, err := ch.QueueDeclare(queueRetry, true, false, false, false, retryArgs); err != nil {
		return  err
	}

	// Parking lot - durable queue we inspect by hand; bound with the "dead" key
	if _, err := ch.QueueDeclare(queueDead, true, false, false, false, nil); err != nil {
		return err
	}

	if err := ch.QueueBind(queueDead, routingDead, exchangeName, false, nil); err != nil {
		return err
	}

	return nil
}
