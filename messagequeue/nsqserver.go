package messagequeue

import (
	"log"

	"github.com/nsqio/go-nsq"
)

// NSQD configuration
var (
	cfg = nsq.NewConfig()
	// domain = "localhost:4150"
	domain = "172.17.0.1:4150"
	topic  = "GALLERY"
)

// Notify the modification done in gallery
func Notify(message string) {

	cfg := nsq.NewConfig()
	w, err := nsq.NewProducer(domain, cfg)
	if err != nil {
		log.Printf("NSQD:: Unable tp create producer topic[%s], Message[%s]. err:(%q)", topic, message, err)
		return
	}

	if err := w.Publish(topic, []byte(message)); err != nil {
		log.Printf("NSQD:: Unable notify topic[%s], Message[%s]. err:(%q)", topic, message, err)
	}

	w.Stop()
}
