package main

import (
	"context"
	"flag"
	"log"

	"github.com/callisto13/mashed-potatoes/jokeprovider/joke"
	"github.com/callisto13/mashed-potatoes/party/natsemitter"
	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

const (
	DEFAULT_NATS_SERVER = "http://localhost:4222"
	SUBJECT             = "jokeprovider"
)

func main() {
	var (
		natsAddress string
		ctx         = context.Background()
	)

	flag.StringVar(&natsAddress, "nats-address", DEFAULT_NATS_SERVER, "address + port of running NATs service")

	p, err := cenats.NewConsumer(natsAddress, SUBJECT, cenats.NatsOptions())
	if err != nil {
		log.Fatalf("failed to create nats protocol, %s", err.Error())
	}

	defer p.Close(ctx)

	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %s", err.Error())
	}

	log.Println("starting provider, waiting for events")

	for {
		if err := c.StartReceiver(ctx, receive); err != nil {
			log.Printf("failed to start nats receiver, %s", err.Error())
		}
	}
}

func receive(ctx context.Context, event cloudevents.Event) error {
	data := &natsemitter.EventData{}
	if err := event.DataAs(data); err != nil {
		log.Printf("got data error: %s\n", err.Error())
	}

	log.Printf("received action: %s", data.Action)

	switch data.Action {
	case "enrol":
		if err := joke.Enrol(); err != nil {
			return err
		}
	default:
		log.Printf("no response found for action: %s", data.Action)
	}

	log.Printf("action %s completed successfully", data.Action)

	return nil
}
