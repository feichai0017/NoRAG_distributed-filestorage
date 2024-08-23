package mq

import "log"

var done chan bool

// StartConsumer: Start a consumer and monitor the message queue
func StartConsumer(qName, cName string, callback func(msg []byte) bool) {
	if !initChannel() {
		log.Println("Failed to initialize channel")
		return
	}
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
