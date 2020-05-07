package main

import (
	"log"
	"sync"

	"github.com/nsqio/go-nsq"
)

// NSQD configuration
var (
	cfg = nsq.NewConfig()
	// domain = "localhost:4150"
	domain = "172.17.0.1:4150"
	topic  = "GALLERY"
)

func main() {

	// Add waitgroup to keep the server alive
	wg := &sync.WaitGroup{}
	wg.Add(1)

	// Create consumer to the receive Notification regarding album/image modification
	q, err := nsq.NewConsumer(topic, "ch", cfg)
	if err != nil {
		log.Panicf("Unable to connect NSQD. domain[%s], topi[%s], err:%q", domain, topic, err)
	}

	log.Printf("Connect to the NSQD[%s], topic[%s]", domain, topic)

	counter := 0
	q.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		counter++
		log.Printf("Seq: [%d], Message: %s", counter, string(message.Body))
		return nil
	}))

	if err := q.ConnectToNSQD(domain); err != nil {
		log.Panicf("Unable to connect NSQD. domain[%s], err:%q", domain, err)
	}

	wg.Wait()
}
