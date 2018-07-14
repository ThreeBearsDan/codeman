package main

import (
	"os"
	"os/signal"

	nsq "github.com/nsqio/go-nsq"
)

func main() {
	conf := nsq.NewConfig()
	conf.MaxInFlight = 1

	// producer
	go func() {
		producer, _ := nsq.NewProducer(":4150", conf)
		// producer message
		for i := 0; i < 10; i++ {
			producer.Publish("apple", []byte("Message from producer\n"))
		}
	}()

	// consumer
	go func() {
		consumer, _ := nsq.NewConsumer("apple", "green", conf)
		// add handler
		consumer.AddHandler(nsq.HandlerFunc(func(m *nsq.Message) error {
			m.WriteTo(os.Stdout)
			return nil
		}))

		consumer.ConnectToNSQLookupd(":4161")
	}()

	exit := make(chan os.Signal)
	signal.Notify(exit, os.Kill, os.Interrupt)
	<-exit
}
