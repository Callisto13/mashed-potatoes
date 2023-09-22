package httpsender

import (
	"context"
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
)

type EventData struct {
	Action string
	Data   map[string]interface{}
}

func SendEvent(target, action string) error {
	ctx := cloudevents.ContextWithTarget(context.Background(), target)

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Printf("failed to create protocol: %s", err.Error())

		return err
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Printf("failed to create client, %v", err)

		return err
	}

	e := cloudevents.NewEvent()
	e.SetType("com.mashedpotato.party.sent")
	e.SetSource("https://github.com/callisto13/mashed-potatoes/party/httpsender")
	e.SetTime(time.Now())
	e.SetDataSchema("some-schema://party.httpsender.http.EventData")
	if err := e.SetData(cloudevents.ApplicationJSON, &EventData{
		Action: action,
		Data:   map[string]interface{}{"unstructured": "data"},
	},
	); err != nil {
		return err
	}

	result := c.Send(ctx, e)
	if cloudevents.IsUndelivered(result) {
		log.Printf("failed to send: %v", result)

		return err
	}

	var httpResult *cehttp.Result
	if cloudevents.ResultAs(result, &httpResult) {
		log.Printf("sent with status code %d", httpResult.StatusCode)
	} else {
		log.Printf("send did not return an HTTP response: %s", result)
	}

	return nil
}
