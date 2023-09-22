package main

import (
	"context"
	"log"

	// _ "github.com/cloudevents/sdk-go/binding/format/protobuf/v2"
	"github.com/callisto13/mashed-potatoes/jokeprovider/joke"
	"github.com/callisto13/mashed-potatoes/party/httpsender"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

func main() {
	ctx := context.Background()
	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Println("starting provider on :8080, waiting for events")

	if err := c.StartReceiver(ctx, receive); err != nil {
		log.Printf("failed to start nats receiver, %s", err.Error())
	}
}

func receive(ctx context.Context, event cloudevents.Event) error {
	// payload := &proto.Sample{}
	data := &httpsender.EventData{}
	if err := event.DataAs(data); err != nil {
		log.Printf("failed to decode data: %s", err)

		return err
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
