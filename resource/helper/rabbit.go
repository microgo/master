package helper

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type HandlerFunc func(*amqp.Delivery)

func (r *Helper) MakeChanel(name string) (*amqp.Channel, error) {
	ch, err := r.Rabbit.Channel()
	if err != nil {
		fmt.Println("[ERROR] Failed open chanel", err)
		return nil, err
	}
	_, err = ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return ch, err
	}
	err = ch.Qos(1, 0, false)
	return ch, err
}

func (r *Helper) GetNumberOfMessages(name string, ch *amqp.Channel) (int, error) {
	queue, err := ch.QueueDeclare(name, false, false, false, false, nil)
	return queue.Messages, err
}

func (r *Helper) MakeConsume(name string, ch *amqp.Channel, handler HandlerFunc) {
	msgs, err := ch.Consume(
		name,  // queue
		"Boo", // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		fmt.Println("[ERROR] Make consume failed", err)
		panic(err)
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			handler(&d)
			err := d.Ack(false)
			if err != nil {
				d.Ack(false)
			}
		}
	}()
	<-forever
}

func (r *Helper) MakeConsumeWithTag(name string, tag string, ch *amqp.Channel, handler HandlerFunc) {
	msgs, err := ch.Consume(
		name,  // queue
		tag,   // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)

	if err != nil {
		fmt.Println("[ERROR] Make consume with tag failed", err)
		panic(err)
	}
	forever := make(chan bool)
	go func() {
		for d := range msgs {
			handler(&d)
			err := d.Ack(false)
			if err != nil {
				d.Ack(false)
			}
		}
	}()
	<-forever
}

func (r *Helper) PublishMessage(ch *amqp.Channel, queue string, form interface{}) error {
	value, err := json.Marshal(form)
	if err != nil {
		return err
	}
	err = ch.Publish(
		"",    // exchange
		queue, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        value,
		})
	return err
}
