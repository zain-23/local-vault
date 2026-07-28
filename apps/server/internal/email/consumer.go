package email

import (
	"context"
	"encoding/json"
	"log"

	ampq "github.com/rabbitmq/amqp091-go"
)

const headerRetries = "x-retries" // custom header tracking how many times a job has been retried

// Consumer reads jobs, sends them, and routes failures to retry/park - held by the worker
type Consumer struct {
	ch			*ampq.Channel
	sender		*Sender
	maxRetries	int
}

// NewConsumer wires the channel, the Resend sender, and the retry limit together
func NewConsumer(ch *ampq.Channel, sender *Sender, maxRetries int) *Consumer {
	return &Consumer{
		ch: ch,
		sender: sender,
		maxRetries: maxRetries,
	}
}

// Run starts comsuming and blocks, handling one delivery at a time until ctx is cancelled
func (c *Consumer) Run (ctx context.Context) error {
	// Consume returns a Go channel streaming deliveries. autoAck=false: we ack manually
	msgs, err := c.ch.Consume(queueSend, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Printf("Email worker waiting for jobs on %s", queueSend)

	for {
		select {
		case <- ctx.Done():
			return ctx.Err()
		case d, ok := <-msgs:
			if !ok {
				return nil
			}
			c.handle(d)
		}
	}
}

// currentRetries reads x-retries from the header, defaulting to 0 for a first attempt.
func currentRetries(d ampq.Delivery) int {
	if v, ok := d.Headers[headerRetries]; ok {
		if n, ok := v.(int32); ok { // AMQP encodes integer headers as int32
			return int(n)
		}
	}
	return 0
}

// retry republishes to email.retry, which holds the job (TTL) then bounces it back to send.
func (c *Consumer) retry(d ampq.Delivery, attempt int) {
	headers := ampq.Table{headerRetries: int32(attempt)} // carry the new count forward
	// Publishing to the default exchange ("") with routing key = queue name targets that queue.
	c.ch.PublishWithContext(context.Background(), "", queueRetry, false, false, ampq.Publishing{
		ContentType:  "application/json",
		DeliveryMode: ampq.Persistent,
		Headers:      headers,
		Body:         d.Body,
	})
}

// park sends an exhausted/broken job to the parking lot with the reason it failed.
func (c *Consumer) park(d ampq.Delivery, reason string) {
	headers := ampq.Table{"x-error": reason} // record why it ended up here
	c.ch.PublishWithContext(context.Background(), exchangeName, routingDead, false, false, ampq.Publishing{
		ContentType:  "application/json",
		DeliveryMode: ampq.Persistent,
		Headers:      headers,
		Body:         d.Body,
	})
}


// handle processes one delivery: send, then ack + (retry | park) as needed.
func (c *Consumer) handle(d ampq.Delivery) {
	var job EmailJob
	if err := json.Unmarshal(d.Body, &job); err != nil {
		log.Printf("⚠️  bad message, parking: %v", err) // unparseable can never succeed
		c.park(d, "invalid json")
		d.Ack(false)
		return
	}
	log.Printf("Sending email to:%s", job.To)
	if err := c.sender.Send(job); err != nil {
		retries := currentRetries(d) // how many attempts this message has already had
		if retries+1 < c.maxRetries {
			log.Printf("↻ send failed (attempt %d), retrying: %v", retries+1, err)
			c.retry(d, retries+1)
		} else {
			log.Printf("✗ send failed permanently, parking: %v", err)
			c.park(d, err.Error())
		}
		d.Ack(false) // ack the original — we've re-published it to retry/park
		return
	}

	d.Ack(false) // success — remove it from the work queue
}
