package mq

import "log"

var done chan bool

// StartConsume: Start a consumer and monitor the message queue
func StartConsume(qName, cName string, callback func(msg []byte) bool) {

	// Start consumer
	msgs, err := channel.Consume(
		qName,
		cName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Println(err.Error())
		return
	}

	done = make(chan bool)

	go func() {
		for msg := range msgs {
			// Process the message
			processSuc := callback(msg.Body)
			if !processSuc {

			}
		}
	}()

	// Waiting for exit signal
	<-done

	// Close the channel
	channel.Close()
}

// StopConsume: Stop consuming messages
func StopConsume() {
	done <- true
}
