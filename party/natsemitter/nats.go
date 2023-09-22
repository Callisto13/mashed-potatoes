package natsemitter

import (
	"context"
	"errors"
	"log"
	"time"

	cenats "github.com/cloudevents/sdk-go/protocol/nats/v2"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
)

type EventData struct {
	Action string
}

func SendEvent(server, subject, action string) error {
	p, err := cenats.NewSender(server, subject, cenats.NatsOptions())
	if err != nil {
		return err
	}

	defer p.Close(context.Background())

	c, err := cloudevents.NewClient(p)
	if err != nil {
		return err
	}

	e := cloudevents.NewEvent()
	e.SetID(uuid.New().String())
	e.SetType("com.mashedpotato.party.sent")
	e.SetTime(time.Now())
	e.SetSource("https://github.com/callisto13/mashed-potatoes/party/natsemitter")
	if err := e.SetData(cloudevents.ApplicationJSON, &EventData{
		Action: action,
	}); err != nil {
		return err
	}

	result := c.Send(context.Background(), e)
	if cloudevents.IsUndelivered(result) {
		log.Println("failed to send event")

		return errors.New("failed to send event")
	}

	log.Printf("sent, accepted: %t", cloudevents.IsACK(result))

	return nil
}
