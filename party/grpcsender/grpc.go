package grpcsender

import (
	"context"
	"encoding/json"
	"time"

	"github.com/callisto13/mashed-potatoes/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type EventData struct {
	Action string
	Data   map[string]interface{}
}

func SendEvent(target, action string) error {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	c := proto.NewProviderServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	dat := EventData{
		Action: "enrol",
		Data:   map[string]interface{}{"unstructured": "data"},
	}

	data, err := json.Marshal(dat)
	if err != nil {
		return err
	}

	_, err = c.Event(ctx, &proto.EventRequest{
		Event: &proto.CloudEvent{
			Source: "https://github.com/callisto13/mashed-potatoes/party/grpcsender",
			Type:   "com.mashedpotato.party.sent",
			Data: &proto.CloudEvent_BinaryData{
				BinaryData: data,
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
